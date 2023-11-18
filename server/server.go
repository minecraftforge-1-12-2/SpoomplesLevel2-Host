package server

import (
	"TogetherForever/config"
	"TogetherForever/util"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
)

const VERSION = "1.2.4"

var MainServer Server

// --- Server Struct --- //

type Server struct {
	Up bool

	Host       string
	Port       int
	Timeout    int
	MaxPlayers int
	Anticheat  bool

	Clients map[int]*Client

	ClientsMut sync.Mutex

	Keys []string
	Bans []string
}

// --- Server Methods --- //

func InitServer(conf config.Config) {
	MainServer = Server{
		Up: false,

		Host:       conf.Host,
		Port:       conf.Port,
		Timeout:    conf.Timeout,
		MaxPlayers: conf.MaxPlayers,
		Anticheat:  conf.Anticheat,

		Clients: make(map[int]*Client),

		ClientsMut: sync.Mutex{},

		Keys: config.ParseLsf(conf.Keys, config.KEYS),
		Bans: config.ParseLsf(conf.Bans, config.BANS),
	}
}

func (server *Server) Start() {
	// Initialize Server
	util.WriteLine("  _____              _   _            \n |_   _|__  __ _ ___| |_| |_  ___ _ _ \n   | |/ _ \\/ _` / -_)  _| ' \\/ -_) '_|\n  _|_|\\___/\\__, \\___|\\__|_||_\\___|_|  \n | __|__ _ |___/__ _____ _ _          \n | _/ _ \\ '_/ -_) V / -_) '_|         \n |_|\\___/_| \\___|\\_/\\___|_|")
	util.WriteLine("Server starting on {0}:{1}...", server.Host, strconv.Itoa(server.Port))
	defer server.Stop()
	server.Up = true

	// Start TCP Listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))
	if err != nil {
		panic("Error opening server: " + err.Error())
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	util.WriteLine("Good to go!")

	for server.Up {
		conn, err := listener.Accept()
		if err != nil {
			util.WriteLine("Error accepting connection: {0}", err.Error())
			continue
		}

		// Generate Client
		client := Client{
			ID:    rand.Intn(9999),
			Conn:  conn,
			Ip256: util.MakeIp256(strings.Split(conn.RemoteAddr().String(), ":")[0]),
		}
		util.WriteLine("Client {0} connected from {1}", strconv.Itoa(client.ID), client.Ip256)
		go client.Accept()
	}
}

func (server *Server) Stop() {
	server.Up = false
	for _, client := range server.Clients {
		_ = client.Conn.Close()
	}
}

func (server *Server) Announce(msg string) {
	util.WriteLine(msg)
	server.ClientsMut.Lock()
	for _, client := range server.Clients {
		client.Queue = append(client.Queue, CompMessage{Msg: msg, Type: OmsgAnnouncement})
	}
	server.ClientsMut.Unlock()
}

func (server *Server) Kick(id int, reason string) {
	server.ClientsMut.Lock()
	if client, ok := server.Clients[id]; ok {
		server.ClientsMut.Unlock()
		if client.Admin {
			return
		}
		client.Close(OmsgKick, reason)
	} else {
		server.ClientsMut.Unlock()
	}
}

func (server *Server) Ban(id int, reason string) {
	server.ClientsMut.Lock()
	if client, ok := server.Clients[id]; ok {
		server.ClientsMut.Unlock()
		if client.Admin {
			return
		}
		server.Bans = append(server.Bans, server.Clients[id].Ip256)
		config.SaveLsf("bans.lsf", server.Bans)
		client.Close(OmsgKick, reason)
	} else {
		server.ClientsMut.Unlock()
	}
}

func (server *Server) CheckKey(key string) bool {
	for _, k := range server.Keys {
		if k == key {
			return true
		}
	}
	return false
}

func (server *Server) CheckBanned(ip256 string) bool {
	for _, b := range server.Bans {
		if b == ip256 {
			return true
		}
	}
	return false
}

func (server *Server) LobbyCnt(lobby string) int {
	server.ClientsMut.Lock()
	cnt := 0
	for _, client := range server.Clients {
		if client.Lobby == lobby {
			cnt++
		}
	}
	server.ClientsMut.Unlock()
	return cnt
}

// Broadcast to all clients in a lobby
// DO NOT USE. CURRENTLY BROKEN!
func (server *Server) Broadcast(msg string, lobby string) {
	server.ClientsMut.Lock()
	for _, client := range server.Clients {
		if client.Lobby == lobby && client.LoggedIn && client.Active {
			client.ServerPm(msg)
		}
	}
	server.ClientsMut.Unlock()
}
