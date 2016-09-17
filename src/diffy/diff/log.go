package diff

import (
	"fmt"
	"os"
)

// -------------------------------------------
// ------------------------------------------- type SimpleLogger
// -------------------------------------------

// Simple logger interface that diagnostic code can use to direct
// output to various destinations.  For instance we may want
// to output diagnostic text using testing.T's log methods when
// running unit tests, but not when running manual tests.  In
// the latter case a testing.T won't even be available.

type SimpleLogger interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
}

// -------------------------------------------
// ------------------------------------------- SimpleStdoutLogger global
// -------------------------------------------

// The global variable "SimpleStdoutLogger" points to an object which
// implements the "SimpleLogger" interface, and you can use it any time
// you want to direct SimpleLogger output to stdout.  You could create
// your own instance of type "tSimpleStdoutLogger", but you don't need to.

type tSimpleStdoutLogger struct {}

var SimpleStdoutLogger tSimpleStdoutLogger

func (s tSimpleStdoutLogger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (s tSimpleStdoutLogger) Println(a ...interface{}) {
	fmt.Println(a...)
}

// -------------------------------------------
// ------------------------------------------- SimpleStderrLogger global
// -------------------------------------------

type tSimpleStderrLogger struct {}

var SimpleStderrLogger tSimpleStderrLogger

func (s tSimpleStderrLogger) Printf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func (s tSimpleStderrLogger) Println(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}
