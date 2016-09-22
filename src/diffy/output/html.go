package output

import (
	"fmt"
	"html"
	"path/filepath"
	"strconv"
	"strings"

	"diffy/diff"
)

// ------------------------------------------- type SourceLinesRec

type SourceLinesRec struct {
	Lines diff.ComparableLines
	FilePath string
}

func NewSourceLinesRec(lines diff.ComparableLines, filePath string) *SourceLinesRec {
	return &SourceLinesRec{Lines: lines, FilePath: filePath}
}

func (source *SourceLinesRec) GetFileName() string {
	return filepath.Base(source.FilePath)
}

func (source *SourceLinesRec) GetAbsoluteFilePath() string {
	absolutePath, err := filepath.Abs(source.FilePath)
	if err != nil {
		return source.FilePath
	}
	return absolutePath
}

// ------------------------------------------- type CssStyle
//
// CssStyle records represent a CSS "style", which for our purposes is just
// a list of CSS properties and their values, with each property/value pair
// represented as a single string.  Multiple CssStyle records can be 
// combined into a single inline HTML "style" attribute.  In the future
// there may be support for generating a proper CSS style sheet.

type CssStyle struct{
	className string
	properties []string
}

func MakeCssStyle(className string, properties ...string) CssStyle {
	return CssStyle{
		className:className,
		properties:properties,
	}
}

func ConcatCssStyles(styles ...CssStyle) string {
	var properties []string
	for _, style := range styles {
		properties = append(properties, style.properties...)
	}
	return strings.Join(properties, ";")
}

func (style CssStyle) when(cond bool) CssStyle {
	if cond {
		return style
	} else {
		return nullStyle
	}
}

// ------------------------------------------- CSS style definitions

// ........................................... null style

var nullStyle CssStyle = MakeCssStyle("null")

// ........................................... title headings table and friends

var titleHeadingsTableStyle CssStyle = MakeCssStyle("title-headings-table",
	"width: 100%",
	"margin-bottom: 0px",
	"border-left: solid #696969 2px",
	"border-right: solid #696969 2px",
	"border-collapse: collapse",
	"border-spacing: 0px",
	"table-layout: fixed",
	"color: white",
	"font-family: monospace",
)

var titleHeadingBoxStyle CssStyle = MakeCssStyle("title-heading-box",
	"border: solid black 1px",
	"background-color: #4682B4",
)

var headingTitleStyle CssStyle = MakeCssStyle("heading-title",
	"padding: 5px",
	"font-size: 20pt",
	"font-weight: bold",
)

var headingSubtitleStyle CssStyle = MakeCssStyle("heading-subtitle",
	"padding: 5px",
	"font-size: 12pt",
	"font-style: italic",
)

// ........................................... two line diff table and friends

var twoLineDiffStyle CssStyle = MakeCssStyle("two-line-diff",
	"width: 100%",
	"border-collapse: collapse",
	"border-spacing: 0px",
	"table-layout: fixed",
)

var lineNumStyle CssStyle = MakeCssStyle("line-num",
	"width: 5ex",
	"padding-right: 5px",
	"background-color: #EEE",
	"white-space: pre",
	"font-family: monospace",
	"font-size: 9pt",
	"text-align: right",
)

var codeLineStyle CssStyle = MakeCssStyle("code-line",
	"overflow: hidden",
	"text-overflow: ellipsis",
	"padding-left: 5px",
	"padding-right: 5px",
	"font-family: monospace",
	"font-size: 9pt",
	"white-space: pre",
)

var codeLineLinesDifferStyle CssStyle = MakeCssStyle("code-line-lines-differ",
	"background-color: #FFFFE0",
)

var codeLineOnlyOneStyle CssStyle = MakeCssStyle("code-line-only-one",
	"background-color: #FFEC8B",
)

var codeLineNoneStyle CssStyle = MakeCssStyle("code-line-none",
	"background-color: #F0F0F0",
)

var twoLineDiffGutterStyle CssStyle = MakeCssStyle("two-line-diff-gutter",
	"height: 3px",
	"width: 1px",
	"border-left: solid black 2px",
	"border-right: solid black 2px",
)

