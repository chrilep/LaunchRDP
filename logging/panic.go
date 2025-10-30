/*
Panic handler/logger:
	- It catches any panic that occurs in the application and logs the reason and stack trace to the log file.
	- It can be used by calling PanicHandler() at the start of the main function and every goroutine that might panic.

	Usage:
	defer PanicHandler()
*/

package logging

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

// PanicHandler catches any panic that occurs in the application.
// It logs the panic reason and stack trace to the log file.
func PanicHandler() {
	if r := recover(); r != nil {
		if fileLog != nil {
			// Write to log file if it is ready
			Log(true, "HEARTBEAT: CRITICAL - Application panic detected!")
			Log(true, "==== PANIC ====")
			Log(true, fmt.Sprintf("Time  : %s", time.Now().Format("2006-01-02 15:04:05")))
			Log(true, fmt.Sprintf("Reason: %v", r))
			Log(true, string(debug.Stack()))
			Log(true, "==== END PANIC ====")
		} else {
			// Append to the log file if it is not ready
			f, err := os.OpenFile(strLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return
			}
			defer f.Close()
			fmt.Fprintln(f, "HEARTBEAT: CRITICAL - Application panic detected!")
			fmt.Fprintln(f, "==== PANIC ====")
			fmt.Fprintf(f, "Time  : %s\n", time.Now().Format("2006-01-02 15:04:05"))
			fmt.Fprintln(f, "Reason:", r)
			fmt.Fprintln(f, string(debug.Stack()))
			fmt.Fprintln(f, "==== END PANIC ====")
			fmt.Fprintln(f, "HEARTBEAT: Application may terminate due to panic")
		}
		// Re-panic to let the application handle it properly
		panic(r)
	}
}

// SafeCallback wraps a function with panic handling
func SafeCallback(function func()) func() {
	return func() {
		defer PanicHandler()
		function()
	}
}
