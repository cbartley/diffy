package diff

import (
	"math"
	"testing"
)

// ------------------------------------------- constants

const CHAR_SET = "abcd"
const MAX_INITIAL_STRING_LEN = 3
const MAX_OP_COUNT = 4

// -------------------------------------------
// ------------------------------------------- TestDiffHash
// -------------------------------------------

func TestDiffHash(tt *testing.T) {

	tester := NewTester(tt, "Testing DiffHash", nil)
	tester.PrintBanner("Testing DiffHash")

	// The hashing algorithm that we're using is basically designed around a fixed
	// size window (four runes, currently).  Of course it has to support strings
	// shorter than this window size as well, and this is obviously going to
	// require some sort of special handling for these cases.  So we want to test
	// a full set of strings shorter than the window size and, in fact, all the
	// way down to the empty string.

	// WARNING: A hash algorithm, by it's very nature, may occasionally compute the 
	// same hash for different strings.  Our test strings here should be chosen so
	// that this does not occur.  If the hash algorithm is any good, this should be
	// easy.

	testStrings := []string{
		"",
		"1",
		"12",
		"123",
		"1234",
		"12345",
		"123456",
		"1234567",
		"12345678",
	}

	// Compare all possible pairs of test strings.  This includes comparing each string to itself as well.
	for _, s := range testStrings {
		for _, t := range testStrings {

			var diffHashS, diffHashT DiffHash
			diffHashS.Init(s)
			diffHashT.Init(t)
			similarityFactor := diffHashS.Similarity(diffHashT)

			// Check for NaN similarity factor.
			if math.IsNaN(float64(similarityFactor)) {
				tt.Errorf("DiffHash: %q and %q should have a similarity factor that is a number, not %f!", s, t, similarityFactor)
			}

			// Check for a similarity factor that's completely out of range.
			if similarityFactor < 0.0 {
				tt.Errorf("DiffHash: %q and %q should have a similarity factor >= 0.0, but instead have %f", s, t, similarityFactor)
			} else if similarityFactor > 1.0 {
				tt.Errorf("DiffHash: %q and %q should have a similarity factor <= 1.0, but instead have %f", s, t, similarityFactor)
			}

			// Check for 100% similarity or not 100% similarity.
			if s == t {
				if similarityFactor != 1.0 {
					tt.Errorf("DiffHash: %q and %q should have a 1.0 similarity factor (they are the same), " +
								"however a '%f' similarity factor was reported.", s, t, similarityFactor)
				}
			} else {
				if similarityFactor == 1.0 {
					tt.Errorf("DiffHash: %q and %q are different but a '%f' similarity factor was reported.", s, t, similarityFactor)
				}
			}
		}
	}
}

// -------------------------------------------
// ------------------------------------------- TestTextLine
// -------------------------------------------

func TestTextLine(tt *testing.T) {

	tester := NewTester(tt, "TextLine", nil)

	tester.PrintBanner("Testing TextLine")

	testStrings := []string{
		"",
		"1",
		"12",
		"123",
		"1234",
		"12345",
		"123456",
		"1234567",
		"12345678",
	}

	// Compare all possible pairs of test strings.  This includes comparing each string to itself as well.
	for _, s := range testStrings {
		for _, t := range testStrings {
			lineS := NewTextLine(s)
			lineT := NewTextLine(t)

			similarityFactor := lineS.Similarity(lineT)

			// Check for NaN similarity factor.
			if math.IsNaN(float64(similarityFactor)) {
				tt.Errorf("TextLine: %q and %q should have a similarity factor that is a number, not %f!", s, t, similarityFactor)
			}

			// Check for a similarity factor that's completely out of range.
			if similarityFactor < 0.0 {
				tt.Errorf("TextLine: %q and %q should have a similarity factor >= 0.0, but instead have %f", s, t, similarityFactor)
			} else if similarityFactor > 1.0 {
				tt.Errorf("TextLine: %q and %q should have a similarity factor <= 1.0, but instead have %f", s, t, similarityFactor)
			}

			// Check for 100% similarity or not 100% similarity.
			if s == t {
				if similarityFactor != 1.0 {
					tt.Errorf("TextLine: %q and %q should have a 1.0 similarity factor (they are the same), " +
								"however a '%f' similarity factor was reported.", s, t, similarityFactor)
				}
			} else {
				if similarityFactor == 1.0 {
					tt.Errorf("TextLine: %q and %q are different but a '%f' similarity factor was reported.", s, t, similarityFactor)
				}
			}
		}
	}
}

