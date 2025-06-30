package ini

import "strings"

const (
	ParseStepLookup = iota
	ParseStepComment
	ParseStepFieldComment
	ParseStepSection
	ParseStepKey
	ParseStepValue
)

type docOrSection interface {
	Set(key string, value string)
	SetBool(key string, value bool)
	SetFieldComment(fieldKey string, value string)
	SetFloat(key string, value float64)
	SetInt(key string, value int64)
	SetUint(key string, value uint64)
	AddComment(value string)
	AddHashComment(value string)
}

func Parse(content string) *IniDoc {
	var key string
	var value string

	step := ParseStepLookup
	escaped := false
	commentType := ";"
	buff := make([]rune, 0, 16)

	doc := NewDoc()

	var currentDoc docOrSection
	currentDoc = doc

	for _, char := range content {
		if char == '\\' && !escaped {
			escaped = true
			continue
		}

		switch step {
		case ParseStepLookup:
			if !escaped {
				switch char {
				case '[':
					step = ParseStepSection
					continue
				case ';':
					step = ParseStepComment
					commentType = ";"
					continue
				case '#':
					step = ParseStepComment
					commentType = "#"
					continue
				case ' ', '\n':
					continue
				}
			} else {
				escaped = false
			}
			step = ParseStepKey
			buff = append(buff, char)
		case ParseStepKey:
			if char == '=' && !escaped {
				key = strings.Trim(string(buff), " ")
				buff = make([]rune, 0, 16)
				step = ParseStepValue
			} else {
				escaped = false
				buff = append(buff, char)
			}
		case ParseStepValue:
			if !escaped {
				switch char {
				case ';', '#':
					value = strings.Trim(string(buff), " ")
					buff = make([]rune, 0, 16)
					currentDoc.Set(key, value)
					step = ParseStepFieldComment
					continue
				case '\n':
					value = strings.Trim(string(buff), " ")
					buff = make([]rune, 0, 16)
					currentDoc.Set(key, value)
					key = ""
					value = ""
					step = ParseStepLookup
					continue
				}
			} else {
				escaped = false
			}
			buff = append(buff, char)
		case ParseStepSection:
			if char == ']' && !escaped {
				currentDoc = doc.Section(strings.Trim(string(buff), " "))
				buff = make([]rune, 0, 16)
				step = ParseStepLookup
			} else {
				escaped = false
				buff = append(buff, char)
			}
		case ParseStepFieldComment:
			if char == '\n' && !escaped {
				currentDoc.SetFieldComment(key, strings.Trim(string(buff), " "))
				buff = make([]rune, 0, 16)
				key = ""
				value = ""
				step = ParseStepLookup
			} else {
				escaped = false
				buff = append(buff, char)
			}
		case ParseStepComment:
			if char == '\n' && !escaped {
				if commentType == ";" {
					currentDoc.AddComment(strings.Trim(string(buff), " "))
				} else {
					currentDoc.AddHashComment(strings.Trim(string(buff), " "))
				}
				buff = make([]rune, 0, 16)
				step = ParseStepLookup
			} else {
				escaped = false
				buff = append(buff, char)
			}
		}
	}

	if key != "" {
		currentDoc.Set(key, value)
		if step == ParseStepFieldComment && len(buff) > 0 {
			currentDoc.SetFieldComment(key, strings.Trim(string(buff), " "))
		}
	}

	return doc
}
