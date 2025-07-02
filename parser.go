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

type DocOrSection interface {
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
	ToString() string
}

func Parse(content string) *IniDoc {
	var key string

	step := ParseStepLookup
	escaped := false
	commentType := ';'
	buff := make([]rune, 0, 16)

	doc := NewDoc()

	var currentDoc DocOrSection
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
					commentType = ';'
					continue
				case '#':
					step = ParseStepComment
					commentType = '#'
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
			switch char {
			case '=':
				if !escaped {
					key = strings.Trim(string(buff), " ")
					buff = make([]rune, 0, 16)
					step = ParseStepValue
				} else {
					buff = append(buff, char)
				}
			case '\n':
				buff = make([]rune, 0, 16)
				step = ParseStepLookup
			default:
				buff = append(buff, char)
			}
			escaped = false
		case ParseStepValue:
			if !escaped {
				switch char {
				case ';', '#':
					currentDoc.Set(key, strings.Trim(string(buff), " "))
					buff = make([]rune, 0, 16)
					step = ParseStepFieldComment
					continue
				case '\n':
					currentDoc.Set(key, strings.Trim(string(buff), " "))
					buff = make([]rune, 0, 16)
					key = ""
					step = ParseStepLookup
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
		case ParseStepSection:
			switch char {
			case ']':
				if !escaped {
					currentDoc = doc.Section(strings.Trim(string(buff), " "))
					buff = make([]rune, 0, 16)
					step = ParseStepLookup
				} else {
					buff = append(buff, char)
				}
			case '\n':
				buff = make([]rune, 0, 16)
				step = ParseStepLookup
			default:
				buff = append(buff, char)
			}
			escaped = false
		case ParseStepFieldComment:
			if char == '\n' && !escaped {
				currentDoc.SetFieldComment(key, strings.Trim(string(buff), " "))
				buff = make([]rune, 0, 16)
				key = ""
				step = ParseStepLookup
			} else {
				escaped = false
				buff = append(buff, char)
			}
		case ParseStepComment:
			if char == '\n' && !escaped {
				if commentType == ';' {
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

	if key != "" && len(buff) > 0 {
		if step == ParseStepFieldComment {
			currentDoc.SetFieldComment(key, strings.Trim(string(buff), " "))
		} else {
			currentDoc.Set(key, strings.Trim(string(buff), " "))
		}
	}

	return doc
}