// -------------------------------------------
// ------------------------------------------- TestLevenshteinDistance
// -------------------------------------------

func TestLevenshteinDistance(t *testing.T) {

	tester := NewTester(t, "TestLevenshteinDistance", nil)

	charSet := []rune(CHAR_SET)

	tester.PrintBanner("generating initial strings")
	initialStrings := generateInitialStrings(t, charSet, 0, MAX_INITIAL_STRING_LEN)

	tester.PrintBanner("generating test cases")
	testCase := generateMasterTestCases(initialStrings, charSet, MAX_OP_COUNT)

	runTests(NewTester(t, "LevenshteinDistance_v1", LevenshteinDistance_v1), testCase)
	runTests(NewTester(t, "LevenshteinDistance_v2", LevenshteinDistance_v2), testCase)
	runTests(NewTester(t, "LevenshteinDistance_v3", LevenshteinDistance_v3), testCase)
	runTests(NewTester(t, "LevenshteinDistance_v4", LevenshteinDistance_v4), testCase)
	runTests(NewTester(t, "LevenshteinDistance_v5", LevenshteinDistance_v5), testCase)
	runTests(NewTester(t, "LevenshteinDistance_v6", LevenshteinDistance_v6), testCase)

	// Single ad hoc test for now...
	// singleTestCase := &tLDTestCase{
	// 	"Now is the time for men to come to the aid of their country",  
	// 	"Now is the time for all good men to come to the aid of their country",
	// 	9,
	// }
	// runTests(t, singleTestCase, "LevenshteinDistance_v5", LevenshteinDistance_v5)
}

// -------------------------------------------
// ------------------------------------------- TestDiff
// -------------------------------------------

var pairsOfSimilarStrings [][]string = [][]string{
	[]string{"He’s Alive!", "It’s Alive!"},
	[]string{"Mirror, mirror on the wall, who is the fairest of them all?", "Magic Mirror on the Wall, who is the Fairest one of all?"},
	[]string{"Toto, I don't think we're in Kansas anymore.", "Toto, I've a feeling we're not in Kansas anymore."},
	[]string{"Beam me up Scotty.", "Beam us up Scotty."},
	[]string{"Mrs. Robinson, are you trying to seduce me?", "Mrs. Robinson, you're trying to seduce me. Aren't you?"},
	[]string{"May the Force Be With You", "Remember, the Force will be with you ... always"},
	[]string{"Luke, I am your father.", "No, I am your father."},
	[]string{"If you build it, they will come", "If you build it, he will come"},
	[]string{"Hello, Clarice", "Good evening, Clarice"},
	[]string{"Ah, Houston, we've had a problem.", "Houston, we have a problem"},
}

func TestDiff(t *testing.T) {

	tester := NewTester(t, "TestDiff", nil)

	var lines1 ComparableLines
	var lines2 ComparableLines
	for _, stringPair := range pairsOfSimilarStrings {
		lines1 = append(lines1, NewTextLine(stringPair[0]))
		lines2 = append(lines2, NewTextLine(stringPair[1]))
	}

	tester.PrintBanner("Diff_v2")
	distance, alignment := Diff_v2(lines1, lines2)
	alignment.Dump(lines1, lines2, int(distance), tester)
}

