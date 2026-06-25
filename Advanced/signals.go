package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// setupSignalChannel creates and returns a buffered channel that listens
// for the specified OS signals.
//
// NOTE: Always use a buffered channel (capacity >= 1) for signal.Notify.
// The OS sends signals asynchronously. If the channel is unbuffered and
// your goroutine isn't ready to receive at that exact moment, the signal
// is dropped silently — your handler never runs.
// Example: signal.Notify(ch, syscall.SIGINT) with an unbuffered ch can
// miss Ctrl+C presses if the goroutine is busy.

// Signal support is platform dependent.
// Signals such as SIGUSR1 and SIGHUP are Unix-specific and are not
// available on Windows. Programs that rely on these signals should
// account for platform differences.
func setupSignalChannel() chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGINT,  // Ctrl+C from terminal
		syscall.SIGTERM, // kill <PID> or docker stop
		syscall.SIGHUP,  // terminal closed / config reload convention
		syscall.SIGUSR1, // user-defined: trigger a custom action without stopping
	)
	return sigCh
}

// handleSignals waits for a signal on sigCh, logs it, then sends true on
// done to tell the worker goroutine to stop.
//
// NOTE — Signal behavior when you call signal.Notify:
// Registering a signal with signal.Notify overrides its default OS behavior.
// Default behavior for SIGINT/SIGTERM is immediate process termination.
// Once registered, Go intercepts the signal and delivers it to your channel
// instead of killing the process. Your code is now responsible for exiting.
//
// Practical consequence: if you register SIGINT but never read from the
// channel (or never call os.Exit / return from main), pressing Ctrl+C will
// appear to do nothing. The process keeps running indefinitely.
//
// SIGUSR1 is intentionally handled with `continue` — it triggers a custom
// action but does NOT stop the program. This is the standard pattern for
// live config reloads, log rotation triggers, or debug dumps.
//
// Signals NOT listed in signal.Notify retain their default OS behavior.
// Example: SIGKILL and SIGSTOP can never be caught or overridden — the OS
// always handles them directly. `kill -9 <PID>` works even when all other
// signals are intercepted.

// The signal listener runs in its own goroutine because receiving from
// sigCh blocks until a signal arrives. Running it in a separate goroutine
// allows the rest of the application (workers, HTTP server, etc.) to
// continue running while waiting for signals asynchronously.
func handleSignals(sigCh chan os.Signal, done chan bool, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for sig := range sigCh {
			switch sig {
			case syscall.SIGINT:
				fmt.Println("Received SIGINT (Ctrl+C) — shutting down...")
				done <- true
				return

			case syscall.SIGTERM:
				fmt.Println("Received SIGTERM (kill / docker stop) — shutting down...")
				done <- true
				return

			case syscall.SIGHUP:
				// SIGHUP originally meant "terminal hang-up" (modem era).
				// Modern convention: reload configuration without restarting.
				// Nginx, for example, does `kill -HUP <pid>` to reload nginx.conf.
				fmt.Println("Received SIGHUP — reloading configuration (no shutdown)...")
				// In a real app you would call something like reloadConfig() here.
				// We do NOT send to done; the worker keeps running.

			case syscall.SIGUSR1:
				// SIGUSR1/SIGUSR2 have no fixed meaning — you define them.
				// Common uses: dump goroutine stack traces, toggle debug logging,
				// flush in-memory buffers to disk.
				fmt.Println("Received SIGUSR1 — executing user-defined action...")
				// We do NOT send to done; the worker keeps running.
			}
		}
	}()
}

// doWork simulates a long-running task (e.g., processing a job queue).
// It checks the done channel on every iteration and exits cleanly when signaled.
//
// NOTE — Why `select` with a `default` case instead of blocking on <-done:
// A blocking receive (`<-done`) would pause the goroutine at that line and
// prevent it from doing any work. Using `select` with `default` makes the
// done-check non-blocking: if nothing is in done, execution falls through to
// `default` immediately and the work continues.
// This is the standard pattern for a "cancellable worker loop".
func Dowork(done chan bool, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Println("Worker received stop signal — cleaning up and exiting...")
				return
			default:
				fmt.Println("Worker: processing...")
			}
			// Simulate work taking 1 second per unit.
			// In production this could be an HTTP request, DB query, file write, etc.
			time.Sleep(time.Second)
		}
	}()
}

func Signals() {
	fmt.Println("Process ID:", os.Getpid())
	// TIP: note the PID printed above. You can send signals manually from
	// another terminal to test each case:
	//   kill -SIGINT  <PID>   → same as Ctrl+C
	//   kill -SIGTERM <PID>   → graceful shutdown request
	//   kill -SIGHUP  <PID>   → reload config (no shutdown)
	//   kill -SIGUSR1 <PID>   → user-defined action (no shutdown)
	//   kill -9       <PID>   → SIGKILL: cannot be caught, always kills immediately

	var wg sync.WaitGroup

	// done is the coordination channel between the signal handler and the worker.
	// Capacity 1 ensures the sender (signal handler) never blocks even if the
	// worker hasn't read yet.
	done := make(chan bool, 1)

	sigCh := setupSignalChannel()

	handleSignals(sigCh, done, &wg)
	Dowork(done, &wg)

	fmt.Println("Application running. Send a signal to stop.")

	// wg.Wait() blocks main until both goroutines call wg.Done().
	// Without this, main() would return immediately after starting the
	// goroutines and the process would exit before any signal is received.
	wg.Wait()

	// signal.Stop tells the runtime to stop delivering signals to sigCh.
	// Close sigCh so the range loop inside handleSignals can exit cleanly
	// if it is still running (e.g., after a SIGHUP/SIGUSR1 that didn't stop the worker).

	// Always call signal.Stop before closing the channel.
	// Otherwise the runtime may still attempt to deliver signals
	// to a closed channel, causing a panic.
	signal.Stop(sigCh)
	close(sigCh)

	fmt.Println("Application shut down gracefully.")
}

// OS signals are not queued indefinitely.
// If multiple identical signals arrive before your program reads them,
// the operating system (and Go runtime) may coalesce them into a single
// notification. You should never rely on receiving every occurrence.
