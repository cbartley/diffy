package diff

import "fmt"

var _ = fmt.Printf

/*
	Levenshtein Distance
	====================

	The Levenshtein Distance between two strings is the minimum number of
	single character insertions, deletions, or substitutions needed to
	convert one string into the other.  It is an "edit distance" string
	metric.

	* https://en.wikipedia.org/wiki/Levenshtein_distance
	* https://en.wikipedia.org/wiki/Edit_distance

	A function to compute the Levenshtein distance can be defined like so:

	L("", t) = len(t)			# base case
	L(s, "") = len(s)			# base case

	L(s + c, t + d) = 			# where c and d are characters

		let cost = 0 if c == d else 1 in
			min (
				L(s, t) + cost,		# substitution or characters match
				L(s, t + d) + 1,	# character inserted into s
				L(s + c, t) + 1,	# character deleted from s
			)

	This definition can be naively implemented as a recursive function:

*/

// -------------------------------------------
// ------------------------------------------- LevenshteinDistance_v1
// -------------------------------------------

/*
	A Naive Recursive Implementation of Levenshtein Distance
	========================================================

	The mathematical definition of Levenshtein Distance above can easily
	be translated into just about any programming language which supports
	recursion.  There is one caveat, though: this function will be
	unusably slow for all but the shortest strings.
*/

func LevenshteinDistance_v1(sc, td string) int {

	// Base case: L("", t) = len(t)
	if len(sc) == 0 {
		return len(td)
	}

	// Base case: L(s, "") = len(s)
	if len(td) == 0 {
		return len(sc)
	}

	// Split each string into prefix and final character

	// L(s + c, t + d) = ...

	s, c := sc[:len(sc) - 1], sc[len(sc) - 1]	// i.e. sc = s + c
	t, d := td[:len(td) - 1], td[len(td) - 1]	// i.e. td = t + d

	// Compute cost: let cost = 0 if c == d else 1 
	cost := 0;
	if c != d {
		cost = 1
	}

	// min (
	// 	L(s, t) + cost,		# substitution or characters match
	// 	L(s, t + d) + 1,	# character inserted into s
	// 	L(s + c, t) + 1,	# character deleted from s
	// )

	return min_int_3(
		LevenshteinDistance_v1(s, t) + cost,	// substitution or characters match
		LevenshteinDistance_v1(s, td) + 1,		// character inserted into s
		LevenshteinDistance_v1(sc, t) + 1,		// character deleted from s
	)
}

// -------------------------------------------

func min_int_3(a, b, c int) int {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
}

// -------------------------------------------
// -------------------------------------------
// -------------------------------------------

