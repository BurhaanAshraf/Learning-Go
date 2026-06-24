// launch and control external system commands — the same commands you would
// type in a terminal.
//
// Core idea: exec.Command() creates a Cmd struct that describes what you want
// to run. Nothing actually executes until you call one of the Run methods:
//
//	cmd.Output()        — runs the command, waits, returns stdout as []byte
//	cmd.Run()           — runs the command, waits, returns only an error
//	cmd.Start()         — starts the command WITHOUT waiting (async)
//	cmd.CombinedOutput() — like Output() but also captures stderr

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// 1. BASIC COMMAND — exec.Command + Output()
// ---------------------------------------------------------------------------

// simpleTask shows the minimal pattern: build a command, run it, read output.
//
// exec.Command(name, arg1, arg2, ...) — the first argument is the program
// name (what you would type in the terminal), every argument after that is
// passed to that program separately. Do NOT combine them into one string:
//
//	WRONG:  exec.Command("echo Hello World")    // tries to find a binary literally named "echo Hello World"
//	RIGHT:  exec.Command("echo", "Hello World") // runs "echo" with argument "Hello World"
//
// Output() blocks until the command finishes and returns stdout as a []byte.
// Convert it with string(output) to get a readable string.
func simpleTask() {
	cmd := exec.Command("echo", "Hello World")

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Output: %s", string(output)) // output already has a trailing newline from echo
}

// ---------------------------------------------------------------------------
// 2. PIPING STDIN TO A COMMAND — cmd.Stdin
// ---------------------------------------------------------------------------

// grep shows how to feed input INTO a command via stdin, just like the pipe
// operator in a shell: echo "some text" | grep foo
//
// By default, cmd.Stdin is nil (no input). You can point it at anything that
// satisfies io.Reader — here we use strings.NewReader() to turn a Go string
// into a reader without needing a real file or pipe.
//
// EXIT CODES MATTER:
// grep returns exit code 0 if it found a match, exit code 1 if it found
// nothing (not an error — just "no results"). cmd.Output() converts any
// non-zero exit code into a *exec.ExitError, so you have to distinguish
// "no match" (code 1) from "command failed" (code 2+ or not found).
func grep() {
	cmd := exec.Command("grep", "foo")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text to grep for 'foo': ")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// TrimSpace removes the trailing newline from ReadString — without this,
	// the input passed to grep would include the newline character and could
	// produce unexpected match behaviour.
	input = strings.TrimSpace(input)

	// Point the command's stdin at our string instead of the terminal.
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.Output()
	if err != nil {
		// Type-assert to *exec.ExitError to inspect the exit code.
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				fmt.Println("Result: No match found.")
				return
			}
		}
		// Any other error means the command itself couldn't run.
		fmt.Println("Error executing command:", err)
		return
	}

	fmt.Printf("Result: Match found!\n%s\n", string(output))
}

// ---------------------------------------------------------------------------
// 3. ASYNC EXECUTION — cmd.Start() vs cmd.Wait()
// ---------------------------------------------------------------------------

// timeConsuming shows the difference between Start() and Wait(), and how to
// kill a process that's been running too long.
//
// cmd.Output() / cmd.Run() = synchronous — your Go code STOPS at that line
// until the command finishes. Fine for fast commands; bad for long ones.
//
// cmd.Start() = asynchronous — launches the command in the background and
// returns IMMEDIATELY. Your Go code keeps running. The command runs
// concurrently alongside your program.
//
// cmd.Wait() = blocks until the background command finishes. You must call
// Wait() after Start() to release the process resources (even if you killed
// it), otherwise you get a zombie process.
//
// cmd.Process.Kill() = sends SIGKILL to the process immediately. The killed
// process still needs to be Wait()-ed on to clean up — hence calling Wait()
// after Kill() even though we know it will return an error.
func timeConsuming() {
	cmd := exec.Command("sleep", "50") // would block for 50s if we used cmd.Run()

	// Start() returns as soon as the OS has launched the process.
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting process:", err)
		return
	}

	fmt.Println("Process started, doing other work for 2 seconds...")
	time.Sleep(2 * time.Second) // simulates your program doing other things

	// Kill sends SIGKILL — immediate, no cleanup on the child's side.
	err = cmd.Process.Kill()
	if err != nil {
		fmt.Println("Error killing process:", err)
		return
	}

	// Always Wait() after Kill() to release OS resources for this process.
	// The error here will be "signal: killed" — that's expected, ignore it.
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error killing process:", err)
		return
	}

	fmt.Println("Process killed.")
}

// ---------------------------------------------------------------------------
// 4. READING ENVIRONMENT VARIABLES FROM A COMMAND
// ---------------------------------------------------------------------------

