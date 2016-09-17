package diff

import (
	"fmt"
	"sort"
)

// "text-line.go" - Types, methods, and functions for working with lines of text.

// -------------------------------------------
// ------------------------------------------- type DiffHash
// -------------------------------------------

/* 	.............................................................
	A DiffHash for one string can be compared to the DiffHash for
	another string in order to approximate their similarity.  The
	computed similarity is represented as a 32 bit floating point
	number between 0.0 and 1.0, where the former is (approximately)
	100% different and the latter is (approximately) 100% the
	same.  The only real guarantees are that the similarity factor
	is definitely between 0.0 and 1.0 and that the DiffHash values
	computed from identical strings will always result in a
	similarity factor of 1.0.  The primary value of the DiffHash
	is that the similarity factor can be computed much more quickly
	than Levenshtein distance or similar algorithm, at least for
	longer strings.
 	.............. */

type DiffHash struct {
	hashes []uint32
}

/*	......................................................................................................
	Basically DiffHash converts a string into a sequence of runes, and then generates a sequence of 32-bit
	hash values from the runes, by sliding a fixed-size "window" over the sequence of runes, where the
	window at each step overlaps all but the first rune of the window from the previous step.

	The algorithm is complicated by the need to cleanly handle strings shorter than the window size.  It
	is also necessary to be able to compare the DiffHash for strings smaller than the window size with
	strings that are as long or longer than the window size.  It is perfectly reasonable to compare "abc"
	with "abcde" and expect a greater than zero similarity factor.

	So the way the algorithm actually works is that it computes a sequence of hashes for each individual
	character as well as the sequence of hashes for all the windows, and then the two sequences of hashes
	are concatentated.  The final step is to sort the hashes.  This simplifies counting the number of
	hashes that two DiffHash values have in common.

	Some examples:

	"" => - => -																	# 0 runes, 0 hashes
	"a" => 'a' => |a|																# 1 rune, 1 hash
	"ab" => 'a', 'b' => |a|, |b|													# 2 runes, 2 hashes
	"abc" => 'a', 'b', 'c' => |a|, |b|, |c|											# 3 runes, 3 hashes
	"abcd" => 'a', 'b', 'c', 'd', ["abcd"] => |a|, |b|, |c|, |d|, ||abcd||			# 4 runes, 5 hashes
	"abcde" => 'a', 'b', 'c', 'd', 'e', "abcd", "bcde"
					=> |a|, |b|, |c|, |d|, |e|, ||abcd||, ||bcde||					# 5 runes, 7 hashes
	"abcded" => 'a', 'b', 'c', 'd', 'e', 'f', "abcd", "bcde", "cdef"
					=> |a|, |b|, |c|, |d|, |e|, |f|, ||abcd||, ||bcde||, ||cdef||	# 6 runes, 9 hashes

	where:

		"..." is a sequence of runes
		'x' is a single rune
		|x| is a 32 bit hash computed from a single rune
			- this *could* just be the unicode code point
		||xxxx|| is a 32-bit hash computed from exactly 4 runes
			- the 1-rune and 4-rune hash functions don't need to be the same

	Note that things get interesting when we get to "abcd" and have enough runes for a single 4-rune "window".
	With the next example, "abcde", there are enough runes for two 4-rune windows: "abcd" and "bcde".  For a
	window size of four runes, the formula for hash length as a function of rune length is

		hashLen = runeLen + max(0, runeLen - 3)

	........................................... */

// ------------------------------------------- DiffHash Init method

func (diffHash *DiffHash) Init(s string) {

	// Convert the string to a slice of runes.
	runes := []rune(s)
	runesLen := len(runes)

	// Create the hashes slice and initialize it with the runes.
	diffHash.hashes = make([]uint32, runesLen)
	for i, rune := range(runes) {
		diffHash.hashes[i] = uint32(rune)	// in this part, we are simply using rune values as "hashes"
	}

	// Add proper hashes to the hashes slice, if we can.
	if runesLen > 3 {
		hashCount := runesLen - 3 		// we will slide a 4-rune window down the length of the rune slice
		diffHash.hashes = append(diffHash.hashes, make([]uint32, hashCount)...)
		for i := 0; i < hashCount; i++ {

			// For each 4-rune window, we will compute a hash and append it to the hashes slice.
			// Note that each subsequent window overlaps the last 3 runes of the previous window.
			r0 := uint32(runes[i + 0])
			r1 := uint32(runes[i + 1])
			r2 := uint32(runes[i + 2])
			r3 := uint32(runes[i + 3])
			hash := rotateLeft(r0, 24) ^ rotateLeft(r1, 16) ^ rotateLeft(r2, 8) ^ r3
			diffHash.hashes[runesLen + i] = hash
		}
	}

	// Sort the hashes.
	sort.Sort(uint32_slice_sortAdaptor(diffHash.hashes))
}