var codeRunDifferentStyle CssStyle = MakeCssStyle("code-run-different",
	"background-color: lightgreen",
)

// ------------------------------------------- GenerateHtmlDiffPage
//
func GenerateHtmlDiffPage(alignment *diff.Alignment, leftSource, rightSource *SourceLinesRec) {

	// Re-jigger the alignment to make it more suitable for display.
	alignment = alignment.RealignUsingThreshold(leftSource.Lines, rightSource.Lines, 0.4)

	// Print the page prologue.
	fmt.Println("<!DOCTYPE html>")
	fmt.Println("<html>")
	fmt.Println("	<head>")
	fmt.Println("		<title>Diff</title>")
	fmt.Println("")
	fmt.Println("		<meta charset=\"utf-8\"/>")
	fmt.Println("	</head>")
	fmt.Println("	<body>")

	// Print the heading.
	fmt.Println("")

	fmt.Printf("		%s\n", generateStartTag("table", titleHeadingsTableStyle))
	fmt.Printf("			%s\n", generateStartTag("tr"))
	fmt.Printf("				%s\n", generateStartTag("td", titleHeadingBoxStyle))
	fmt.Printf("					%s\n", generateElement("div", leftSource.GetFileName(), headingTitleStyle))
	fmt.Printf("					%s\n", generateElement("div", leftSource.GetAbsoluteFilePath(), headingSubtitleStyle))
	fmt.Printf("				%s\n", generateEndTag("td"))
	fmt.Printf("				%s\n", generateElement("td", "", twoLineDiffGutterStyle))
	fmt.Printf("				%s\n", generateStartTag("td", titleHeadingBoxStyle))
	fmt.Printf("					%s\n", generateElement("div", rightSource.GetFileName(), headingTitleStyle))
	fmt.Printf("					%s\n", generateElement("div", rightSource.GetAbsoluteFilePath(), headingSubtitleStyle))
	fmt.Printf("				%s\n", generateEndTag("td"))
	fmt.Printf("			%s\n", generateEndTag("tr"))
	fmt.Printf("		%s\n", generateEndTag("table"))
	fmt.Println("")

	// Generate an empty initial "code-line" table to provide some extra spacing.
	fmt.Printf("		%s\n", generateStartTag("table", twoLineDiffStyle))
	fmt.Printf("			%s\n", generateStartTag("tr"))
	fmt.Printf("				%s\n", generateElement("td", "", lineNumStyle))
	fmt.Printf("				%s\n", generateElement("td", "", codeLineStyle))
	fmt.Printf("				%s\n", generateElement("td", "", twoLineDiffGutterStyle))
	fmt.Printf("				%s\n", generateElement("td", "", codeLineStyle))
	fmt.Printf("				%s\n", generateElement("td", "", lineNumStyle))
	fmt.Printf("			%s\n", generateEndTag("tr"))
	fmt.Printf("		%s\n", generateEndTag("table"))
	fmt.Println("")

	// For each link in the alignment generate a side-by-side diff of the corresponding
	// pair of lines.  We will just use blank lines when one line is missing.
	for _, link := range alignment.Links {

		// Figure out what type of link we've got.
		var leftItem, rightItem diff.Comparable = nil, nil
		switch link.LinkType {
		case diff.Matching, diff.Different:
			leftItem, rightItem = leftSource.Lines[link.LeftIndex], rightSource.Lines[link.RightIndex]
		case diff.LeftOnly:
			leftItem = leftSource.Lines[link.LeftIndex]
		case diff.RightOnly:
			rightItem = rightSource.Lines[link.RightIndex]
		default:
			panic("not reached")
		}

		// Generate the HTML for the left and right lines.
		leftHtml, rightHtml := "", ""
		if link.LinkType == diff.Different {
			leftHtml, rightHtml = generateLineHtml(leftItem.(*diff.TextLine).Text, rightItem.(*diff.TextLine).Text)
		} else {
			if leftItem != nil {
				leftHtml = html.EscapeString(leftItem.(*diff.TextLine).Text)
			}
			if rightItem != nil {
				rightHtml = html.EscapeString(rightItem.(*diff.TextLine).Text)
			}
		}

		// Figure out the appropriate styles for the left and right lines.
		leftLineStyle := []CssStyle{
			codeLineStyle,
			codeLineLinesDifferStyle.when(link.LinkType == diff.Different),
			codeLineOnlyOneStyle.when(link.LinkType == diff.LeftOnly),
			codeLineNoneStyle.when(leftItem == nil),
		}
		rightLineStyle := []CssStyle{
			codeLineStyle,
			codeLineLinesDifferStyle.when(link.LinkType == diff.Different),
			codeLineOnlyOneStyle.when(link.LinkType == diff.RightOnly),
			codeLineNoneStyle.when(rightItem == nil),
		}

		// Line numbers.  Remember that slice indexes start from zero, but line numbers start from 1!
		leftLineNumHtml, rightLineNumHtml := "", ""
		if link.LeftIndex >= 0 {
			leftLineNumHtml = strconv.FormatInt(int64(link.LeftIndex + 1), 10)
		}
		if link.RightIndex >= 0 {
			rightLineNumHtml = strconv.FormatInt(int64(link.RightIndex + 1), 10)
		}

		// Output the HTML for these two lines.
		fmt.Printf("		%s\n", generateStartTag("table", twoLineDiffStyle))
		fmt.Printf("			%s\n", generateStartTag("tr"))
		fmt.Printf("				%s\n", generateElement("td", leftLineNumHtml, lineNumStyle))
		fmt.Printf("				%s\n", generateElement("td", leftHtml, leftLineStyle...))
		fmt.Printf("				%s\n", generateElement("td", "", twoLineDiffGutterStyle))
		fmt.Printf("				%s\n", generateElement("td", rightHtml, rightLineStyle...))
		fmt.Printf("				%s\n", generateElement("td", rightLineNumHtml, lineNumStyle))
		fmt.Printf("			%s\n", generateEndTag("tr"))
		fmt.Printf("		%s\n", generateEndTag("table"))
	}
	fmt.Println("")

	// Generate an empty final "code-line" table to provide some extra spacing.
	fmt.Printf("		%s\n", generateStartTag("table", twoLineDiffStyle))
	fmt.Printf("			%s\n", generateStartTag("tr"))
	fmt.Printf("				%s\n", generateElement("td", "", lineNumStyle))
	fmt.Printf("				%s\n", generateElement("td", "", codeLineStyle))
	fmt.Printf("				%s\n", generateElement("td", "", twoLineDiffGutterStyle))
	fmt.Printf("				%s\n", generateElement("td", "", codeLineStyle))
	fmt.Printf("				%s\n", generateElement("td", "", lineNumStyle))
	fmt.Printf("			%s\n", generateEndTag("tr"))
	fmt.Printf("		%s\n", generateEndTag("table"))
	fmt.Println("")

	// Print the page epilogue.
	fmt.Println("	</body>")
	fmt.Println("</html>")
}

