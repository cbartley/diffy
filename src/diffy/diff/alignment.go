package diff

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

// -------------------------------------------

func (alignment *tAlignment) dump(left, right ComparableSequence, computedEditDistance int, s SimpleLogger) {

	s.Printf(".................................................... ")
	s.Printf("%s/%s (edit distance: %d)\n", left.GetDescription(), right.GetDescription(), computedEditDistance)
	s.Println()

	s.Printf("edit sequence\n")
	s.Printf("=============\n")

	s.Println()
	matchingCount := 0
	for _, link := range alignment.links {
		codeChar := " "
		var leftItem, rightItem Comparable = NewTextLine("-"), NewTextLine("-")
		switch link.linkType {
		case Matching:
			codeChar = " "
			leftItem, rightItem = left.GetItemAt(link.leftIndex), right.GetItemAt(link.rightIndex)
			matchingCount++
		case Different:
			codeChar = "*"
			leftItem, rightItem = left.GetItemAt(link.leftIndex), right.GetItemAt(link.rightIndex)
		case LeftOnly:
			codeChar = "-"
			leftItem = left.GetItemAt(link.leftIndex)
		case RightOnly:
			codeChar = "+"
			rightItem = right.GetItemAt(link.rightIndex)
		default:
			panic("Missing case")
		}
		s.Printf("%s %2d %-30s %-30s %2d\n", codeChar, link.leftIndex, leftItem.Stringify(30), rightItem.Stringify(30), link.rightIndex)
	}
	s.Println()

	s.Printf("first column legend\n")
	s.Printf("-------------------\n")
	s.Printf("%q copy\n", " ")
	s.Printf("%q change\n", "*")
	s.Printf("%q insert\n", "+")
	s.Printf("%q delete\n", "-")

	s.Println()
	nonMatchingCount := len(alignment.links) - matchingCount
	s.Printf("non-matching count, computed edit distance = %d, %d\n", nonMatchingCount, computedEditDistance)
	s.Println()
}

