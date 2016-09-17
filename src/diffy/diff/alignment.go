package diff

// -------------------------------------------
// ------------------------------------------- type Alignment
// -------------------------------------------

// The "Alignment" type represents an "alignment" between two sequences of items
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

type Alignment struct {
	Links []Link
}

// -------------------------------------------

type LinkType int

const (
	Matching LinkType = iota	// we have a left item index and a right item index, and the items *match*
	Different 					// we have a left item index and a right item index, and the items are *different*
	LeftOnly 					// we only have a left item index
	RightOnly 					// we only have a right item index
)

// -------------------------------------------

type Link struct {
	LinkType LinkType 	// link type
	LeftIndex int 		// -1 or zero-based index into the left or first sequence
	RightIndex int 		// -1 or zero-based index into the right or second sequence
}

// ------------------------------------------- Alignment RealignUsingThreshold
//
// Generate a nicer alignment using a thresholded similarity comparison.
//
func (alignment *Alignment) RealignUsingThreshold(left, right ComparableSequence, threshold float32) *Alignment {

	leftItem := func (link Link) Comparable {
		return left.GetItemAt(link.LeftIndex)
	}

	rightItem := func (link Link) Comparable {
		return right.GetItemAt(link.RightIndex)
	}

	// The naive diff algorithm will often pair up items that are not remotely similar, and just label
	// them as "different".  In many of these cases what we really want is one or more LeftOnly links
	// followed by one or more RightOnly links.  This algorithm uses a similarity measure and a threshold
	// to figure out when lines that are "aligned" shouldn't be.  In these cases it puts LeftOnly links
	// first followed by RightOnly links, until, of course, the run is interrupted by a pair of lines
	// that really are similar enough that they should be treated as fully aligned.

	var newLinks, rightLinks []Link
	for _, link := range alignment.Links {
		if link.LinkType == Different && leftItem(link).Compare(rightItem(link)) > threshold {
			newLinks = append(newLinks, Link{LeftOnly, link.LeftIndex, -1})
			rightLinks = append(rightLinks, Link{RightOnly, -1, link.RightIndex})
		} else {
			newLinks = append(newLinks, rightLinks...)	// append outstanding right links, if any
			rightLinks = rightLinks[:0]					// reset outstanding right links slice
			newLinks = append(newLinks, link)			// append the current link as-is
		}
	}
	newLinks = append(newLinks, rightLinks...)	// we might have some outstanding right links, append them
	return &Alignment{newLinks}
}

// ------------------------------------------- Alignment Dump

func (alignment *Alignment) Dump(left, right ComparableSequence, computedEditDistance int, s SimpleLogger) {

	s.Printf(".................................................... ")
	s.Printf("%s/%s (edit distance: %d)\n", left.GetDescription(), right.GetDescription(), computedEditDistance)
	s.Println()

	s.Printf("edit sequence\n")
	s.Printf("=============\n")

	s.Println()
	matchingCount := 0
	for _, link := range alignment.Links {
		codeChar := " "
		var leftItem, rightItem Comparable = NewTextLine("-"), NewTextLine("-")
		switch link.LinkType {
		case Matching:
			codeChar = " "
			leftItem, rightItem = left.GetItemAt(link.LeftIndex), right.GetItemAt(link.RightIndex)
			matchingCount++
		case Different:
			codeChar = "*"
			leftItem, rightItem = left.GetItemAt(link.LeftIndex), right.GetItemAt(link.RightIndex)
		case LeftOnly:
			codeChar = "-"
			leftItem = left.GetItemAt(link.LeftIndex)
		case RightOnly:
			codeChar = "+"
			rightItem = right.GetItemAt(link.RightIndex)
		default:
			panic("Missing case")
		}
		s.Printf("%s %2d %-30s %-30s %2d\n", codeChar, link.LeftIndex, leftItem.Stringify(30), rightItem.Stringify(30), link.RightIndex)
	}
	s.Println()

	s.Printf("first column legend\n")
	s.Printf("-------------------\n")
	s.Printf("%q copy\n", " ")
	s.Printf("%q change\n", "*")
	s.Printf("%q insert\n", "+")
	s.Printf("%q delete\n", "-")

	s.Println()
	nonMatchingCount := len(alignment.Links) - matchingCount
	s.Printf("non-matching count, computed edit distance = %d, %d\n", nonMatchingCount, computedEditDistance)
	s.Println()
}
