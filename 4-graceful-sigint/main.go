//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	proc := MockProcess{}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		<-signalCh
		close(done)

		<-signalCh
		fmt.Print("\nForcing process to stop.")
		os.Exit(0)
	}()

	go func() {
		proc.Run() // Run the process (blocking)
	}()

	<-done
	proc.Stop()
}
