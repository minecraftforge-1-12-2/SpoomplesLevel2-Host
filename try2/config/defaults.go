package config

const (
	CONFIG = `# Together Server Config

# Public Host
host = "0.0.0.0"
# Default Port
port = 6666
# Timeout in seconds
timeout = 10
# Max players
maxplayers = 999


# Enable Anticheat
anticheat = true


# Name the files (no directory)
# Keys file
keys = "keys.lsf"
# Bans file
bans = "bans.lsf"
# Bad words file
badwords = "badwords.lsf"


# Enable GUI
enable = true
# Enable GUI logging
log = true
# Enable GUI input
input = true
`
	KEYS = `# Together Server Admin Keys File
# This file contains a list of admin keys.
# Each key should be on a new line.
# You may use the keygen tool to generate keys! But any valid string will work.
# Hashtags are comments!

# Example:
# 1234567890
`
	BANS = `# Together Server Bans File
# This file contains a list of banned Adr256s.
# Each Adr256 should be on a new line.
# You're probably not gonna want to edit this file manually. Use the /ban command instead.
`
	BADWORDS = `# Together Server Bad Words File
# This file contains a list of bad words.
# Each word should be on a new line.

# Example:
fuck
bitch
ass
whore
cunt
nigger
nigga
niga
niger
tranny
trany
tranni
trani
faggot
fagot
fag
faag
`
)