/*
	A More Efficient Approach Using Dynamic Programming
	===================================================

	The simple recursive definition above is really slow for strings of
	any significant size at all.  This is because the same intermediate
	values are computed over and over again.

	One approach for optimizing the algorigthm is to use *memoization*,
	where the results of functions calls are stored in a cache keyed
	by the function arguments.  This way each intermediate result only
	has to be calculated once.

	In the particular case of the Levenshtein Distance algorithm it
	turns out that we need *all* of the possible intermediate results,
	so another option is to simply compute all the intermediate results
	from the bottom up, and store them in a matrix for easy retrieval.
	Compared to implementing your own memoization scheme for a
	programming language that doesn't support it, this is quite a bit
	easier.  It's also very fast since there is minimal overhead.

	The pre-computing approach is a classic example of *dynamic
	programming*.

	The Levenshtein Distance Algorithm in Dynamic Programming Form
	--------------------------------------------------------------

	Recall the Levenshtein Distance algorithm from above:

	L("", t) = len(t)			# base case
	L(s, "") = len(s)			# base case

	L(s + c, t + d) = 			# where c and d are characters

		let cost = 0 if c == d else 1 in
			min (
				L(s, t) + cost,		# substitution or characters match
				L(s, t + d) + 1,	# character inserted into s
				L(s + c, t) + 1,	# character deleted from s
			)

	Now imagine that we want to store all the intermediate results in
	a matrix, which we'll call M.

	First, assume we are comparing the strings s and t, and len(s) = m
	and len(t) = n.

	+--------------------------------------------------------------------------
	|
	| Aside: An example showing *all* of the prefixes of a string
	| ===========================================================
	|
	| All of the prefixes of "cat" including the empty string and "cat" itself
	|
	| indexes  range notation  substring  characters
	| -------  --------------  ---------  ----------
	|          s[:0]           ""
	| 0        s[:1]           "c"        s[0] => "c"
	| 01       s[:2]           "ca"       s[0] => "c", s[1] => "a"
	| 012      s[:3]           "cat"      s[0] => "c", s[1] => "a", s[2] => "t"
	|
	+--------------------------------------------------------------------------

	The string s has m + 1 prefixes ranging from the empty string s[:0]
	all the way up to s itself, s[:m].  So for a string of length m,
	there are m + 1 prefixes including (pedantically) both the empty
	string and the entire source string.

	Likewise, the string t has n + 1 prefixes ranging from the empty
	string t[:0] to t[:n], t itself.

	We can re-express the base cases `L("", t) = len(t)` and `L(s, "") = len(s)`
	in terms of the prefix substrings:

		L(s[:0], t[:j]) = j, for 0 <= j <= m
		L(s[:i], t[:0]) = i, for 0 <= i <= n

	where s[:0] is the first zero characters of s (i.e. the empty string) and
	t[:j] is the first j characters of t.  Likewise s[:i] is the first i 
	characters of s, and t[:0] is the first zero characters of t (again, the
	empty string).

	Now imagine a new function M, where M(i, j) = L(s[:i], t[:j]).  The base
	cases now look like:

	M(0, j) = j, for 0 <= j <= m
	M(i, 0) = i, for 0 <= i <= n

	Similarly:

		L(s + c, t + d) = 				# where c and d are characters

			let cost = 0 if c == d else 1 in
				min (
					L(s, t) + cost,		# substitution or characters match
					L(s, t + d) + 1,	# character inserted into s
					L(s + c, t) + 1,	# character deleted from s
				)

	can be transformed first into:

		L(s[:i] + s[i], t[:j] + t[j]) =

			let cost = 0 if s[i] == t[j] else 1 in
				min (
					L(s[:i], t[:j]) + cost,			# substitution or characters match
					L(s[:i], t[:j + 1]) + 1,		# character inserted into s
					L(s[:i + 1], t[:j]) + 1			# character deleted from s
				)

	and then into:

		L(s[:i + 1], t[:j + 1]) =		# s[:i] + s[i] => s[:i + 1] and t[:j] + t[j] => t[:j + 1]

			o
			o
			o

	and again into:

		M(i + 1, j + 1) =

			let cost = 0 if s[i] == t[j] else 1 in
				min (
					M(i, j) + cost,			# substitution or characters match
					M(i, j + 1) + 1,		# character inserted into s
					M(i + 1, j) + 1			# character deleted from s
				)

	Finally, we recognize that M could simply be a matrix rather than a function:

		M[0, j] = j, for 0 <= j <= m
		M[i, 0] = i, for 0 <= i <= n

		M[i + 1, j + 1] =

			let cost = 0 if s[i] == t[j] else 1 in
				min (
					M[i, j] + cost,			# substitution or characters match
					M[i, j + 1] + 1,		# character inserted into s
					M[i + 1, j] + 1			# character deleted from s
				)

	However, unlike functions, matrices don't normally compute themselves, so we need
	a function to populate the matrix.  In pseudo-code, this would look something like:

	M = new Matrix[m + 1, n + 1]

	for j in 0..n + 1 do M[0, j] = j 	# last j is j = n
	for i in 1..m + 1 do M[i, 0] = i 	# we have already initialized M[0, 0], last i is i = m

	for i in 0..m do
		for j in 0..n do
			cost = 0 if s[i] == t[j] else 1
			M[i + 1, j + 1] = min(
				M[i, j] + cost,			# substitution or characters match
				M[i, j + 1] + 1,		# character inserted into s
				M[i + 1, j] + 1			# character deleted from s
			)

	The final value in the matrix is M[m, n], and is the Levenshtein Distance for our
	original strings s and t.

*/


