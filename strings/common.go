package strings

import "strings"

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
