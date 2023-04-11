package strings

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