// -------------------------------------------
// ------------------------------------------- LevenshteinDistance_v2
// -------------------------------------------

func LevenshteinDistance_v2(s, t string) int {

	m, n := len(s), len(t)

	// Go doesn't have natural two-dimensional arrays.  One option
	// is to pack the two-dimensional array into a one-D array.
	// The "offset" function abstracts out the offset calculation
	// so we can pretend that we really have a two-D array.
	matrix := make([]int, (m + 1) * (n + 1))	// number of rows * number of columns
	offset := func (i, j int) int { return i * (n + 1) + j }

	for j := 0; j < n + 1; j++ {
		matrix[offset(0, j)] = j
	}
	for i := 1; i < m + 1; i++ {
		matrix[offset(i, 0)] = i
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			var cost int
			if s[i] == t[j] {
				cost = 0
			} else {
				cost = 1
			}
			matrix[offset(i + 1, j + 1)] = min_int_3(
				matrix[offset(i, j)] + cost,
				matrix[offset(i, j + 1)] + 1,
				matrix[offset(i + 1, j)] + 1,
			)
		}
	}

	return matrix[offset(m, n)]
}

// -------------------------------------------
// ------------------------------------------- LevenshteinDistance_v3
// -------------------------------------------

// Just one change: Initialze the cells in the first column just before
// we need them.  This will make sense when you see the next version.

func LevenshteinDistance_v3(s, t string) int {

	m, n := len(s), len(t)

	// Go doesn't have natural two-dimensional arrays.  One option
	// is to pack the two-dimensional array into a one-D array.
	// The "offset" function abstracts out the offset calculation
	// so we can pretend that we really have a two-D array.
	matrix := make([]int, (m + 1) * (n + 1))	// number of rows * number of columns
	offset := func (i, j int) int { return i * (n + 1) + j }

	// Initialize the first row now.  However, we will wait and initialize
	// cells in the first column row-by-row, just before we need them.
	for j := 0; j < n + 1; j++ {
		matrix[offset(0, j)] = j
	}

	for i := 0; i < m; i++ {
		matrix[offset(i + 1, 0)] = i + 1 // yo
		for j := 0; j < n; j++ {
			var cost int
			if s[i] == t[j] {
				cost = 0
			} else {
				cost = 1
			}
			matrix[offset(i + 1, j + 1)] = min_int_3(
				matrix[offset(i, j)] + cost,
				matrix[offset(i, j + 1)] + 1,
				matrix[offset(i + 1, j)] + 1,
			)
		}
	}

	return matrix[offset(m, n)]
}

// -------------------------------------------
// ------------------------------------------- LevenshteinDistance_v4
// -------------------------------------------

// One more change, we only ever need to keep the last two rows in memory!

func LevenshteinDistance_v4(s, t string) int {

	m, n := len(s), len(t)

	// Strictly speaking, we don't need to save the entire matrix
	// of intermediate results, we only need the row right before
	// the one we are currently computing.
	rowCount := 2

	// Go doesn't have natural two-dimensional arrays.  One option
	// is to pack the two-dimensional array into a one-D array.
	// The "offset" function abstracts out the offset calculation
	// so we can pretend that we really have a two-D array.
	matrix := make([]int, rowCount * (n + 1))	// number of rows * number of columns
	offset := func (i, j int) int { return (i % rowCount) * (n + 1) + j }

	// Initialize the first row now.  However, we will wait and initialize
	// cells in the first column row-by-row, just before we need them.
	for j := 0; j < n + 1; j++ {
		matrix[offset(0, j)] = j
	}

	for i := 0; i < m; i++ {
		matrix[offset(i + 1, 0)] = i + 1 // yo
		for j := 0; j < n; j++ {
			var cost int
			if s[i] == t[j] {
				cost = 0
			} else {
				cost = 1
			}
			matrix[offset(i + 1, j + 1)] = min_int_3(
				matrix[offset(i, j)] + cost,
				matrix[offset(i, j + 1)] + 1,
				matrix[offset(i + 1, j)] + 1,
			)
		}
	}

	return matrix[offset(m, n)]
}

