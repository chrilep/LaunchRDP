/*
Lancer's simple logging module
	This module provides a simple logging mechanism that writes debug messages to both the console and a log file.
	The log file is created in the user's application data directory, and it is purged at the start of the application.
	It is designed to be used for debugging purposes, allowing developers to track the flow of the application.

	It needs three global variables to work:
	- strPublisherName: The name of the publisher, e.g. "Lancer"
	- strProductName: The name of the product, e.g. "LaunchRDP"
	- strVersion: The version of the product, e.g. "1.2.3.4"

	Usage:
	Log(true, "Some var", "is", var)
*/

package logging

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var strLogFilePath string // eg. <dataFolder>\Lancer\<Product>\log.txt
var fileLog *os.File
var strAppTempDir string // like %LOCALAPPDATA%\Lancer\<Product>\

// Global variables for application info
var strPublisherName = "Lancer"
var strProductName = "LaunchRDP"
var strVersion = "1.0.0"

// Log writes a message to the log file and console.
// If debug is false, it does nothing. If debug is true, it writes the message to the log file and console.
// It can take multiple arguments, which will be converted to strings.
// The first argument is the debug flag, the rest are the message parts.
// It also includes the name of the function that called it.
func Log(debug bool, arrMessageParts ...any) {
	if !debug {
		return
	}
	// check if the logfile is ready
	if fileLog == nil {
		err := activateLogging()
		if err != nil {
			fmt.Println("Error activating logging:", err)
			return
		}
	}
	// Get current time and format it as HH:mm:ss.fff
	now := time.Now()
	timestamp := now.Format("15:04:05.000")
	strParentName := `main.unknown`
	// Get the parent function's name
	ptrCaller, _, _, isSuccess := runtime.Caller(1)
	if isSuccess {
		funcCaller := runtime.FuncForPC(ptrCaller)
		if funcCaller != nil {
			strParentName = funcCaller.Name()
		}
	}
	// Convert inputs to strings
	arrMessages := make([]string, len(arrMessageParts))
	for i, v := range arrMessageParts {
		arrMessages[i] = fmt.Sprint(v)
	}
	fmt.Println(timestamp, `[`+strParentName+`]`, strings.Join(arrMessages, " "))
	fmt.Fprintln(fileLog, timestamp, `[`+strParentName+`]`, strings.Join(arrMessages, " "))
}

// activateLogging activates the logging module. See function Log() for details.
func activateLogging() error {
	// Can have no logging since logger not ready yet!
	debug := true // set to true to enable debug logging
	switch runtime.GOOS {
	case `windows`:
		strTempDir := os.Getenv("LOCALAPPDATA")
		if strTempDir == `` {
			strTempDir = os.Getenv("TMP")
			if strTempDir == `` {
				strTempDir = os.Getenv("TEMP")
			}
		}
		strAppTempDir = strTempDir + `\` + strPublisherName + `\` + strProductName
		strLogFilePath = strAppTempDir + `\log.txt`
	case `linux`:
		strHomeDir, err := os.UserHomeDir()
		if err != nil {
			strAppTempDir = ``
		} else {
			strAppTempDir = strHomeDir + `/.local/` + strProductName
			strLogFilePath = strAppTempDir + `/log.txt`
		}
	default:
		strAppTempDir = ``
	}
	// Check if directory exists.
	if _, err := os.Stat(strAppTempDir); os.IsNotExist(err) {
		// If not, create the directory.
		os.MkdirAll(strAppTempDir, 0755)
	}
	var err error
	fileLog, err = os.OpenFile(strLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Unable to open log file at '"+strLogFilePath+"':", err)
		return err
	}
	Log(debug, strProductName, strVersion, `is logging to `+strLogFilePath)
	return nil
}

// GetLogFilePath returns the current log file path
func GetLogFilePath() string {
	return strLogFilePath
}

// CloseLog closes the log file
func CloseLog() {
	if fileLog != nil {
		fileLog.Close()
		fileLog = nil
	}
}
