
#  INI Parser & Reader Writer Library

## Introduction

The INI Parser & Reader Writer Library is a fast and easy-to-use library for parsing and manipulating INI files in the Go programming language. It provides functionality to read INI files from both strings and files, and offers options to marshal and unmarshal INI data into Go structs, and write data back to files.

## Features
* **Read by string**: The library allows you to parse INI data stored in a string.
* **Read by file**: You can also read INI data directly from a file.
* **Unmarshal to Struct**: It provides the ability to map INI data to Go structs, making it convenient to work with structured data.
* **Marshal from Struct**: You can easily convert Go struct into INI data using the library's marshal functionality.
* **Write to File**: The library allows you to write INI data back to files.
* **Custom marshaling**: MarshalINI and UnmarshalINI methods, similar to the [custom JSON marshaling](https://pkg.go.dev/encoding/json#example-package-CustomMarshalJSON)
* **Subsections**: Easily access and add subsections

## Table of Contents

1. [Installation](#installation)
2. [Document parsing, reading and writing](#document-parsing-reading-and-writing)
3. [Unmarshal Struct](#unmarshal-struct)
4. [Marshal Struct](#marshal-struct)
5. [Sections](#sections)
6. [Subsections](#subsections)
7. [Custom Marshal/Unmarshal](#custom-marshalunmarshal)
8. [Parse File](#parse-file)
9. [Write File](#write-file)

## Installation

```shell
go get github.com/ncpa0cpl/ini
```

## Document parsing, reading and writing

```go
iniFile := `
topLevelValue=foo

[section]
k=v

[section1]
k1=v1
k2=1
k3=3.5
k4=0.0.0.0
`

doc := ini.Parse(iniFile)
fmt.Println(doc.Get("topLevelValue")) // -> "foo"

section1 := doc.Section("section1")
fmt.Println(section1.Get("k1")) // -> "v1"

k2, err := section1.GetInt("k2")
fmt.Println(k2) // -> 1

doc.Set("k", "hello")
section1.SetInt("k2", -2)

iniFile = doc.ToString()
```

## Unmarshal Struct

```go
type MyStruct struct {
	Foo  string
	Bar  bool
	Baz  uint8
	// override the key names through tags
	K    string  `ini:"k"`
	K1   int     `ini:"k1"`
	User User    `ini:"user"`
}

type User struct {
	Name string `ini:"name"`
	Age  int    `ini:"age"`
}

doc := `
Foo=Lorem Ipsum
Bar=true
Baz=2
k=val
k1=-5

[user]
name=tom
age=23
`

cfg := MyStruct{}

ini.Unmarshal(doc, &cfg)
fmt.Println("MyStruct:", cfg) // -> MyStruct: {Lorem Ipsum true 2 val -5 {tom 23}}
```

## Marshal Struct

```go
type User struct {
	Name string `ini:"name"`
	Age  int    `ini:"age"`
}

type MyStruct struct {
	Foo  string
	Bar  bool
	Baz  uint8
	// override the key names through tags
	K    string  `ini:"k"`
	K1   int     `ini:"k1"`
	User User    `ini:"user"`
}

cfg := MyStruct{
	Foo: "Lorem Ipsum",
	Bar: true,
	Baz: 2,
	K: "val",
	K1: -5,
	User: User{
		Name: "Tom",
		Age: 23,
	},
}

iniFile, err := ini.Marshal(&cfg)
fmt.Println(iniFile)
```

Output:

```
Foo=Lorem Ipsum
Bar=true
Baz=2
k=val
k1=-5

[user]
name=Tom
age=23
```

## Sections

When marshaling/unmarshaling sections can be either nested structs, struct pointers or maps of string keys.

```go
type MySectionStruct struct {
	K string
}

type MyIni struct {
	TopLevelValue string

	Section1 MySectionStruct

	Section2 *MySectionStruct

	Section3 map[string]string
}
```

When using maps for sections, it is required that the map key type is `string`. If the key type is different it will be ignored. Also any values in the map that have a non-primitive type will be ignored as well (for example given a map like this: `map[string]any{"foo": 1, "bar": "hello", "baz": time.Now()}` - only `foo` and `bar` will be marshaled into the ini doc, since `time.Now()` is not of a primitive type.)

### Subsections

Subsections are a way to nest sections by using dots as path delimiters. For example `foo.bar` is a subsection of a section `foo`.

```go
doc := ini.NewDoc()

foo := doc.Section("foo")
fooBar := foo.Section("bar")

fooBar.Set("key", "value")

fmt.Println(doc.ToString())
```

Output:
```
[foo.bar]
key=value
```

Subsections can also be accessed by specifying the whole path as the argument for `Subsection()` method (e.x. `doc.Section("foo.bar")`)

When marshaling and un-marshaling nested structs will also create or read subsections. Maps cannot have subsections.

## Custom Marshal/Unmarshal

Custom marshaling an un-marshaling can be achieved by implementing these interfaces:

```go
type Marshalable interface {
	MarshalINI() (DocOrSection, error)
}

type Unmarshalable interface {
	UnmarshalINI(DocOrSection) error
}
```

### Example

```go
type MySection struct {
	myPrivateValue string
}

func (m *MySection) UnmarshalINI(doc DocOrSection) error {
	m.myPrivateValue = doc.Get("Value")
	return nil
}

func (m *MySection) MarshalINI() (DocOrSection, error) {
	section := ini.NewSection()
	section.Set("Value", m.myPrivateValue)
	return doc, nil
}

type MyIniFile struct {
	MySection *MySection
}

doc := "[MySection]\nValue=hello world\n"

myinifile := MyIniFile{}
err := ini.Unmarshal(doc, &myinifile)

fmt.Println(myinifile.MySection.myPrivateValue) // -> "hello world"
myinifile.MySection.myPrivateValue = "bye bye"

str, err := ini.Marshal(&myinifile)

fmt.Println(str) // -> "[MySection]\nValue=bye bye\n"
```

## Parse File

ini file

```ini
; this is comment
; author levene
; date 2021-8-1


a='23'34?::'<>,.'
c=d

[s1]
k=67676
k1 =fdasf
k2= sdafj3490&@)34 34w2

# comment
# 12.0.0.1
[s2]

k=3


k2=945
k3=-435
k4=0.0.0.0

k5=127.0.0.1
k6=levene@github.com

k7=~/.path.txt
k8=./34/34/uh.txt

k9=234@!@#$%^&*()324
k10='23'34?::'<>,.'

```

```go
file := "./test.ini"
doc, err := ini.Load(file)
k2 := doc.Section("s2").Get("k2")

fmt.Println(k2) // -> 945
```

## Write File

```go
filename := "./save.ini"
doc := ini.NewDoc()

doc.Set("a1", 1)

doc.Section("FooBar").Set("b", "hello")

err := doc.Save(filename)
```

save.ini
```ini
a1=1

[FooBar]
b=hello
```