// -------------------------------------------
// ------------------------------------------- LevenshteinDistance_v5
// -------------------------------------------

func LevenshteinDistance_v5(s, t string) int {
	distance, alignment := Diff_v1(s, t)

	// --- display the alignment ---
	alignment.Dump(ComparableString(s), ComparableString(t), distance, SimpleStdoutLogger)

	return distance
}

// -------------------------------------------
// -------------------------------------------
// -------------------------------------------

func Diff_v1(s, t string) (distance int, alignment *Alignment) {

	alignment = new(Alignment)

	// --- compute the edit distance matrix

	m, n := len(s), len(t)

	// Go doesn't have natural two-dimensional arrays.  One option
	// is to pack the two-dimensional array into a one-D array.
	// The "offset" function abstracts out the offset calculation
	// so we can pretend that we really have a two-D array.
	matrix := make([]int, (m + 1) * (n + 1))	// number of rows * number of columns
	offset := func (i, j int) int { return i * (n + 1) + j }

	for j := 0; j < n + 1; j++ {
		matrix[offset(0, j)] = j
	}
	for i := 1; i < m + 1; i++ {
		matrix[offset(i, 0)] = i
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			var cost int
			if s[i] == t[j] {
				cost = 0
			} else {
				cost = 1
			}
			matrix[offset(i + 1, j + 1)] = min_int_3(
				matrix[offset(i, j)] + cost,
				matrix[offset(i, j + 1)] + 1,
				matrix[offset(i + 1, j)] + 1,
			)
		}
	}

	// --- extract an alignment from the computed matrix ---

	for i, j := m, n; i > 0 || j > 0; {

		var iNext, jNext int

		var link Link

		// We'll use "sIndex" and "tIndex" when referring to the "s" and "t" sequences,
		// and "i" and "j" when referring to coordinates into the computation matrix.
		// This makes the code a little easier to read.
		sIndex := i - 1
		tIndex := j - 1

		if i < 1 {
			link, iNext, jNext = Link{RightOnly, -1, tIndex}, 0, j - 1
		} else if j < 1 {
			link, iNext, jNext = Link{LeftOnly, sIndex, -1}, i - 1, 0
		} else {

			var cost int
			if s[i - 1] == t[j - 1] { cost = 0 } else { cost = 1 }

			a := matrix[offset(i - 1, j - 1)] + cost
			b := matrix[offset(i - 1, j)] + 1
			c := matrix[offset(i, j - 1)] + 1

			// Another readability improvement: Use boolean temporaries rather than inlining the expressions.  
			aIsOK := a <= b && a <= c
			bIsOK := b <= a && b <= c
			cIsOK := c <= a && c <= b

			if aIsOK {
				if cost == 0.0 {
					link, iNext, jNext = Link{Matching, sIndex, tIndex}, i - 1, j - 1
				} else {
					link, iNext, jNext = Link{Different, sIndex, tIndex}, i - 1, j - 1
				}
			} else if bIsOK {
				link, iNext, jNext = Link{LeftOnly, sIndex, -1}, i - 1, j
			} else if cIsOK {
				link, iNext, jNext = Link{RightOnly, -1, tIndex}, i, j - 1
			} else {
				panic("not reached")
			}
		}

		alignment.Links = append(alignment.Links, link)

		i, j = iNext, jNext
	}

	// The links are supposed to be in ascending order, but we've extracted them
	// in descending order, so now we need to reverse them.
	for low, high := 0, len(alignment.Links) - 1; low < high; low, high = low + 1, high - 1 {
		alignment.Links[low], alignment.Links[high] = alignment.Links[high], alignment.Links[low]
	}

	return matrix[offset(m, n)], alignment
}