// -------------------------------------------
// ------------------------------------------- TestDiff2
// -------------------------------------------

func TestDiff2(t *testing.T) {

	tester := NewTester(t, "Diff_v2", nil)

	tester.PrintBanner("Diff_v2")

	// If we generate all the numbers between 0 (inclusive) and 2^n (exclusive), then
	// we also generate all the possible bit patterns for n bits.  This can be an easy
	// way to generate an exhaustive (for some definition) list of test cases of a
	// given length.  In this case we construct each test case from a sequence of m
	// 2-bit control codes (n = 2 * m in this case).  We use the i'th 2-bit sequence
	// (bits 2 * i and 2 * i + 1) to determine what to do with the pair of strings at
	// position i in the source test data.

	tester.PrintHeader("Constructing Test Cases")

	var codeCount uint32 = 5					// we will build test cases from the first m string pairs
	codeBitCount := 2 * codeCount				// we need 2 bits for each string pair
	codeBitStringCount := 1 << codeBitCount		// we will generate the range [0..2 ^ (2 * codeCount)]

	// Construct test cases.
	var testCases []TestCase
	for controlCodeBitString := 0; controlCodeBitString < codeBitStringCount; controlCodeBitString++ {
		var leftLines, rightLines ComparableLines
		remainingControlCodeBitString := controlCodeBitString
		for i := uint32(0); i < codeCount; i++ {
			controlCode := remainingControlCodeBitString & 3					// read the two lowest order bits
			remainingControlCodeBitString = remainingControlCodeBitString >> 2	// shift away the two lowest order bits
			switch controlCode {
			case 0:
				leftLines = append(leftLines, NewTextLine(pairsOfSimilarStrings[i][0]))
			case 1:
				rightLines = append(rightLines, NewTextLine(pairsOfSimilarStrings[i][1]))
			case 2:
				leftLines = append(leftLines, NewTextLine(pairsOfSimilarStrings[i][0]))
				rightLines = append(rightLines, NewTextLine(pairsOfSimilarStrings[i][0]))
			case 3:
				leftLines = append(leftLines, NewTextLine(pairsOfSimilarStrings[i][0]))
				rightLines = append(rightLines, NewTextLine(pairsOfSimilarStrings[i][1]))
			default:
				panic("not reached")
			}
		}
		testCases = append(testCases, NewDiffTestCase(leftLines, rightLines))
	}

	// Execute the test cases.
	tester.PrintHeader("Testing", len(testCases), "test cases.")
	for _, testCase := range testCases {
		testCase.execute(tester)
	}

	tester.PrintSummary("Done.")
}

// -------------------------------------------
// ------------------------------------------- Levenshtein Distance functions
// -------------------------------------------

