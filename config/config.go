package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"strings"
)

func LoadConfig(name string) Config {
	name = "conf/" + name
	// Create directory if it doesn't exist.
	if _, err := os.Stat("conf"); os.IsNotExist(err) {
		err := os.Mkdir("conf", 0755)
		if err != nil {
			panic(err)
		}
	}

	// Create file if it doesn't exist.
	if _, err := os.Stat(name); os.IsNotExist(err) {
		_, err := os.Create(name)
		err2 := os.WriteFile(name, []byte(CONFIG), 0755)
		if err != nil || err2 != nil {
			panic(err)
		}
	}

	var cfg Config
	_, err := toml.DecodeFile(name, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func ParseLsf(name string, def string) []string {
	name = "conf/" + name

	// Create directory if it doesn't exist.
	if _, err := os.Stat("conf"); os.IsNotExist(err) {
		err := os.Mkdir("conf", 0755)
		if err != nil {
			panic(err)
		}
	}

	// Create file if it doesn't exist.
	if _, err := os.Stat(name); os.IsNotExist(err) {
		_, err := os.Create(name)
		err2 := os.WriteFile(name, []byte(def), 0755)
		if err != nil || err2 != nil {
			panic(err)
		}
	}

	// Open file for reading.
	file, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	// Read file line by line, ignoring comments and empty lines.
	var lines []string
	for _, line := range strings.Split(string(file), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		lines = append(lines, line)
	}

	return lines
}

func SaveLsf(name string, lines []string) {
	name = "conf/" + name

	// Create directory if it doesn't exist.
	if _, err := os.Stat("conf"); os.IsNotExist(err) {
		err := os.Mkdir("conf", 0755)
		if err != nil {
			panic(err)
		}
	}

	// Create file if it doesn't exist.
	if _, err := os.Stat(name); os.IsNotExist(err) {
		_, err := os.Create(name)
		if err != nil {
			panic(err)
		}
	}

	// Open file for writing.
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	// Write lines to file.
	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}
}
