
#  INI Parser & Reader Writer Library

## Introduction

The INI Parser & Reader Writer Library is a fast and easy-to-use library for parsing and manipulating INI files in the Go programming language. It provides functionality to read INI files from both strings and files, and offers options to marshal and unmarshal INI data into Go structs, and write data back to files.

## Features
* **Read by string**: The library allows you to parse INI data stored in a string.
* **Read by file**: You can also read INI data directly from a file.
* **Unmarshal to Struct**: It provides the ability to map INI data to Go structs, making it convenient to work with structured data.
* **Marshal from Struct**: You can easily convert Go struct into INI data using the library's marshal functionality.
* **Write to File**: The library allows you to write INI data back to files.


## Installation

```shell
go get github.com/ncpa0cpl/ini
```


## Example

### GetValue

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

fmt.Println(section1.Get("k1")) // -> "v"

k2, err := section1.GetInt("k2")
fmt.Println(k2) // -> 1
```

### Unmarshal Struct

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

### Marshal Struct

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

### Parse File

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

### Write Ini

```go
filename := "./save.ini"
doc := ini.NewDoc()

doc.Set("a1", 1)

doc.Section("FooBar").Set("b", "hello")

err := ini.Save(filename)
```

save.ini
```ini
a1=1

[FooBar]
b=hello
```
