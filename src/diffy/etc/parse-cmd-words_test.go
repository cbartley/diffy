package etc

import (
	"fmt"
	"strings"
	"testing"
)

// -------------------------------------------
// ------------------------------------------- helper functions for building test cases
// -------------------------------------------

func convertTabsToEscapeSequences(text string) string {
	return strings.Replace(text, "\t", `\t`, -1)
}

// Dump all of the strings in a slice, each preceded by its index.
func dumpStrings(slice []string, escapeTabs bool) {
	for index, s := range slice {
		if escapeTabs { s = convertTabsToEscapeSequences(s) }
		fmt.Printf("... %d |%s|\n", index, s)
	}
}

// Concatenate multiple string slices into one larger
// string slice and return it.
func concat(slices ...[]string) []string {
	count := 0
	for _, slice := range slices {
		count += len(slice)
	}

	resultSlice := make([]string, 0, count)
	for _, slice := range slices {
		resultSlice = append(resultSlice, slice...)
	}

	return resultSlice
}

// Compute the SQL style cross product of multiple string slices.  However,
// instead of returning tuples, concatenate each tuple into a single string.
// Then return a slice containing all of the new strings.
//
// crossProduct({"a"}, {"b"}) => {"ab"}
// crossProduct({"a", "b"}, {"1", "2"}) => {"a1", "a2", "b1", "b2"}
// crossProduct({"a"}, {"b"}, {"c", "d"}) => {"abc", "abd"}
// etc.
func crossProduct(slice []string, slices ...[]string) []string {
	if len(slices) == 0 {
		return slice
	}

	var resultSlice []string
	slice2 := crossProduct(slices[0], slices[1:]...)
	for _, s := range slice {
		for _, s2 := range slice2 {
			resultSlice = append(resultSlice, s + s2)
		}
	}
	return resultSlice
}

// Return a slice of strings the same size as the input slice, but 
// with each element replaced with the "replacement" argument.
// Typically "replacement" will either be the empty string or a
// separator character such as ",".  This is often useful for
// building expected output slices in cases where the outputs will
// generally be simpler than the inputs.  This is especially true
// when the inputs and expected outputs are constructed with
// "crossProduct()".
func replaceWith(slice []string, replacement string) []string {
	result := make([]string, len(slice), len(slice))
	for index := range result {
		result[index] = replacement
	}
	return result
}

// Return a new slice of strings where each output string is the
// corresponding input string, but wrapped in single quotes.
func singleQuoteItems(strings []string) []string {
	singleQuote := []string{`'`}
	return crossProduct(singleQuote, strings, singleQuote)
}

// Return a new slice of strings where each output string is the
// corresponding input string, but wrapped in double quotes.
func doubleQuoteItems(strings []string) []string {
	doubleQuote := []string{`"`}
	return crossProduct(doubleQuote, strings, doubleQuote)
}

// Return a new slice of strings where each string is a "sentence" made
// up of one or more words where each word is drawn from the "words" 
// slice.  Conceptually, sentences have leading and trailing whitespace
// drawn from "wsEnds" and each adjacent word (if there is more than
// one) is separated by whitespace drawn from "wsSeparators".
//
// In general: sentences(...) => crossProduct(wsEnds, words, {wsSeparators, words}*, wsEnds)
//
// For example: sentences(3, {"a"}, {" "}, {" "}) => {" a a a "}
//
// Since this function is built on top of "crossProduct()", it should be
// used judiciously as it could generate a very large output slice.
func sentences(count int, words, wsEnds, wsSeparators []string) []string {
	if count < 1 { panic("'count' must be at least '1'") }
	var arguments [][]string
	arguments = append(arguments, words)
	for i := 1; i < count; i++ {
		arguments = append(arguments, wsSeparators)
		arguments = append(arguments, words)
	}
	arguments = append(arguments, wsEnds)
	return crossProduct(wsEnds, arguments...)
}

// -------------------------------------------
// ------------------------------------------- test ParseWords
// -------------------------------------------

// ------------------------------------------- run_ParseWords_Tests

func run_ParseWords_Tests(t *testing.T, inputs []string, expectedOutputs []string, separator string) {

	if len(inputs) != len(expectedOutputs) {
		panic("The 'inputs' and 'outputs' slices must be the same length!")
	}

	for index := range inputs {
		input := inputs[index]
		expectedOutput := expectedOutputs[index]
		outputWords := ParseWords(input)
		output := strings.Join(outputWords, separator)
		if output != expectedOutput {
			inputT := convertTabsToEscapeSequences(input)
			outputT := convertTabsToEscapeSequences(output)
			expectedOutputT := convertTabsToEscapeSequences(expectedOutput)
			msg := fmt.Sprintf("ParseWords: |%s| => |%s|; expected |%s|", inputT, outputT, expectedOutputT)
			t.Error(msg)
		}
	}
}

