package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"diffy/diff"
	"diffy/etc"
	"diffy/output"
)

// ------------------------------------------- flags

var openWithPtr = flag.String("open-with", "", "open with")

// ------------------------------------------- main

func main() {

	// We must parse the flags before we do anything else.
	flag.Parse()

	// Do we have the right number of arguments?
	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s FILE1 FILE2\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Exit 1.")
		os.Exit(1)
	}

	// Extract our arguments.
	pathToFile1, pathToFile2 := flag.Arg(0), flag.Arg(1)

	// Do the specified files exist?
	if !checkThatPathExists(pathToFile1) || !checkThatPathExists(pathToFile2) {
		exitWithNotification(1)
	}

	// Are the files actually files?
	if !checkThatPathIsAFile(pathToFile1) || !checkThatPathIsAFile(pathToFile2) {
		exitWithNotification(1)
	}

	// Try to read the files.
	lines1, err := readFile(pathToFile1)
	if err != nil {
		exitWithNotification(2)
	}
	lines2, err := readFile(pathToFile2)
	if err != nil {
		exitWithNotification(3)
	}

	_, alignment := diff.Diff_v2(lines1, lines2)
	// alignment.Dump(lines1, lines2, 0, diff.SimpleStderrLogger)

	sourceLines1 := output.NewSourceLinesRec(lines1, pathToFile1)
	sourceLines2 := output.NewSourceLinesRec(lines2, pathToFile2)

	// We will output to stdout or a temporary file, depending.
	outputFile := os.Stdout
	if *openWithPtr != "" {
		outputFile, err = ioutil.TempFile("", "diffy")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open the temporary file; error = %v\n", err)
			exitWithNotification(4)
		}
		defer outputFile.Close()
	}

	output.GenerateHtmlDiffPage(outputFile, alignment, sourceLines1, sourceLines2)

	// If we are doing "--open-with" then we need to invoke the open command on the temp file.
	if *openWithPtr != "" {
		err := executeCommand(*openWithPtr, outputFile.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, 
						"Tried to execute the %q command %q, but got an error.\n", 
						"--open-with", *openWithPtr)
			fmt.Fprintf(os.Stderr, "The error was %v", err)
			exitWithNotification(4)
		}
	}
}

// ------------------------------------------- executeCommand

func executeCommand(cmdText string, extraArgs ...string) error {

	// Figure out the executable name and assemble the arguments.
	cmdWords := etc.ParseWords(cmdText)
	cmdName := cmdWords[0]
	cmdArgs := cmdWords[1:]
	cmdArgs = append(cmdArgs, extraArgs...)

	// Create the command and run it.
	command := exec.Command(cmdName, cmdArgs...)
	return command.Run();
}

// ------------------------------------------- checkThatPathExists

func checkThatPathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		fmt.Fprintf(os.Stderr, "The path %q does not exist.\n", path)
		fmt.Fprintln(os.Stderr)
		return false
	}
	return true
}

// ------------------------------------------- checkThatPathIsAFile

func checkThatPathIsAFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't stat the path %q.\n", path)
		fmt.Fprintln(os.Stderr)
		return false
	}
	if fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "The path %q points to a directory, not a file.\n", path)
		fmt.Fprintln(os.Stderr)
		return false
	}
	return true
}

// ------------------------------------------- readFile

func readFile(pathToFile string) (diff.ComparableLines, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var lines diff.ComparableLines
	for {
		strLine, err := reader.ReadString('\n')
		if len(strLine) > 0 {
			strLine = expandTabsAndStripLineEndings(strLine, 4)
			lines = append(lines, diff.NewTextLine(strLine))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return lines, nil
}

// ------------------------------------------- expandTabsAndStripLineEndings

func expandTabsAndStripLineEndings(s string, tabSize int) string {
	result := ""
	for _, char := range s {
		if char == '\t' {
			spaceCount := tabSize - len(result) % tabSize
			for i := 0; i < spaceCount; i++ {
				result += " "
			}
		} else if char == '\n' || char == '\r' {
			// do nothing
		} else {
			result += string(char)
		}
	}
	return result
}

// ------------------------------------------- exitWithNotification

func exitWithNotification(exitCode int) {
	fmt.Fprintf(os.Stderr, "Exit %d.\n", exitCode)
	os.Exit(exitCode)
}