// ------------------------------------------- generateLineHtml
//
// Generate HTML which highlights the differences between two different but similar lines.
func generateLineHtml(leftLine, rightLine string) (string, string) {

	// Generate a diff for the two lines.
	leftLineRunes, rightLineRunes := diff.MakeComparableString(leftLine), diff.MakeComparableString(rightLine)
	_, alignment := diff.Diff_v2(leftLineRunes, rightLineRunes)

	// Use the "alignment" generated above to generate HTML which highlights the differences.
	leftRunPositions, rightRunPositions := findAlternatingRunPositions(alignment, diff.Matching)
	leftSpansHtml := constructEvenOddSpans(leftLineRunes, leftRunPositions, nullStyle, codeRunDifferentStyle)
	rightSpansHtml := constructEvenOddSpans(rightLineRunes, rightRunPositions, nullStyle, codeRunDifferentStyle)

	return leftSpansHtml, rightSpansHtml
}

// ------------------------------------------- findAlternatingRunPositions
//
// Based on the provided alignment and link type, generate "run positions" (one set each) for the
// left and right sequences that the diff was generated from.  The run positions can then be used
// to split each source sequence up into alternating runs for display.
//
// Notes:
// - the diff is split up into "left links" and "right links", which are then treated separately
// - the left links will have links for every item in the left source sequence, in order
// - the right links will have links for every item in the right source sequence, in order
// - even runs (0, 2, 4, etc.) will be runs where all the links match the specified link type
// - odd runs (1, 3, 5, etc.) will be runs where all the links do *not* match the link type
// - the first (zero'th!) run may be empty; all other runs should be non-empty
// - the final run position in the left run positions will be len(<left-sequence>)
// - the final run position in the right run positions will be len(<right-sequence>)
//
func findAlternatingRunPositions(alignment *diff.Alignment, linkType diff.LinkType) ([]int, []int) {

	findRunPositions := func (links []diff.Link) []int {
		runPositions := []int{0}
		prevLinkIsType := true
		for index, link := range links {
			currLinkIsType := link.LinkType == linkType
			if currLinkIsType != prevLinkIsType {
				runPositions = append(runPositions, index)
			}
			prevLinkIsType = currLinkIsType
		}
		runPositions = append(runPositions, len(links))
		return runPositions
	}

	leftLinks := getLeftLinks(alignment)
	rightLinks := getRightLinks(alignment)
	return findRunPositions(leftLinks), findRunPositions(rightLinks)
}

