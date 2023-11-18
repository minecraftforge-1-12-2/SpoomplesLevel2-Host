package server

import (
	"TogetherForever/util"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// --- Client Structs --- //

// Client is a client. It contains the basic information needed to process one.
type Client struct {
	ID       int
	Conn     net.Conn
	Ip256    string
	Name     string
	Admin    bool
	Active   bool
	Paused   bool
	LoggedIn bool
	Lobby    string

	Data    CliData
	LastMsg int64

	Queue    []CompMessage
	QueueMut sync.Mutex
	MsgTries int

	// 1.2.4
	ParseFails int
	Color      string
	Chat       []Message `json:"chat"`
	ChatMut    sync.Mutex
}

// CliData is the data of a client. This is separated because it is constantly updated.
type CliData struct {
	Type    int    `json:"type"`
	Msg     string `json:"msg"`
	Name    string `json:"name"`
	Version string `json:"ver"`
	Lobby   string `json:"lobby"`

	Key  string  `json:"key"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Room uint16  `json:"room"`

	Sprite         string `json:"sprite"`
	Frame          uint8  `json:"frame"`
	Direction      int    `json:"dir"`
	Palette        uint8  `json:"palette"`
	PaletteSprite  string `json:"paletteSprite"`
	PaletteTexture string `json:"paletteTexture"`
	Color          string `json:"color"`

	MsgId int `json:"msgId"`
}

// --- Client Functions --- //

// Accept is called when a client is accepted.
// It sets the client to active, and starts the loop.
func (client *Client) Accept() {
	i := 0
	MainServer.ClientsMut.Lock()
	for _, c := range MainServer.Clients {
		if c.Ip256 == client.Ip256 {
			i++
		}

		if i >= 3 {
			MainServer.ClientsMut.Unlock()
			client.Close(OmsgKick, "You are already connected with the max amount of connections.")
			return
		}
	}
	MainServer.ClientsMut.Unlock()

	client.Chat = make([]Message, 0)
	client.Active = true
	client.LastMsg = time.Now().Unix()

	MainServer.ClientsMut.Lock()
	MainServer.Clients[client.ID] = client
	MainServer.ClientsMut.Unlock()

	client.loop()
}

// loop is the heart of the client. This is where the client is constantly reading and writing!
// It is started when the client is accepted.
func (client *Client) loop() {
	t := time.NewTicker(time.Second / 60)
	for range t.C {
		if !client.Active {
			t.Stop()
			break
		}

		// Reading
		err := client.Conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		buff := make([]byte, 1024)
		_, err = client.Conn.Read(buff)
		data, err := bufio.NewReader(bytes.NewReader(buff)).ReadString('\n')

		if err != nil {
			client.Close(MsgNone, "Timed out")
			break
		}
		
		client.Parse([]byte(data))
		client.LastMsg = time.Now().Unix()

		if client.ParseFails >= 10 {
			client.Close(MsgNone, "Too many invalid packets")
			return
		}

		// Writing
		var js []byte
		if len(client.Queue) > 0 {
			msg := client.Queue[0]

			client.QueueMut.Lock()
			client.Queue = client.Queue[1:]
			client.QueueMut.Unlock()

			js, _ = json.Marshal(map[string]interface{}{
				"type": msg.Type,
				"msg":  msg.Msg,
			})

			_, err := client.Conn.Write(js)
			if err != nil {
				client.Active = false
				client.Close(MsgNone, "Average Disconnection")
				break
			}
		} else {
			// Compacting Clients
			var clients = make([]CompClient, 0)
			MainServer.ClientsMut.Lock()
			for _, c := range MainServer.Clients {
				if c.ID == client.ID || c.Lobby != client.Lobby || c.Data.Room != client.Data.Room {
					continue
				}
				cc := CompClient{
					ID:    c.ID,
					X:     c.Data.X,
					Y:     c.Data.Y,
					Name:  c.Name,
					Admin: c.Admin,
					Room:  c.Data.Room,

					Sprite:         c.Data.Sprite,
					Frame:          c.Data.Frame,
					Direction:      c.Data.Direction,
					Palette:        c.Data.Palette,
					PaletteSprite:  c.Data.PaletteSprite,
					PaletteTexture: c.Data.PaletteTexture,
					Color:          c.Color,
				}

				clients = append(clients, cc)
			}
			MainServer.ClientsMut.Unlock()

			// Marshalling Message
			data := map[string]interface{}{
				"type":      OmsgDefault,
				"loggedIn":  client.LoggedIn,
				"admin":     client.Admin,
				"name":      client.Name,
				"id":        client.ID,
				"onlineCnt": MainServer.LobbyCnt(client.Lobby),
				"clients":   clients,
			}

			client.ChatMut.Lock()
			data["msgs"] = client.Chat
			client.ChatMut.Unlock()

			js, _ = json.Marshal(data)

			_, err := client.Conn.Write(js)
			if err != nil {
				client.Active = false
				client.Close(MsgNone, "Average Disconnection")
				break
			}
		}
	}
}

// Close is called when a client is disconnected.
// It closes the connection, and removes the client from the server.
func (client *Client) Close(typ int, msg string) {
	client.DirWrite(typ, msg)
	util.WriteLine("Client {0} Disconnected: {1}", strconv.Itoa(client.ID), msg)
	_ = client.Conn.Close()
	client.Active = false

	MainServer.ClientsMut.Lock()
	delete(MainServer.Clients, client.ID)
	MainServer.ClientsMut.Unlock()
}

// Parse is called when a client sends a message.
// It parses the message, and sets the client's data. The actions taken depend on the message type.
func (client *Client) Parse(msg []byte) {
	dat := CliData{}

	msg = []byte(strings.Replace(string(msg), "true", "1", -1))
	msg = []byte(strings.Replace(string(msg), "false", "-1", -1))
	err := json.Unmarshal(msg, &dat)

	if err != nil {
		// DISABLED UNTIL IMPROVED
		//util.WriteLine(err.Error())
		//client.ParseFails++
		return
	}

	client.ParseFails = 0

	if MainServer.Anticheat {
		dat.Sprite = util.Anticheat(dat.Sprite)
	}

	client.Paused = false
	switch dat.Type {
	case ImsgDefault:
		if !client.LoggedIn {
			return
		}
		client.Data = dat
		break

	case ImsgPaused:
		if !client.LoggedIn {
			return
		}
		client.Paused = true
		break

	case ImsgLogin:
		if client.LoggedIn {
			return
		}
		// Checking if key is valid
		if MainServer.CheckKey(dat.Key) {
			client.Admin = true
		}
		// Let's not let banned people in...
		if MainServer.CheckBanned(client.Ip256) {
			client.Close(OmsgKick, "You are banned.")
			return
		}
		// NO BAD VERSIONS GOD DAMN IT
		if dat.Version != VERSION {
			client.Close(OmsgKick, "Your client is outdated.")
			return
		}

		dat.Name = util.CleanName(dat.Name)

		// The following section may seem complicated, but all it does is make sure that two people don't have the same name.
		MainServer.ClientsMut.Lock()
		if client.Admin {
			for _, c := range MainServer.Clients {
				if c.Name == dat.Name {
					MainServer.ClientsMut.Unlock()
					c.Close(MsgNone, "")
					MainServer.ClientsMut.Lock()
					break
				}
			}
		} else {
			nn := dat.Name
			for {
				if func() bool {
					for _, c := range MainServer.Clients {
						if c.Name == nn {
							nn = dat.Name + strconv.Itoa(rand.Intn(9999))
							return false
						}
					}
					return true
				}() {
					dat.Name = nn
					break
				}
			}
		}
		MainServer.ClientsMut.Unlock()

		// Setting client data
		client.Data = dat
		client.Name = dat.Name
		client.Lobby = dat.Lobby
		client.LoggedIn = true
		client.Color = dat.Color

		// Telling everyone else that the client has joined.
		util.WriteLine("Client {0} logged in as {1}", strconv.Itoa(client.ID), client.Name)
		MainServer.Broadcast(client.Name+" entered the tower!", client.Lobby)
		client.ServerPm(fmt.Sprintf("Welcome to PTT, %s! Use /help to view commands.", client.Name))
		break

	case ImsgMessage:
		if !client.LoggedIn {
			return
		}

		// Clean that shit out your mouth...
		dat.Msg = util.Clean(dat.Msg, 256)

		if dat.Msg == "" {
			return
		}

		if dat.Msg[0] == '/' {
			client.Command(dat.Msg)
			return
		}

		MainServer.ClientsMut.Lock()
		for _, c := range MainServer.Clients {
			if client.Lobby == c.Lobby && c.Active && c.LoggedIn {
				fmt.Println("preparing message")
				c.Pm(Message{Body: dat.Msg, Username: client.Name, Id: client.ID})
				fmt.Println("msg sent")
			}
		}
		MainServer.ClientsMut.Unlock()
		break
	}
}

// Command is for handling the commands sent by a client.
// These are executed on server-side entirely.
func (client *Client) Command(msg string) {
	arr := strings.Split(msg[1:], " ")
	cmd := arr[0]
	args := arr[1:]

	switch cmd {
	case "help":
		client.ServerPm("- Help -")
		if client.Admin {
			client.ServerPm("/ban <id> <reason> -  Bans a user. Admin only!")
			client.ServerPm("/kick <id> <reason> - Kicks a user. Admin only!")
			client.ServerPm("/announce <message> - Announces a message. UNSTABLE! USE WITH CARE! Admin only!")
		}
		client.ServerPm("/nick <name> - Sets nickname.")
		client.ServerPm("/pm <name> <message> - Private Messages a User.")

	case "ban":
		if !client.Admin {
			client.ServerPm("You don't have permission to use this command.")
			return
		}
		if len(args) < 2 {
			client.ServerPm("Usage")
			client.ServerPm("/ban <id> <reason>")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return
		}
		MainServer.Ban(id, args[1])
		break

	case "kick":
		if !client.Admin {
			client.ServerPm("You don't have permission to use this command.")
			return
		}
		if len(args) < 2 {
			client.ServerPm("Usage")
			client.ServerPm("/kick <id> <reason>")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return
		}
		MainServer.Kick(id, args[1])
		break

	case "announce":
		if !client.Admin {
			client.ServerPm("You don't have permission to use this command.")
			return
		}
		if len(args) < 1 {
			client.ServerPm("Usage")
			client.ServerPm("/announce <message>")
			return
		}
		MainServer.Announce(fmt.Sprintf("%s: %s", client.Name, strings.Join(args, " ")))
		break

	case "nick":
		if len(args) == 0 {
			client.ServerPm("Usage")
			client.ServerPm("/nick <name>")
			return
		}
		name := util.CleanName(strings.Join(args, " "))

		// Same as the login shit
		MainServer.ClientsMut.Lock()
		if client.Admin {
			for _, c := range MainServer.Clients {
				if c.Name == name {
					MainServer.ClientsMut.Unlock()
					c.Close(MsgNone, "")
					MainServer.ClientsMut.Lock()
					break
				}
			}
		} else {
			nn := name
			for {
				if func() bool {
					for _, c := range MainServer.Clients {
						if c.Name == nn {
							nn = name + strconv.Itoa(rand.Intn(9999))
							return false
						}
					}
					return true
				}() {
					name = nn
					break
				}
			}
		}
		MainServer.ClientsMut.Unlock()

		client.Name = name
		break

	case "pm":
		if len(args) < 2 {
			client.ServerPm("Usage")
			client.ServerPm("/pm <name> <message>")
			return
		}
		name := args[0]
		msg := strings.Join(args[1:], " ")

		MainServer.ClientsMut.Lock()
		for _, c := range MainServer.Clients {
			if c.Name == name {
				c.Pm(Message{Body: msg, Username: client.Name + " whispers to you", Id: client.ID})
				break
			}
		}
		MainServer.ClientsMut.Unlock()

		client.Pm(Message{Body: msg, Username: "You whisper to " + name, Id: client.ID})
		break
	}
}

func (client *Client) Append(msg CompMessage) {
	client.QueueMut.Lock()
	client.Queue = append(client.Queue, msg)
	client.QueueMut.Unlock()
}

func (client *Client) DirWrite(typ int, msg string) {
	js, _ := json.Marshal(map[string]interface{}{
		"type": typ,
		"msg":  msg,
	})

	_, _ = client.Conn.Write(js)
}

func (client *Client) Pm(msg Message) {
	msg.Mid = rand.Intn(1000000)
	client.ChatMut.Lock()
	if len(client.Chat) >= 32 {
		client.Chat = client.Chat[1:]
	}
	client.Chat = append(client.Chat, msg)
	client.ChatMut.Unlock()
}

func (client *Client) ServerPm(msg string) {
	client.Pm(Message{Username: "[Server]", Body: msg, Id: -1})
}
