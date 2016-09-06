package diff

// -------------------------------------------
// -------------------------------------------
// -------------------------------------------

type Comparable interface {
	Compare(other Comparable) float32
	Stringify(maxWidth int) string
}

// -------------------------------------------
// -------------------------------------------
// -------------------------------------------

type ComparableRune rune

// Assert that Comparable is implemented by ComparableRune.
var _ Comparable = ComparableRune(0)

// -------------------------------------------

func (c ComparableRune) Compare(d Comparable) float32 {
	if c == d.(ComparableRune) {
		return 0.0
	}
	return 1.0
}

// -------------------------------------------

func (c ComparableRune) Stringify(maxWidth int) string {
	return string(c)
}

// -------------------------------------------
// -------------------------------------------
// -------------------------------------------

type ComparableSequence interface {
	Length() int
	GetItemAt(int) Comparable
	GetDescription() string
}

// -------------------------------------------
// -------------------------------------------
// -------------------------------------------

type ComparableString []rune

// Assert that ComparableSequence is implemented by ComparableString.
var _ ComparableSequence = ComparableString(nil)

// ------------------------------------------- MakeComparableString ComparableString factory function

func MakeComparableString(s string) ComparableString {
	return ComparableString([]rune(s))
}

// ------------------------------------------- ComparableString xxx

func (s ComparableString) Length() int {
	return len(s)
}

// ------------------------------------------- ComparableString xxxx

func (s ComparableString) GetItemAt(index int) Comparable {
	return ComparableRune(s[index])
}

// ------------------------------------------- ComparableString xxxx

func (s ComparableString) GetDescription() string {
	return string(s)
}


