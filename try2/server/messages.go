package server

const (
	MsgNone = iota

	ImsgLogin
	ImsgDefault
	ImsgPaused
	ImsgMessage

	OmsgDefault
	OmsgAccepted
	OmsgDisconnect
	OmsgKick
	OmsgAnnouncement
)

type CompMessage struct {
	Type int
	Msg  string
}

// CompClient is a compacted version of a Client used for sending to clients.
type CompClient struct {
	ID    int     `json:"id"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Name  string  `json:"name"`
	Admin bool    `json:"admin"`
	Room  uint16  `json:"room"`

	Sprite         string `json:"sprite"`
	Frame          uint8  `json:"frame"`
	Direction      int    `json:"dir"`
	Palette        uint8  `json:"palette"`
	PaletteSprite  string `json:"paletteSprite"`
	PaletteTexture string `json:"paletteTexture"`
	Color          string `json:"color"`
}

// Message is NOT a Message for a Client. It is a Text Message!
type Message struct {
	Body     string `json:"body"`
	Username string `json:"username"`
	Id       int    `json:"id"`
	Mid      int    `json:"mid"`
}
