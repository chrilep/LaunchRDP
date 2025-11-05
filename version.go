package main

import "github.com/chrilep/LaunchRDP/logging"

const (
	AppName    = "LaunchRDP"
	ID         = "com.chrilep.launchrdp"
	Version    = "1.3.0.7"
	Author     = "Lancer"
	Repository = "https://github.com/chrilep/LaunchRDP"
)

// PrintVersion prints version information
func PrintVersion() {
	logging.Log(true, AppName+" v"+Version)
	logging.Log(true, "Author: "+Author)
	logging.Log(true, "Repository: "+Repository)
}