// ------------------------------------------- DiffHash Similarity method

func (diffHash DiffHash) Similarity(diffHash2 DiffHash) float32 {

	hashLen, hashLen2 := len(diffHash.hashes), len(diffHash2.hashes)

	if hashLen == 0 && hashLen2 == 0 {
		return 1.0				// the empty string is 100% similar to the empty string!
	} else if hashLen == 0 || hashLen2 == 0 {
		return 0.0				// the empty string and any other string have 0% similarity
	}

	matchCount := 0
	for i, j := 0, 0; i < hashLen && j < hashLen2; {
		hash1, hash2 := diffHash.hashes[i], diffHash2.hashes[j]
		if hash1 == hash2 {
			i, j, matchCount = i + 1, j + 1, matchCount + 1
		} else if hash1 < hash2 {
			i++
		} else { // hash2 < hash1
			j++
		}
	}
	denominator := hashLen
	if hashLen2 > hashLen {
		denominator = hashLen2
	}

	return float32(matchCount) / float32(denominator)
}

// ------------------------------------------- rotateLeft

func rotateLeft(val uint32, shiftCount uint) uint32 {
	var bitCount uint = 32
	shiftLeftCount := shiftCount % bitCount
	shiftRightCount := bitCount - shiftLeftCount
	return (val << shiftLeftCount) | (val >> shiftRightCount)
}

// ------------------------------------------- sort adaptor type for uint32 slices

type uint32_slice_sortAdaptor []uint32
func (slice uint32_slice_sortAdaptor) Len() int           { return len(slice) }
func (slice uint32_slice_sortAdaptor) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }
func (slice uint32_slice_sortAdaptor) Less(i, j int) bool { return slice[i] < slice[j] }


// -------------------------------------------
// ------------------------------------------- type TextLine
// -------------------------------------------

// The TextLine type is used to represent a single line of text.
// Notably, each TextLine has a precomputed DiffHash so rapid 
// similarity computations can be made between two TextLines.

type TextLine struct {
	Text string
	diffHash DiffHash
}

// ------------------------------------------- NewTextLine TextLine factory function

func NewTextLine(text string) *TextLine {
	line := TextLine{Text:text}
	line.diffHash.Init(text)
	return &line
}

// ------------------------------------------- TextLine Similarity method

func (line1 *TextLine) Similarity(line2 *TextLine) float32 {
	similarityFactor := line1.diffHash.Similarity(line2.diffHash)
	if similarityFactor < 0.6 { similarityFactor = 0.0 }
	return similarityFactor
}

// ------------------------------------------- TextLine Compare method

func (line1 *TextLine) Compare(line2 Comparable) float32 {
	return 1.0 - line1.Similarity(line2.(*TextLine))
}

// ------------------------------------------- TextLine Stringify method

func (line *TextLine) Stringify(maxWidth int) string {
	runes := []rune(line.Text)
	if len(runes) > maxWidth {
		runes = runes[:maxWidth]
	}
	for i := maxWidth - 3; i < len(runes); i++ {
		if i > 0 {
			runes[i] = '.'
		}
	}
	return string(runes)
}

// ------------------------------------------- type ComparableLines

// Type ComparableLines is a TextLine slice subtype which implements the
// ComparableSequence interface.

type ComparableLines []*TextLine

// Assert that ComparableSequence is implemented by ComparableLines.
var _ ComparableSequence = (*ComparableLines)(nil)

// -------------------------------------------

func (slice ComparableLines) Length() int {
	return len(slice)
}

// -------------------------------------------

func (slice ComparableLines) GetItemAt(index int) Comparable {
	return slice[index]
}

// -------------------------------------------

func (slice ComparableLines) GetDescription() string {
	return fmt.Sprintf("%d lines", len(slice))
}
