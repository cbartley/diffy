package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"diffy/diff"
	"diffy/output"
)

// ------------------------------------------- main

func main() {

	// Do we have the right number of arguments?
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s FILE1 FILE2\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Exit 1.")
		os.Exit(1)
	}

	// Extract our arguments.
	pathToFile1, pathToFile2 := os.Args[1], os.Args[2]

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
	output.GenerateHtmlDiffPage(alignment, sourceLines1, sourceLines2)
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

