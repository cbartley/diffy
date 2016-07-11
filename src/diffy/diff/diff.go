package diff

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
	matrix := make([]int, (m + 1) * (n + 1))	// number of rows * number of columns

	for j := 0; j < n + 1; j++ {
		matrix[0 * (n + 1) + j] = j
	}
	for i := 1; i < m + 1; i++ {
		matrix[i * (n + 1) + 0] = i
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			var cost int
			if s[i] == t[j] {
				cost = 0
			} else {
				cost = 1
			}
			matrix[(i + 1) * (n + 1) + (j + 1)] = min_int_3(
				matrix[i * (n + 1) + j] + cost,
				matrix[i * (n + 1) + (j + 1)] + 1,
				matrix[(i + 1) * (n + 1) + j] + 1,
			)
		}
	}

	return matrix[m * (n + 1) + n]
}




