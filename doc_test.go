package ini_test

import (
	"testing"

	"github.com/ncpa0cpl/ini"
	"github.com/stretchr/testify/assert"
)

func TestDoc(t *testing.T) {
	assert := assert.New(t)

	doc := ini.NewDoc()
	doc.AddComment("this is a comment")
	doc.Set("str", "hello world")
	doc.SetBool("bl", true)
	doc.SetBool("bl2", false)
	doc.AddComment("numbers:")
	doc.SetFloat("flt", 1.2345)
	doc.SetInt("int", 420)
	doc.SetFieldComment("bl", "this is a boolean")

	stringifiedDoc := doc.ToString()

	expectedResult := `; this is a comment
str=hello world
bl=true ;this is a boolean
bl2=false
; numbers:
flt=1.2345
int=420
`
	assert.Equal(expectedResult, stringifiedDoc)
}

func TestDocSection(t *testing.T) {
	assert := assert.New(t)

	doc := ini.NewDoc()
	doc.Set("top", "lorem ipsum")

	section1 := doc.Section("MySection")
	section1.SetFloat("number", 4.2069)
	section1.Set("foobar", "bazquux")
	section1.SetBool("is", true)

	section2 := doc.Section("AnotherSection")
	section2.SetBool("is", false)
	section2.Set("bazquux", "foobar")
	section2.SetInt("number", -999)

	stringifiedDoc := doc.ToString()

	expectedResult := `top=lorem ipsum

[MySection]
number=4.2069
foobar=bazquux
is=true

[AnotherSection]
is=false
bazquux=foobar
number=-999
`
	assert.Equal(expectedResult, stringifiedDoc)
}