// ------------------------------------------- TestParseWords

func TestParseWords(t *testing.T) {

	// ---
	// --- simple test cases drawn from documentation examples
	// ---

	inputs := []string{
		`abc`,
		` abc `,
		`abc 123`,
		`abc '1 2'`,
		`abc "1 2"`,
		`abc '1 2'"1 2"`,
		`abc "'123'"`,
		`abc '"123"'`,
	}

	expectedOutputs := []string{
		`abc`,
		`abc`,
		`abc,123`,
		`abc,1 2`,
		`abc,1 2`,
		`abc,1 21 2`,
		`abc,'123'`,
		`abc,"123"`,
	}

	run_ParseWords_Tests(t, inputs, expectedOutputs, ",")

	// ---
	// --- define the basic building blocks we'll use for creating test cases
	// ---

	// The very basics.
	emptyString := []string{""}
	spaceString := []string{" "}
	tabString := []string{"\t"}

	// General whitespace.
	ws1x := concat(spaceString, tabString)	// one whitespace character
	ws2x := crossProduct(ws1x, ws1x)		// two whitespace characters

	// Two kinds of end-whitespace: 0/1 whitespace characters and 0/1/2 whitespace characters.
	wsEnds0x1x := concat(emptyString, ws1x)
	wsEnds0x1x2x := concat(emptyString, ws1x, ws2x)
	wsEnds0x1xParsed := replaceWith(wsEnds0x1x, "")
	wsEnds0x1x2xParsed := replaceWith(wsEnds0x1x2x, "")
	_ = wsEnds0x1x
	_ = wsEnds0x1x2x
	_ = wsEnds0x1xParsed
	_ = wsEnds0x1x2xParsed

	// Two kinds of separator whitespace: 1 whitespace character and 1/2 whitespace characters.
	wsSeps1x := ws1x
	wsSeps1x2x := concat(ws1x, ws2x)
	wsSeps1xParsed := replaceWith(wsSeps1x, ",")
	wsSeps1x2xParsed := replaceWith(wsSeps1x2x, ",")

	// Some basic "word" definitions.
	simpleWords := []string{"a", "ab", "abc"}
	words := concat(
		simpleWords,
		singleQuoteItems(simpleWords),
		doubleQuoteItems(simpleWords),
	)
	wordsParsed := concat(
		simpleWords,
		simpleWords,
		simpleWords,
	)

	// ---
	// --- test sentences of 1, 2, and 3 words with at most 1 character of contiguous whitespace
	// ---

	// --- 1 word, whitespace: 0 or 1 chars at beginning and end ---

	sentence1 := sentences(1, words, wsEnds0x1x, nil)
	sentence1Parsed := sentences(1, wordsParsed, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, sentence1, sentence1Parsed, ",")

	// --- simple 2 word sentence, whitespace: 0 or 1 chars at beginning and end, always 1 char separating words ---

	sentence2 := sentences(2, words, wsEnds0x1x, wsSeps1x)
	sentence2Parsed := sentences(2, wordsParsed, wsEnds0x1xParsed, wsSeps1xParsed)
	run_ParseWords_Tests(t, sentence2, sentence2Parsed, ",")

	// --- simple 3 word sentence, whitespace: 0 or 1 chars at beginning and end, always 1 char separating words ---

	sentence3 := sentences(3, words, wsEnds0x1x, wsSeps1x)
	sentence3Parsed := sentences(3, wordsParsed, wsEnds0x1xParsed, wsSeps1xParsed)
	run_ParseWords_Tests(t, sentence3, sentence3Parsed, ",")

	// ---
	// --- test words that are quoted sentences of other words
	// ---

	// --- one single quoted word which may contain component double quoted words ---

	wordsAndDoubleQuotedWords := concat(simpleWords, doubleQuoteItems(simpleWords))
	sentencesOfWordsAndDoubleQuotedWords := sentences(2, wordsAndDoubleQuotedWords, wsEnds0x1x, wsSeps1x)
	singleQuotedSentences := sentences(1, singleQuoteItems(sentencesOfWordsAndDoubleQuotedWords), wsEnds0x1x, nil)
	singleQuotedSentencesParsed := sentences(1, sentencesOfWordsAndDoubleQuotedWords, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, singleQuotedSentences, singleQuotedSentencesParsed, ",")

	// --- one double quoted word which may contain component single quoted words ---

	wordsAndSingleQuotedWords := concat(simpleWords, singleQuoteItems(simpleWords))
	sentencesOfWordsAndSingleQuotedWords := sentences(2, wordsAndSingleQuotedWords, wsEnds0x1x, wsSeps1x)
	doubleQuotedSentences := sentences(1, doubleQuoteItems(sentencesOfWordsAndSingleQuotedWords), wsEnds0x1x, nil)
	doubleQuotedSentencesParsed := sentences(1, sentencesOfWordsAndSingleQuotedWords, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, doubleQuotedSentences, doubleQuotedSentencesParsed, ",")

	// ---
	// --- test sentences of 1 and 2 words with up to 2 characters of contiguous whitespace
	// ---

	// --- 1 word, whitespace: 0, 1, or 2 chars at beginning and end ---

	sentence1max2 := sentences(1, words, wsEnds0x1x2x, nil)
	sentence1max2Parsed := sentences(1, wordsParsed, wsEnds0x1x2xParsed, nil)
	run_ParseWords_Tests(t, sentence1max2, sentence1max2Parsed, ",")

	// --- simple 2 word sentence, whitespace: 0, 1, or 2 chars at beginning and end, 1 or 2 chars separating words ---

	sentence2max2 := sentences(2, words, wsEnds0x1x2x, wsSeps1x2x)
	sentence2max2Parsed := sentences(2, wordsParsed, wsEnds0x1x2xParsed, wsSeps1x2xParsed)
	run_ParseWords_Tests(t, sentence2max2, sentence2max2Parsed, ",")

	// ---
	// --- test adjacent words not separated by whitespace
	// ---

	// --- 2 adjacent words not separated by whitespace should be joined into a single word ---

	// We want to test cases like `'word''word'`, `word'word'`, `"word"word`, etc.
	// Note that we will also test cases like `wordword` which aren't even special
	// cases, but that doesn't cause any practical problems.  In all these examples
	// the expected output is `wordword`.

	twoAdjacentWords := crossProduct(wsEnds0x1x, crossProduct(words, words), wsEnds0x1x)
	twoAdjacentWordsParsed := crossProduct(wsEnds0x1xParsed, crossProduct(wordsParsed, wordsParsed), wsEnds0x1xParsed)
	run_ParseWords_Tests(t, twoAdjacentWords, twoAdjacentWordsParsed, ",")

	// ---
	// --- test nested quotes
	// ---

	stuffToQuote := concat(emptyString, ws1x, []string{"a"})

	// --- level one: single and double quotes, no nesting ---

	containsSingleQuotes_L1 := sentences(1, singleQuoteItems(stuffToQuote), wsEnds0x1x, nil)
	parsedSingleQuotes_L1 := sentences(1, stuffToQuote, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, containsSingleQuotes_L1, parsedSingleQuotes_L1, ",")

	containsDoubleQuotes_L1 := sentences(1, doubleQuoteItems(stuffToQuote), wsEnds0x1x, nil)
	parsedDoubleQuotes_L1 := sentences(1, stuffToQuote, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, containsDoubleQuotes_L1, parsedDoubleQuotes_L1, ",")

	// --- two levels: single around double, double around single ---

	containsSingleQuotes_L2 := sentences(1, singleQuoteItems(containsDoubleQuotes_L1), wsEnds0x1x, nil)
	parsedSingleQuotes_L2 := sentences(1, containsDoubleQuotes_L1, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, containsSingleQuotes_L2, parsedSingleQuotes_L2, ",")

	containsDoubleQuotes_L2 := sentences(1, doubleQuoteItems(containsSingleQuotes_L1), wsEnds0x1x, nil)
	parsedDoubleQuotes_L2 := sentences(1, containsSingleQuotes_L1, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, containsDoubleQuotes_L2, parsedDoubleQuotes_L2, ",")

	// --- three levels: single around double around single, double around single around double ---

	containsSingleQuotes_L3 := sentences(1, singleQuoteItems(containsDoubleQuotes_L2), wsEnds0x1x, nil)
	parsedSingleQuotes_L3 := sentences(1, containsDoubleQuotes_L2, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, containsSingleQuotes_L3, parsedSingleQuotes_L3, ",")

	containsDoubleQuotes_L3 := sentences(1, doubleQuoteItems(containsSingleQuotes_L2), wsEnds0x1x, nil)
	parsedDoubleQuotes_L3 := sentences(1, containsSingleQuotes_L2, wsEnds0x1xParsed, nil)
	run_ParseWords_Tests(t, containsDoubleQuotes_L3, parsedDoubleQuotes_L3, ",")
}