// envVars shows that commands launched by exec inherit the parent process's
// environment by default (same PATH, HOME, SHELL, etc.).
//
// If you need to run a command with a DIFFERENT environment, set cmd.Env
// before calling Output()/Run():
//
//	cmd.Env = append(os.Environ(), "MY_VAR=hello")
//
// Passing os.Environ() as the base keeps existing variables intact;
// leaving it out starts with an empty environment.
func envVars() {
	cmd := exec.Command("printenv", "SHELL") // prints the value of $SHELL

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Print(string(output))
}

// ---------------------------------------------------------------------------
// 5. IO PIPE — streaming data between Go and a command
// ---------------------------------------------------------------------------

// IOPipe shows a more powerful stdin pattern using io.Pipe(), which creates
// a synchronised in-memory pipe:
//
//	pr, pw := io.Pipe()
//	pr = the read end  — hand this to cmd.Stdin so the command reads from it
//	pw = the write end — your Go code writes into this
//
// WHY use io.Pipe() instead of strings.NewReader()?
// strings.NewReader() requires the entire input to already be in memory as a
// string. io.Pipe() lets you STREAM data — you can write chunks into pw from
// a goroutine while the command is simultaneously reading from pr. This is
// essential when the data is large or generated on the fly.
//
// IMPORTANT: always close pw when you're done writing. The command (grep
// here) reads until EOF — closing pw is what signals EOF. Without close(),
// the command hangs forever waiting for more input.
func IOPipe() {
	pr, pw := io.Pipe()

	cmd := exec.Command("grep", "foo")
	cmd.Stdin = pr // command reads from the read-end of the pipe

	// Write into the pipe from a goroutine so it runs concurrently with
	// cmd.Output() below. If we wrote synchronously here, we'd deadlock:
	// pw.Write blocks until the reader consumes data, but cmd.Output()
	// hasn't started yet — neither side would make progress.
	go func() {
		defer pw.Close() // signals EOF to grep when the goroutine exits

		data := []byte("burger\nfood is good\npizza\nfoo faa ji\n")
		pw.Write(data)
	}()

	// Output() starts the command and reads its stdout. Internally, grep is
	// now pulling lines from pr while the goroutine above pushes into pw.
	output, err := cmd.Output()
	if err != nil {
		// grep exits 1 for no match — handle same as in grep() above.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			fmt.Println("No match found.")
			return
		}
		fmt.Println("Error:", err)
		return
	}

	fmt.Print(string(output)) // prints lines containing "foo"
}

// ---------------------------------------------------------------------------
// 6. CAPTURING STDERR TOO — CombinedOutput()
// ---------------------------------------------------------------------------

// ls shows CombinedOutput(), which merges both stdout and stderr into one
// []byte. This is useful when:
//   - You want to log everything the command printed (including errors).
//   - The command writes useful information to stderr on failure and you
//     don't want to miss it.
//
// If you need stdout and stderr SEPARATELY, use cmd.StdoutPipe() and
// cmd.StderrPipe() instead — they give you independent io.Reader handles.
func ls() {
	cmd := exec.Command("ls", "-l")

	// CombinedOutput() = stdout + stderr interleaved in arrival order.
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Even on failure, output may contain useful stderr text.
		fmt.Println("Error:", err)
		fmt.Println("Command output:", string(output))
		return
	}

	fmt.Println(string(output))
}

func Exec() {
	fmt.Println("\n===== 1. BASIC COMMAND =====")
	simpleTask()

	fmt.Println("\n===== 2. STDIN TO COMMAND =====")
	grep()

	fmt.Println("\n===== 3. ASYNC EXECUTION =====")
	timeConsuming()

	fmt.Println("\n===== 4. ENVIRONMENT VARIABLES =====")
	envVars()

	fmt.Println("\n===== 5. IO PIPE =====")
	IOPipe()

	fmt.Println("\n===== 6. COMBINED OUTPUT =====")
	ls()
}

// ---------------------------------------------------------------------------
// QUICK REFERENCE
// ---------------------------------------------------------------------------
//
// Method            Waits?   Captures output?   Use when
// ---------------------------------------------------------------------------
// cmd.Run()         YES      NO (just error)    You only care if it succeeded
// cmd.Output()      YES      stdout only        You need the command's output
// cmd.CombinedOutput() YES   stdout + stderr    You need everything it printed
// cmd.Start()       NO       —                  Long-running / async work
// cmd.Wait()        YES      —                  After Start(), to clean up
// cmd.Process.Kill() NO      —                  Abort a Start()-ed process
//
// stdin patterns:
//   cmd.Stdin = strings.NewReader(s)  — small, known-upfront string
//   pr, pw := io.Pipe(); cmd.Stdin = pr  — large/streaming data from goroutine
//
// environment:
//   cmd.Env = append(os.Environ(), "KEY=value")  — add/override one var
//   cmd.Env = []string{"KEY=value"}              — isolated environment
