package diff

import (
	"fmt"
	"strings"
	"testing"
)

// -------------------------------------------
// ------------------------------------------- type TestCase
// -------------------------------------------

type TestCase interface {
	execute(tester *tTester)
}

// -------------------------------------------
// ------------------------------------------- type MasterTestCase
// -------------------------------------------

// This is a type of TestCase which is simply a container for other
// test cases.  It can be used to package up a group of lower-level
// test cases into one larger test case.  A master test case may
// also contain other master test cases, should you need to do that.

type MasterTestCase struct {
	title string
	testCases []TestCase
}

// Assert that TestCase is implemented by MasterTestCase.
var _ TestCase = (*MasterTestCase)(nil)

func (self *MasterTestCase) execute(tester *tTester) {
	for _, testCase := range self.testCases {
		testCase.execute(tester)
	}
}

func (self *MasterTestCase) appendSimpleTestCases(testCases []TestCase) {
	self.testCases = append(self.testCases, testCases...)
}

// -------------------------------------------
// ------------------------------------------- type tTester
// -------------------------------------------

type tTester struct {
	*testing.T
	name string
	testeeFn func (s, t string) int
}

// Assert that SimpleLogger is implemented by tTester.
var _ SimpleLogger = (*tTester)(nil)

// ------------------------------------------- NewTester tTester factory function

func NewTester(t *testing.T, name string, testeeFn func (s, t string) int) *tTester {
	return &tTester{t, name, testeeFn}
}

// ------------------------------------------- tTester Print

func (tester *tTester) Println(a ...interface{}) {
	fmt.Println(a...)
}

// ------------------------------------------- tTester Printf

func (tester *tTester) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// ------------------------------------------- tTester PrintBanner

func (tester *tTester) PrintBanner(title string) {
	bannerTitle := fmt.Sprintf("----- %s -------", title)
	bannerBar := strings.Repeat("-", len(bannerTitle))

	// log to stdout
	fmt.Println()
	fmt.Println(bannerBar)
	fmt.Println(bannerTitle)
	fmt.Println(bannerBar)
	fmt.Println()

	// log to the test log
	tester.Log("")
	tester.Log(bannerBar)
	tester.Log(bannerTitle)
	tester.Log(bannerBar)
	tester.Log("")
}

// ------------------------------------------- tTester PrintHeader

func (tester *tTester) PrintHeader(a ...interface{}) {

	// log to stdout
	fmt.Println()
	fmt.Println(a...)
	fmt.Println()

	// log to the test log
	tester.Log()
	tester.Log(a...)
	tester.Log()
}

// ------------------------------------------- tTester PrintSummary

func (tester *tTester) PrintSummary(a ...interface{}) {

	// log to stdout
	fmt.Println()
	fmt.Println(a...)
	fmt.Println()

	// log to the test log
	tester.Log()
	tester.Log(a...)
	tester.Log()
}
