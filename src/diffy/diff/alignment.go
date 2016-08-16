package diff

import "fmt"

// -------------------------------------------
// -------------------------------------------
// -------------------------------------------

// -------------------------------------------

// This type represents an "alignment" between two sequences of items
// * each sequence may be a subsequence of a larger sequence
// * the sequences are contiguous -- they don't have any holes
// * each link has either a "left index" a "right index" or both
// * links are never empty, there is at least one index
// * "missing" indexes are represented by a "-1" index
// * present (i.e. not missing) indexes are ascending, that is, if a
//   left index is present in a link, its value is exactly one more
//   than the previous present left index; of course there may be
//   missing left indexes in between; the same is true of right indexes
// * each link has a type, one of:
//   - "Matching":   both indexes are present, and the referenced items are equivalent
//   - "Different":  both indexes are present, but the referenced items are different
//   - "LeftOnly":   only the left index is present, the right index is -1
//   - "RightOnly":  only the right index is present, the left index is -1

type tAlignment struct {
	links []tLink
}

// -------------------------------------------

type tLinkType int

const (
	Matching tLinkType = iota	// we have a left item index and a right item index, and the items *match*
	Different 					// we have a left item index and a right item index, and the items are *different*
	LeftOnly 					// we only have a left item index
	RightOnly 					// we only have a right item index
)

// -------------------------------------------

type tLink struct {
	linkType tLinkType 	// link type
	leftIndex int 		// -1 or zero-based index into the left or first sequence
	rightIndex int 		// -1 or zero-based index into the right or second sequence
}

// ------------------------------------------- tAlignment dump

func (alignment *tAlignment) dump(left, right string, computedEditDistance int) {

	fmt.Printf(".................................................... ")
	fmt.Printf("%s/%s (edit distance: %d)\n", left, right, computedEditDistance)
	fmt.Println()

	fmt.Printf("edit sequence\n")
	fmt.Printf("=============\n")

	fmt.Println()
	matchingCount := 0
	for _, link := range alignment.links {
		codeChar := " "
		leftItem, rightItem := ".", "."
		switch link.linkType {
		case Matching:
			codeChar = " "
			leftItem, rightItem = string(left[link.leftIndex]), string(right[link.rightIndex])
			matchingCount++
		case Different:
			codeChar = "*"
			leftItem, rightItem = string(left[link.leftIndex]), string(right[link.rightIndex])
		case LeftOnly:
			codeChar = "-"
			leftItem = string(left[link.leftIndex])
		case RightOnly:
			codeChar = "+"
			rightItem = string(right[link.rightIndex])
		default:
			panic("Missing case")
		}
		fmt.Printf("%s %2d %s %s %2d\n", codeChar, link.leftIndex, leftItem, rightItem, link.rightIndex)
	}
	fmt.Println()

	fmt.Printf("first column legend\n")
	fmt.Printf("-------------------\n")
	fmt.Printf("%q copy\n", " ")
	fmt.Printf("%q change\n", "*")
	fmt.Printf("%q insert\n", "+")
	fmt.Printf("%q delete\n", "-")

	fmt.Println()
	nonMatchingCount := len(alignment.links) - matchingCount
	fmt.Printf("non-matching count, computed edit distance = %d, %d\n", nonMatchingCount, computedEditDistance)
	fmt.Println()
}
