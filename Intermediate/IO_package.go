package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// ====================================================
// io.Reader
// ====================================================
// WHAT    : Standard interface for reading data
// SIG     : Read(p []byte) (n int, err error)
// WHO     : Files, Strings, Buffers, Network Conns, HTTP Bodies
// WHY     : Accept io.Reader → works with ANY data source
// ====================================================

func readFromReader(r io.Reader) {

	buf := make([]byte, 1024)

	n, err := r.Read(buf)

	if err != nil {
		log.Fatal("Cannot Read From Reader:", err)
	}

	// buf[:n] → only the bytes actually read (not the full 1024)
	fmt.Println(buf[:n])
}

// ====================================================
// io.Writer
// ====================================================
// WHAT    : Standard interface for writing data
// SIG     : Write(p []byte) (n int, err error)
// WHO     : Files, Buffers, HTTP Responses, Network Conns
// WHY     : Accept io.Writer → write anywhere without caring where
// ====================================================

func writeFromWriter(w io.Writer, data string) {

	_, err := w.Write([]byte(data))

	if err != nil {
		log.Fatal("Cannot Write:", err)
	}
}

// ====================================================
// io.Closer
// ====================================================
// WHAT    : Interface for releasing resources
// SIG     : Close() error
// WHO     : Files, DB Connections, HTTP Response Bodies
// WHY     : Any type with Close() satisfies io.Closer automatically
// ====================================================

func closeResource(c io.Closer) {

	if err := c.Close(); err != nil {
		log.Fatal("Cannot Close Resource:", err)
	}
}

// ====================================================
// bytes.Buffer
// ====================================================
// WHAT    : In-memory byte storage (like a RAM-based temp file)
// USE     : Building strings, JSON payloads, HTTP requests
// IMPLS   : io.Reader + io.Writer
// ALLOC   : var buf bytes.Buffer → stack allocation
// ====================================================

func bufferExample() {

	// Stack allocation.
	var buf bytes.Buffer

	buf.WriteString("Hello Buffer")

	fmt.Println(buf.String())
}

// ====================================================
// io.MultiReader
// ====================================================
// WHAT    : Combines multiple readers into one sequential reader
// ORDER   : r1 → r2 → r3 → ...
// USE     : Data arrives from multiple sources, read as one stream
// ====================================================

func multiReaderExample() {

	r1 := strings.NewReader("Hello ")
	r2 := strings.NewReader("World")

	mr := io.MultiReader(r1, r2)

	// Heap allocation.
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(mr)

	if err != nil {
		log.Fatal("Cannot Read:", err)
	}

	fmt.Println(buf.String())
}

// ====================================================
// io.Pipe
// ====================================================
// WHAT    : Creates a connected Writer → Reader pair
// FLOW    : writer.Write() → reader.Read()
// USE     : Streaming, goroutine communication, pipelines
// KEY     : Synchronous — writer BLOCKS until reader consumes
// CLOSE   : writer.Close() signals EOF to reader
// ====================================================

func pipeExample() {

	reader, writer := io.Pipe()

	go func() {

		writer.Write([]byte("Hello Pipe"))

		// Closing signals EOF to reader.
		writer.Close()
	}()

	buf := new(bytes.Buffer)

	buf.ReadFrom(reader)

	fmt.Println(buf.String())
}

// ====================================================
// FILE WRITING — os.OpenFile
// ====================================================
// FLAGS:
//   O_APPEND → append to existing content
//   O_CREATE → create file if it doesn't exist
//   O_WRONLY → write-only access
//
// PERMISSIONS (0644):
//   Owner  → read + write
//   Others → read only
//
// NOTE: *os.File implements io.Writer
//   var writer io.Writer = file  ← valid assignment
// ====================================================

func writeToFile(filepath string, data string) {

	file, err := os.OpenFile(
		filepath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		log.Fatal("Cannot Open File:", err)
	}

	// Ensures cleanup before function exits.
	defer closeResource(file)

	_, err = file.Write([]byte(data))

	if err != nil {
		log.Fatal("Cannot Write To File:", err)
	}

// *os.File implements io.Writer,
// so it can be assigned to an io.Writer variable.

// var writer io.Writer = file
// _, err = writer.Write([]byte(data))
// if err != nil {
// 	log.Fatal("Cannot Write To File", err)
// }

}

// ====================================================
// INTERFACE SATISFACTION (Implicit)
// ====================================================
// RULE    : Implement all required methods → automatically satisfies interface
// NO      : No "implements" keyword needed
// EXAMPLE : myResource has Close() error → satisfies io.Closer
// ====================================================

type myResource struct {
	name string
	time time.Time
}

func (m myResource) Close() error {

	fmt.Println(
		"Closing",
		m.name,
		"at",
		m.time,
	)

	return nil
}

func IO() {

	fmt.Println("==== Reader Example ====")

	// strings.NewReader converts a string
	// into an io.Reader.
	readFromReader(
		strings.NewReader("Hello Reader!"),
	)

	// ------------------------------------------------

	fmt.Println("==== Writer Example ====")

	var writer bytes.Buffer

	// bytes.Buffer implements io.Writer.
	writeFromWriter(
		&writer,
		"Hello Writer",
	)

	fmt.Println(writer.String())

	// ------------------------------------------------

	fmt.Println("==== Buffer Example ====")

	bufferExample()

	// ------------------------------------------------

	fmt.Println("==== MultiReader Example ====")

	multiReaderExample()

	// ------------------------------------------------

	fmt.Println("==== Pipe Example ====")

	pipeExample()

	// ------------------------------------------------

	fmt.Println("==== File Example ====")

	writeToFile(
		"io.txt",
		"Appending New Line\n",
	)

	// ------------------------------------------------

	fmt.Println("==== Closer Example ====")

	resource := &myResource{
		name: "Test Resource",
		time: time.Now(),
	}

	// myResource satisfies io.Closer.
	closeResource(resource)
}

/*
====================================================
QUICK REVISION — io Package
====================================================

INTERFACE       METHOD SIG                  PURPOSE
-----------     --------------------------  -------------------
io.Reader       Read(p []byte) (n, err)     Pull data in
io.Writer       Write(p []byte) (n, err)    Push data out
io.Closer       Close() error               Release resource

----------------------------------------------------
TYPE / FUNC         KEY FACTS
----------------------------------------------------
bytes.Buffer        In-memory R+W | stack: var buf bytes.Buffer
                                  | heap:  new(bytes.Buffer)

strings.NewReader   string → io.Reader

io.MultiReader      Merge readers: r1 → r2 → r3

io.Pipe             Writer → Reader | synchronous | close = EOF

os.OpenFile         Flags: O_APPEND | O_CREATE | O_WRONLY
                    Perms: 0644 (owner rw, others r)

defer               Always used for cleanup: defer file.Close()

Interface Rule      Implement methods → satisfies interface (no keyword)

----------------------------------------------------
INTERVIEW QUESTIONS
----------------------------------------------------
Q1. What is io.Reader? What is its method signature?
Q2. What is io.Writer? What is its method signature?
Q3. What interfaces does bytes.Buffer implement?
Q4. Why use io.Pipe? What makes it synchronous?
Q5. Why use defer with files?
Q6. How does interface satisfaction work in Go?
Q7. Difference between io.Reader and io.Writer?
Q8. What do O_APPEND, O_CREATE, O_WRONLY do?

====================================================
*/