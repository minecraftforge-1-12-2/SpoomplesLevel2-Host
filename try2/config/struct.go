package config

type Config struct {
	// [server]
	Host       string
	Port       int
	Timeout    int
	MaxPlayers int

	// [security]
	Anticheat bool

	// [files]
	Keys     string
	Bans     string
	BadWords string

	// [gui]
	Enable bool
	Log    bool
	Input  bool
}