/*
	The Levenshtein Distance algorithm computes the smallest number of single
	character insertions, deletions, and replacements that will convert one
	string into another one.

	We start with:
	* an initial string
	* a set of characters to use
	* a maximum number of instructions, which we will just call "n"

	We generate and execute all possible legal sequences of n or fewer
	instructions, with the goal of recording all unique output strings
	as well as the shortest instruction sequence which created each
	unique output string.

	Notes:
	* A particular output string can be created by multiple different
	  instruction sequences, some longer, some shorter, and indeed 
	  many of the same the length; we only care about the length of
	  the shortest sequences and it doesn't matter if there are
	  multiple such shortest sequences.
	* By executing each instruction against the result string
	  of the previous instruction in a sequence we can skip illegal
	  cases when we generate the set of all possible following
	  instructions in the sequence.
	* You can think of the algorithm as traversing the tree of depth
	  n of all possible insert, delete, and replace operations that can
	  be performed starting with the initial string.
	* Since the whole goal is to find the shortest sequence of
	  instructions which will produce a particular output string,
	  we must look at result strings and instruction count at each
	  node of the tree, not just at the leaves!
	* We only need to record each unique output string once and
	  then just need to make sure we update the associated
	  instruction count if we encounter a shorter one for that
	  string.  We could even skip the update part if we did a
	  breadth-first traversal, but depth-first is easier.
	* Recording each unique string only once greatly reduces the
	  amount of storage we need, and results in a much faster algorithm,
	  at least in Go v1.6.
	* Assuming that we have truly exhaustively explored our search space,
	  then we've assembled a set of
	    (initial string, result string, minimun instruction count)
	  triples, each of which can be used as a Levenshtein Distance test
	  case.  Note that does not mean the list of triples will
	  exhaustively test our Levenshtein Distance algorithm!

	 Initial Strings
	 ===============

	 Using the algorithm above to generate test cases does not
	 necessarily result in an exhaustive set of test cases for our
	 Levenshtein Distance algorithm.  Empirically speaking, the choice
	 of initial string can make a difference.  So we should probably
	 test over a range of initial strings as well.

	 Ideally we could use another exhaustive method for generating
	 initial strings, but this will result in a much larger set of
	 test cases and much longer test runs.

	 The strategy used here is to generate an exhaustive set of 
	 initial strings of length 0 to n, but then consolidate effectively
	 equivalent strings into single representative cases.

	 If we assume that we are using the same character set for both
	 initial strings and edit instructions, then some generated test
	 cases will be obviously equivalent.

	 Assume we're using the character set "abcd" to build test cases.
	 Now suppose we have an initial string "abb" and an exhaustive set
	 of test cases for that string.  Each test case will correspond to a
	 sequence of insertions, deletions, and replacements where
	 insertion and replacement operations also specify a character
	 drawn from the character set.

	 If we were to map the characters in the initial string and all
	 of the generated edit sequences like so:

	 	"a" => "c"
	 	"b" => "d"
	 	"c" => "a"
	 	"d" => "b"

	 then we should get the initial string "cdd" and all the same
	 edit sequences (and consequently, test cases) as if we just
	 started with "cdd" and used the exhaustive test case algorithm
	 to generate test cases for it directly.  The test cases might
	 be in a different order, but we don't care about the order.

	 So I am asserting that that "abb" and its test cases are
	 isomorphic to "cdd" and its test cases, and so there is no
	 need to generate test cases for both initial strings -- either
	 one will do.

	 Notes:
	 * I am asserting this without any sort of formal proof!
	 * Remember the assumption that we are always using the same
	   source character set.
	 * Also remember the assumption that we are generating
	   exhaustive test cases for each string that we will test

	 So the idea is to generate an exhaustive set of initial strings
	 of 0 to n characters, then reduce isomorphic initial strings
	 to single representative cases, and then generate exhaustive
	 tests for each of strings of this smaller set of initial
	 strings.
*/

// ------------------------------------------- type tTestCaseSpec

type tTestCaseSpec struct {
	initialString string
	charSet []rune
	maxOpCount int
}

// ------------------------------------------- type tLDTestCase

type tLDTestCase struct {
	s, t string
	distance int
}

// Assert that TestCase is implemented by tLDTestCase. 
var _ TestCase = (*tLDTestCase)(nil)

// ------------------------------------------- tLDTestCase execute

func (self *tLDTestCase) execute(tester *tTester) {

	executeTest := func (s, t string, distance int, note string) {
		computedDistance := tester.testeeFn(s, t)
	 	if computedDistance != distance {
	 		noteText := ""
	 		if note != "" {
	 			noteText = " (" + note + ")"
	 		}
	 		tester.Errorf("%s: %-9q %-9q got: %3d expected: %3d%s", tester.name, s, t, computedDistance, distance, noteText)
	 	}
	}

	s, t, distance := string(self.s), string(self.t), self.distance
	executeTest(s, t, distance, "forward")
	executeTest(t, s, distance, "backward")
}

// ------------------------------------------- runTests

