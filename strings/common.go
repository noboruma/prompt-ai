package strings

import "strings"

type MetaString struct {
	Markdown bool
	Content  string
}

func ExtractMarkdownSections(content string) []MetaString {
	sections := []MetaString{}
	splits := strings.Split(content, "```")

	for i := range splits {
		sections = append(sections, MetaString{Markdown: i%2 == 1, Content: splits[i]})
	}
	return sections
}
