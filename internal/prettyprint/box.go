package prettyprint

// asBox prints arbitrary line sof text into a box.
func asBox(lines []string) string {
	longestLength := 0
	for _, l := range lines {
		length := len([]rune(l))
		if length > longestLength {
			longestLength = length
		}
	}

	topLine := "┌─"
	bottomLine := "└─"
	for i := 1; i <= longestLength; i++ {
		topLine += "─"
		bottomLine += "─"
	}
	topLine += "─┐\n"
	bottomLine += "─┘\n"

	finalLines := make([]string, 0, len(lines))
	for i, l := range lines {
		length := len([]rune(l))
		if i == 0 {
			l = asBold(l)
		}
		missingChars := longestLength - length
		line := "│ "
		line += l
		for j := 1; j <= missingChars; j++ {
			line += " "
		}
		line += " │\n"
		finalLines = append(finalLines, line)
	}

	finalString := topLine
	for _, l := range finalLines {
		finalString += l
	}
	finalString += bottomLine

	return finalString
}