func runTests(tester *tTester, testCase TestCase) {
	tester.PrintBanner(tester.name)
	testCase.execute(tester)
}

// ------------------------------------------- generateInitialStrings

func generateInitialStrings(t *testing.T, charSet []rune, minCharCount, maxCharCount int) []string {
	var strings []string
	for charCount := minCharCount; charCount <= maxCharCount; charCount++ {		// "maxCharCount" is included!
		strings = append(strings, generateUniqueCanonicalStrings(t, charSet, charCount)...)
	}
	return strings
}

// ------------------------------------------- generateMasterTestCases

func generateMasterTestCases(initialStrings []string, charSet []rune, maxOpCount int) TestCase {
	var testCases []TestCase
	for _, initialString := range initialStrings {
		testCase := generateTestCases(initialString, charSet, maxOpCount)
		testCases = append(testCases, testCase)
	}

	masterTestCase := MasterTestCase{title: "All test cases"}
	masterTestCase.appendSimpleTestCases(testCases)

	return &masterTestCase
}

// ------------------------------------------- generateTestCases

func generateTestCases(initialString string, charSet []rune, maxOpCount int) TestCase {

	var tStrings []string
	tStringToDistanceMap := make(map[string]int)

	var testCaseRecorder = func (t string, distance int) {
		// We want to identify the minimum edit distance for each unique pair of strings.
		prevDistance, found := tStringToDistanceMap[t]
		if !found {
			// This is the first time we've seen this string.
			tStringToDistanceMap[t] = distance
			tStrings = append(tStrings, t)
		} else {
			if distance < prevDistance {
				// We've seen this string before, but this time we have a shorter distance.
				tStringToDistanceMap[t] = distance
				// "t" has already been added to list, and we don't want to add it again.
			}
		}
	}

	// Note that any string can be converted to itself
	// with zero edits, but the test generator won't generate this
	// specific case, so we need to add it manually.
	testCaseRecorder(initialString, 0)

	spec := tTestCaseSpec{initialString, charSet, maxOpCount}
	createTestCases(&spec, initialString, 0, testCaseRecorder)

	testCases := make([]TestCase, 0, len(tStrings))
	for _, t := range tStrings {
		distance := tStringToDistanceMap[t]
		testCase := tLDTestCase{s: initialString, t: t, distance: distance}
		testCases = append(testCases, &testCase)
	}

	masterTestCase := MasterTestCase{title: initialString}
	masterTestCase.appendSimpleTestCases(testCases)

	return &masterTestCase
}

// ------------------------------------------- createTestCases

func createTestCases(spec *tTestCaseSpec, current string, currentOpCount int,
									recorderFn func (s string, distance int)) {

	if currentOpCount == spec.maxOpCount {
		return
	}

	const opCodeInsert = 1
	const opCodeDelete = 2
	const opCodeReplace = 3

	for opCode := opCodeInsert; opCode <= opCodeReplace; opCode++ {

		// Only *insert* makes sense for an empty string.
		if len(current) == 0 && opCode != opCodeInsert {
			continue
		}

		switch opCode {
		case opCodeInsert:

			// insert argument1, argument2
			//     where argument1 in [0..len(current)+1] -- last argument1 is len(current), an insert at the end
			//     and   argument2 in [0..len(charSet)]   -- last argument2 is len(charSet) - 1, as you'd expect
			for argument1 := 0; argument1 <= len(current); argument1++ {			// note "<="!
				for argument2 := 0; argument2 < len(spec.charSet); argument2++ {
					next := current[:argument1] + string(spec.charSet[argument2]) + current[argument1:]
					recorderFn(next, currentOpCount + 1)
					createTestCases(spec, next, currentOpCount + 1, recorderFn)
				}
			}

		case opCodeDelete:

			// delete argument
			//     where argument1 in [0..len(current)]   -- last argument1 is len(current) - 1, as you'd expect
			for argument := 0; argument < len(current); argument++ {
				next := current[:argument] + current[argument + 1:]
				recorderFn(next, currentOpCount + 1)
				createTestCases(spec, next, currentOpCount + 1, recorderFn)
			}

		case opCodeReplace:

			// replace argument1, argument2
			//     where argument1 in [0..len(current)]   -- last argument1 is len(current) - 1, as you'd expect
			//     and   argument2 in [0..len(charSet)]   -- last argument2 is len(charSet) - 1, as you'd expect
			for argument1 := 0; argument1 < len(current); argument1++ {
				for argument2 := 0; argument2 < len(spec.charSet); argument2++ {
					next := current[:argument1] + string(spec.charSet[argument2]) + current[argument1 + 1:]
					recorderFn(next, currentOpCount + 1)
					createTestCases(spec, next, currentOpCount + 1, recorderFn)
				}
			}

		default:
			panic("Unknown opcode!")
		}

	}

}

