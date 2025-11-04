package main

import "fmt"

const (
	AppName    = "LaunchRDP"
	ID         = "com.chrilep.launchrdp"
	Version    = "1.2.0.7"
	Author     = "Lancer"
	Repository = "https://github.com/chrilep/LaunchRDP"
)

// PrintVersion prints version information
func PrintVersion() {
	fmt.Printf("%s v%s\n", AppName, Version)
	fmt.Printf("Author: %s\n", Author)
	fmt.Printf("Repository: %s\n", Repository)
}