// -------------------------------------------
// ------------------------------------------- LevenshteinDistance_v6
// -------------------------------------------

func LevenshteinDistance_v6(s, t string) int {
	distance, alignment := Diff_v2(ComparableString(s), ComparableString(t))

	// --- display the alignment ---
	_ = alignment
	// alignment.dump(s, t, distance)

	return int(distance)
}

// -------------------------------------------
// -------------------------------------------
// -------------------------------------------

func Diff_v2(s, t ComparableSequence) (distance float32, alignment *Alignment) {

	alignment = new(Alignment)

	// --- compute the edit distance matrix

	m, n := s.Length(), t.Length()

	// Go doesn't have natural two-dimensional arrays.  One option
	// is to pack the two-dimensional array into a one-D array.
	// The "offset" function abstracts out the offset calculation
	// so we can pretend that we really have a two-D array.
	matrix := make([]float32, (m + 1) * (n + 1))	// number of rows * number of columns
	offset := func (i, j int) int { return i * (n + 1) + j }

	for j := 0; j < n + 1; j++ {
		matrix[offset(0, j)] = float32(j)
	}
	for i := 1; i < m + 1; i++ {
		matrix[offset(i, 0)] = float32(i)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			cost := s.GetItemAt(i).Compare(t.GetItemAt(j))
			matrix[offset(i + 1, j + 1)] = min_float32_3(
				matrix[offset(i, j)] + cost,
				matrix[offset(i, j + 1)] + 1,
				matrix[offset(i + 1, j)] + 1,
			)
		}
	}

	// --- extract an alignment from the computed matrix ---

	for i, j := m, n; i > 0 || j > 0; {

		var iNext, jNext int

		var link Link

		// We'll use "sIndex" and "tIndex" when referring to the "s" and "t" sequences,
		// and "i" and "j" when referring to coordinates into the computation matrix.
		// This makes the code a little easier to read.
		sIndex := i - 1
		tIndex := j - 1

		if i < 1 {
			link, iNext, jNext = Link{RightOnly, -1, tIndex}, 0, j - 1
		} else if j < 1 {
			link, iNext, jNext = Link{LeftOnly, sIndex, -1}, i - 1, 0
		} else {

			cost := s.GetItemAt(i - 1).Compare(t.GetItemAt(j - 1))

			a := matrix[offset(i - 1, j - 1)] + cost
			b := matrix[offset(i - 1, j)] + 1
			c := matrix[offset(i, j - 1)] + 1

			// Another readability improvement: Use boolean temporaries rather than inlining the expressions.  
			aIsOK := a <= b && a <= c
			bIsOK := b <= a && b <= c
			cIsOK := c <= a && c <= b

			if aIsOK {
				if cost == 0.0 {
					link, iNext, jNext = Link{Matching, sIndex, tIndex}, i - 1, j - 1
				} else {
					link, iNext, jNext = Link{Different, sIndex, tIndex}, i - 1, j - 1
				}
			} else if bIsOK {
				link, iNext, jNext = Link{LeftOnly, sIndex, -1}, i - 1, j
			} else if cIsOK {
				link, iNext, jNext = Link{RightOnly, -1, tIndex}, i, j - 1
			} else {
				panic("not reached")
			}
		}

		alignment.Links = append(alignment.Links, link)

		i, j = iNext, jNext
	}

	// The links are supposed to be in ascending order, but we've extracted them
	// in descending order, so now we need to reverse them.
	for low, high := 0, len(alignment.Links) - 1; low < high; low, high = low + 1, high - 1 {
		alignment.Links[low], alignment.Links[high] = alignment.Links[high], alignment.Links[low]
	}

	return matrix[offset(m, n)], alignment
}

// -------------------------------------------

func min_float32_3(a, b, c float32) float32 {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
}

