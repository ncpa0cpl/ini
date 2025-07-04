package ini

import (
	"strings"
)

const (
	parseStepLookup = iota
	parseStepComment
	parseStepFieldComment
	parseStepSection
	parseStepKey
	parseStepValue
)

type docOrSection interface {
	Del(key string)
	Get(key string) string
	GetBool(key string) (bool, error)
	GetFloat(key string) (float64, error)
	GetInt(key string) (int64, error)
	GetUint(key string) (uint64, error)
	Set(key string, value string)
	SetBool(key string, value bool)
	SetFieldComment(fieldKey string, value string)
	SetFloat(key string, value float64)
	SetInt(key string, value int64)
	SetUint(key string, value uint64)
	AddComment(value string)
	AddHashComment(value string)
	AddWhiteLine()
	Section(name string) *IniSection
	ToString() string
	addParsedSection(name string) *IniSection
}

func Parse(content string) *IniDoc {
	var key string

	step := parseStepLookup
	escaped := false
	commentType := ';'
	buff := make([]rune, 0, 16)

	doc := NewDoc()

	var currentDoc docOrSection
	currentDoc = doc

	for idx, char := range content {
		if char == '\\' && !escaped {
			escaped = true
			continue
		}

		switch step {
		case parseStepLookup:
			if !escaped {
				switch char {
				case '[':
					step = parseStepSection
					continue
				case ';':
					step = parseStepComment
					commentType = ';'
					continue
				case '#':
					step = parseStepComment
					commentType = '#'
					continue
				case '\n':
					if idx == 0 || content[idx-1] == '\n' {
						currentDoc.AddWhiteLine()
					}
					continue
				case ' ':
					continue
				}
			} else {
				escaped = false
			}
			step = parseStepKey
			buff = append(buff, char)
		case parseStepKey:
			switch char {
			case '=':
				if !escaped {
					key = strings.Trim(string(buff), " ")
					buff = make([]rune, 0, 16)
					step = parseStepValue
				} else {
					buff = append(buff, char)
				}
			case '\n':
				buff = make([]rune, 0, 16)
				step = parseStepLookup
			default:
				buff = append(buff, char)
			}
			escaped = false
		case parseStepValue:
			if !escaped {
				switch char {
				case ';', '#':
					currentDoc.Set(key, strings.Trim(string(buff), " "))
					buff = make([]rune, 0, 16)
					step = parseStepFieldComment
					continue
				case '\n':
					currentDoc.Set(key, strings.Trim(string(buff), " "))
					buff = make([]rune, 0, 16)
					key = ""
					step = parseStepLookup
					continue
				}
			} else {
				escaped = false
				if char == 'N' || char == 'n' {
					buff = append(buff, '\n')
					continue
				}
			}
			buff = append(buff, char)
		case parseStepSection:
			switch char {
			case ']':
				if !escaped {
					currentDoc = doc.addParsedSection(strings.Trim(string(buff), " "))
					buff = make([]rune, 0, 16)
					step = parseStepLookup
				} else {
					buff = append(buff, char)
				}
			case '\n':
				buff = make([]rune, 0, 16)
				step = parseStepLookup
			default:
				buff = append(buff, char)
			}
			escaped = false
		case parseStepFieldComment:
			if char == '\n' && !escaped {
				currentDoc.SetFieldComment(key, strings.Trim(string(buff), " "))
				buff = make([]rune, 0, 16)
				key = ""
				step = parseStepLookup
			} else {
				escaped = false
				buff = append(buff, char)
			}
		case parseStepComment:
			if char == '\n' && !escaped {
				if commentType == ';' {
					currentDoc.AddComment(strings.Trim(string(buff), " "))
				} else {
					currentDoc.AddHashComment(strings.Trim(string(buff), " "))
				}
				buff = make([]rune, 0, 16)
				step = parseStepLookup
			} else {
				escaped = false
				buff = append(buff, char)
			}
		}
	}

	if key != "" && len(buff) > 0 {
		if step == parseStepFieldComment {
			currentDoc.SetFieldComment(key, strings.Trim(string(buff), " "))
		} else {
			currentDoc.Set(key, strings.Trim(string(buff), " "))
		}
	}

	return doc
}