// ------------------------------------------- getLeftlinks

func getLeftLinks(alignment *diff.Alignment) []diff.Link {
	var leftLinks []diff.Link
	for _, link := range alignment.Links {
		if link.LeftIndex >= 0 {
			leftLinks = append(leftLinks, link)
		}
	}
	return leftLinks
}

// ------------------------------------------- getRightLinks

func getRightLinks(alignment *diff.Alignment) []diff.Link {
	var rightLinks []diff.Link
	for _, link := range alignment.Links {
		if link.RightIndex >= 0 {
			rightLinks = append(rightLinks, link)
		}
	}
	return rightLinks
}

// ------------------------------------------- constructEvenOddSpans
//
// Convert the literal text (or a subset thereof) in "runes" into HTML, where each "run" is
// represented as a single SPAN element, and where even spans are styled with "evenStyle" and
// odd runs are styled with "oddStyle".
//
// Notes:
// - each run position denotes the *beginning* of a run
// - run positions should be in ascending order
// - the final run position should be the position *after* the last rune of the last run
// - when the runs cover the whole rune slice, the first run position will be 0
// - when the runs cover the whole rune slice, the last run position will be len(runes)
//
func constructEvenOddSpans(runes []rune, runPositions []int, evenStyle, oddStyle CssStyle) string {
	var spansHtml []string
	for i := 0; i < len(runPositions) - 1; i++ {	// note: last iteration is i = len(runPositions) - 2
		runIsEven := i % 2 == 0
		runIsOdd := !runIsEven
		runStartIndex := runPositions[i + 0]
		runEndIndex := runPositions[i + 1]
		spanText := runes[runStartIndex:runEndIndex]
		spanTextEscaped := html.EscapeString(string(spanText))
		span := generateElement("span", spanTextEscaped, evenStyle.when(runIsEven), oddStyle.when(runIsOdd))
		spansHtml = append(spansHtml, span)
	}
	return strings.Join(spansHtml, "")
}

// ------------------------------------------- generateElement
//
// generateElement("div" ...) => "<div>...</div>" or "<div style='...'>...</div>"
// This function will generate no additional newlines, although the body may
// contain newlines which will be retained.
func generateElement(tagName string, body string, styles ...CssStyle) string {
	return generateStartTag(tagName, styles...) + body + generateEndTag(tagName)
}

// ------------------------------------------- generateStartTag
//
// generateStartTag("div" ...) => "<div>" or "<div style='...'>" as appropriate,
// depending on whether any styles are generated or not.
func generateStartTag(tagName string, styles ...CssStyle) string {

	startTagText := "<" + tagName

	stylePropertyText := ConcatCssStyles(styles...)
	if stylePropertyText != "" {
		startTagText += " style='" + stylePropertyText + "'"
	}

	return startTagText + ">"
}

// ------------------------------------------- generateEndTag
//
// generateEndTag("div") => "</div>"
func generateEndTag(tagName string) string {
	return "</" + tagName + ">"
}