// ------------------------------------------- generateUniqueCanonicalStrings

func generateUniqueCanonicalStrings(t *testing.T, charSet []rune, charCount int) []string {
	strings := generateStrings(charSet, charCount)
	return selectUniqueCanonicalStrings(t, strings, charSet)
}

// ------------------------------------------- selectUniqueCanonicalStrings

func selectUniqueCanonicalStrings(t *testing.T, strings []string, charSet []rune) []string {
	var results []string
	seenItMap := make(map[string]bool)
	for _, s := range strings {
		canonicalKey := canonicalizeString(s, charSet)
		_, seenIt := seenItMap[canonicalKey]
		if !seenIt {
			seenItMap[canonicalKey] = true
			results = append(results, canonicalKey)
			t.Logf("%q ok\n", canonicalKey)
		} else {
			t.Logf("%q skipped (%q is equivalent)\n", s, canonicalKey)
		}
	}
	return results
}

// ------------------------------------------- canonicalizeString

func canonicalizeString(s string, charSet []rune) string {
	nextIndex := 0
	charMap := make(map[rune]rune)
	result := ""
	for _, char := range s {
		replacementChar, found := charMap[char]
		if !found {
			replacementChar = charSet[nextIndex]
			nextIndex++
			charMap[char] = replacementChar
		}
		result += string(replacementChar)
	}
	return result
}

// ------------------------------------------- generateStrings

func generateStrings(charSet []rune, charCount int) []string {
	var accumulator []string
	generateStringsIntoAccumulator("", charSet, charCount, &accumulator)
	return accumulator
}

// ------------------------------------------- generateStringsIntoAccumulator

func generateStringsIntoAccumulator(prefix string, charSet []rune, charCount int, accumulator *[]string) {
	if charCount < 1 {
		*accumulator = append(*accumulator, prefix)
	} else {
		for _, char := range charSet {
			newPrefix := prefix + string(char)
			generateStringsIntoAccumulator(newPrefix, charSet, charCount - 1, accumulator)
		}
	}
}

// -------------------------------------------
// ------------------------------------------- Diff stuff
// -------------------------------------------

// ------------------------------------------- type DiffTestCase

type DiffTestCase struct {
	leftLines, rightLines ComparableSequence
}

// Assert that TestCase is implemented by DiffTestCase. 
var _ TestCase = (*DiffTestCase)(nil)

// ------------------------------------------- NewDiffTestCase DiffTestCase factory function

func NewDiffTestCase(leftLines, rightLines ComparableSequence) *DiffTestCase {
	return &DiffTestCase{leftLines: leftLines, rightLines: rightLines}
}

// ------------------------------------------- DiffTestCase execute

func (self *DiffTestCase) execute(tester *tTester) {
	distance, alignment := Diff_v2(self.leftLines, self.rightLines)
	alignment.Dump(self.leftLines, self.rightLines, int(distance), tester)
	// TODO: Short of a panic, the test will never actually fail!
}
