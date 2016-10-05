package etc

// ------------------------------------------- parseWords
// Parse the contents of "text" into a list of "words" and return the words
// as a slice of strings.  We're using "word" in the Unix shell sense:
// basically strings of characters separated from other strings of characters
// by whitespace.  Words may also be quoted with single or double quotes or
// may contain parts which are quoted with double or single quotes.  In
// these cases words may also contain whitespace.
//
// The "ParseWords()" function allows the nesting of alternating single and
// double quotes (or alternating double and single quotes) as many levels
// deep as you may want.  This capability is probably overkill, but it's
// there if you want it.  Note that top-level quotes are removed and nested
// quotes are preserved.
//
// ParseWords(`abc`) 			=> {`abc`}
// ParseWords(` abc `) 			=> {`abc`}
// ParseWords(`abc 123`) 		=> {`abc`, `123`}
// ParseWords(`abc '1 2'`) 		=> {`abc`, `1 2`}		# note single quotes are stripped
// ParseWords(`abc "1 2"`) 		=> {`abc`, `1 2`}		# note double quotes are stripped
// ParseWords(`abc '1 2'"1 2"`)	=> {`abc`, `1 21 2`}	# quotes are stripped and adjacent text is combined
// ParseWords(`abc "'123'"`) 	=> {`abc`, `'123'`}		# top-level quotes stripped, nested quotes preserved
// ParseWords(`abc '"123"'`) 	=> {`abc`, `"123"`}		# top-level quotes stripped, nested quotes preserved
// etc.
//
func ParseWords(text string) []string {
	var words []string
	runes := []rune(text)
	for index := 0; index < len(runes); {
		if word, next, matched := parseTopLevelWord(runes, index); matched {
			words = append(words, string(word))
			index = next
		} else if char := runes[index]; char == ' ' || char == '\t' {
			index += 1
		} else {
			panic("not reached")
		}
	}
	return words
}

// ------------------------------------------- parseTopLevelWord
// Parse a "top-level" word starting at position "start" in the "runes" slice.
// If a word is matched, return the matched word as a rune slice, the next
// position in the "runes" slice *after* the last matched rune, and true.
// Otherwise return false.  
// 
// Notes: 
// 
// * The function only returns true if at least one character in "runes" is matched.
//   However, the matched word itself could be empty if quotes are used.
// * Top-level words may be quoted or may contain one or more quoted parts.  
// * Adjacent "parts" (unquoted, single quoted, or double quoted) will be concatenated,
//   minus any top-level quotes.  
// * Top-level quotes are stripped but any embedded (nested) quotes will be preserved.
// * Whitespace may only appear within quoted parts.  Any other whitespace would mark
//   the end of the word.
//
func parseTopLevelWord(runes []rune, start int) ([]rune, int, bool) {
	var accumulator []rune
	matchedSomething := false
	index := start
	for ; index < len(runes); {
		if next, matched := parseDoubleQuotedString(runes, index); matched {
			matchedSomething = true		// but we might have matched a quoted empty string!
			accumulator = append(accumulator, runes[index + 1:next - 1]...)
			index = next
		} else if next, matched := parseSingleQuotedString(runes, index); matched {
			matchedSomething = true		// but we might have matched a quoted empty string!
			accumulator = append(accumulator, runes[index + 1:next - 1]...)
			index = next
		} else if char := runes[index]; char != ' ' && char != '\t' {
			matchedSomething = true
			accumulator = append(accumulator, char)
			index += 1
		} else {
			break						// we hit whitespace or the end of the string
		}
	}
	return accumulator, index, matchedSomething
}

// ------------------------------------------- parseDoubleQuotedString
// Parse a double quoted string starting at position "start" in the "runes" slice.
// If a string is matched, return the next position in the "runes" slice *after*
// the last matched rune and true.  Otherwise return false.
//
func parseDoubleQuotedString(runes []rune, start int) (int, bool) {
	// We must start with a double quote, otherwise we're done.
	if start < len(runes) && runes[start] != '"' {
		return start, false		// no starting quote
	}

	// Find the matching end quote, skipping over any single quoted substrings.
	for index := start + 1; index < len(runes); {
		if runes[index] == '"' {
			return index + 1, true
		} else if next, matched := parseSingleQuotedString(runes, index); matched {
			index = next
		} else {
			index += 1
		}
	}
	return start, false		// unmatched quote
}

// ------------------------------------------- parseSingleQuotedString
// Parse a single quoted string starting at position "start" in the "runes" slice.
// If a string is matched, return the next position in the "runes" slice *after*
// the last matched rune and true.  Otherwise return false.
//
func parseSingleQuotedString(runes []rune, start int) (int, bool) {
	// We must start with a single quote, otherwise we're done.
	if start < len(runes) && runes[start] != '\'' {
		return start, false		// no starting quote
	}

	// Find the matching end quote, skipping over any double quoted substrings.
	for index := start + 1; index < len(runes); {
		if runes[index] == '\'' {
			return index + 1, true
		} else if next, matched := parseDoubleQuotedString(runes, index); matched {
			index = next
		} else {
			index += 1
		}
	}
	return start, false		// unmatched quote
}
