package strings

import "strings"

const (
	markdown_marker = "```"
)

type MetaString struct {
	Markdown bool
	Content  string
}

func ExtractMarkdownSections(content string) []MetaString {
	sections := []MetaString{}
	splits := SplitWithoutCopy(content, markdown_marker)

	for i := range splits {
		sections = append(sections, MetaString{Markdown: i%2 == 1, Content: splits[i]})
	}
	return sections
}

func SplitWithoutCopy(s, sep string) []string {
	separatorIndices := []int{}
	startIndex := 0
	for {
		index := strings.Index(s[startIndex:], sep)
		if index == -1 {
			break
		}
		separatorIndices = append(separatorIndices, startIndex+index)
		startIndex += index + len(sep)
	}

	substrings := make([]string, len(separatorIndices)+1)
	startIndex = 0
	for i, endIndex := range separatorIndices {
		substrings[i] = s[startIndex:endIndex]
		startIndex = endIndex + len(sep)
	}
	substrings[len(separatorIndices)] = s[startIndex:]

	return substrings
}
