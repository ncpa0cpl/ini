

![logo](./logo.png)

#  INI Parser & Reader Writer Library

## Introduction

The INI Parser & Reader Writer Library is a fast and easy-to-use library for parsing and manipulating INI files in the Go programming language. It provides functionality to read INI files from both byte slices and files, supports real-time file monitoring, and offers options to unmarshal INI data into Go structs, marshal data to JSON, and write data back to files.


[![Build Status](https://travis-ci.org/meolu/walden.svg?branch=master)](https://github.com/ncpa0cpl/ini)
![version](https://img.shields.io/badge/version-0.1.5-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/ncpa0cpl/ini)](https://goreportcard.com/report/github.com/ncpa0cpl/ini)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)

## Features
* **Read by []byte**: The library allows you to parse INI data stored in a byte slice.
* **Read by file**: You can also read INI data directly from a file.
* **Real-time file monitoring**: The library supports file monitoring, allowing you to observe changes in INI files without the need for manual reloading.
* **Unmarshal to Struct**: It provides the ability to map INI data to Go structs, making it convenient to work with structured data.
* **Marshal to JSON**: You can easily convert INI data to JSON format using the library's marshal functionality.
* **Write to File**: The library allows you to write INI data back to files, preserving the original file format.


## Installation

```shell
go get github.com/ncpa0cpl/ini
```


## Example

```go
import (
    "fmt"
    "github.com/ncpa0cpl/ini"
)
```

### GetValue

```go
doc := `
[section]
k=v

[section1]
k1=v1
k2=1
k3=3.5
k4=0.0.0.0
`
v1 := ini.New().Load([]byte(doc)).Section("section1").Get("k1")
fmt.Println(v1)
```

Output

```
v1
```


```go
i := ini.New().Load([]byte(doc))
v1 := i.Section("section1").Get("k1")
v2 := i.GetInt("k2")
v3 := i.GetFloat64("k3")
v4 := i.Get("k4")
v5 := i.GetIntDef("keyint", 10)
v6 := i.GetDef("keys", "defualt")

fmt.Printf("v1:%v v2:%v v3:%v v4:%v v5:%v v6:%v\n", v1, v2, v3, v4, v5, v6)
```

Output

```
v1:v1 v2:1 v3:3.5 v4:0.0.0.0 v5:10 v6:defualt
```


### Marshal2Json

```go
fmt.Println(string(i.Marshal2Json()))
```

Output

```json
{"section":{"k":"v"},"section1":{"k1":"v1","k2":"1","k3":"3.5","k4":"0.0.0.0"}}
```

### Unmarshal Struct

```go
type TestConfig struct {
	K    string  `ini:"k" json:"k,omitempty"`
	K1   int     `ini:"k1" json:"k1,omitempty"`
	K2   float64 `ini:"k2"`
	K3   int64   `ini:"k3"`
	User User    `ini:"user"`
}

type User struct {
	Name string `ini:"name"`
	Age  int    `ini:"age"`
}

```

```go
doc := `
k=v
k1=2
k2=2.2
k3=3

[user]
name=tom
age=-23
`

cfg := TestConfig{}

ini.Unmarshal([]byte(doc), &cfg)
fmt.Println("cfg:", cfg)
```

Output

```
cfg: {v 2 2.2 3 {tom -23}}
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
ini.New().LoadFile(file).Section("s2").Get("k2")

fmt.Println(string(ini.Marshal2Json()))
```

Output

```
945
```



### Watch File

```go
file := "./test.ini"

idoc := ini.New().WatchFile(file)
v := idoc.Section("s2").Get("k1")
fmt.Println("v:", v1)

// modify k1=v1   ==> k1=v2
time.Sleep(10 * time.Second)

v = idoc.Section("s2").Get("k1")
fmt.Println("v:", v1)
```

Output

```
v: v1
v: v2
```



Print file with json

```go
file := "./test.ini"
fmt.Println(string(ini.New().LoadFile(file).Marshal2Json()))
```

Output

```json
{
  "a": "'23'34?::'<>,.'",
  "c": "d",
  "s1": {
    "k": "67676",
    "k2": "34w2"
  },
  "s2": {
    "k": "3",
    "k10": "'23'34?::'<>,.'",
    "k2": "945",
    "k3": "-435",
    "k4": "0.0.0.0",
    "k5": "127.0.0.1",
    "k6": "levene@github.com",
    "k7": "~/.path.txt",
    "k8": "./34/34/uh.txt",
    "k9": "234@!@#$%^&*()324"
  }
}
```

### Set Ini 
```go
doc := `
k =v
[section]
a=b
c=d
`
ini := New().Load([]byte(doc)).Section("section")
fmt.Println("--------------------------------")
ini.Dump()

fmt.Println("--------------------------------")
ini.Set("a", 11).Set("c", 12.3).Section("").Set("k", "SET")
ini.Dump()

v := ini.Section("section").GetInt("a")

if v != 11 {
    t.Errorf("Error: %d", v)
}

v1 := ini.GetFloat64("c")

if v1 != 12.3 {
    t.Errorf("Error: %f", v1)
}

v2 := ini.Section("").Get("k")
if v2 != "SET" {
    t.Errorf("Error: %s", v2)
}
```



### Wirte Ini

```go
filename := "./save.ini"
ini := New().Set("a1", 1)
ini.Save(filename)
fmt.Println(ini.Err())

ini2 := New().Set("a1", 1).Section("s1").Set("a2", "v2")
ini2.Save(filename)
fmt.Println(ini2.Err())

// ------

doc := `
; 123
c11=d12312312
# 434

[section]
k=v
; dsfads 
;123
#3452345


[section1]
k1=v1

[section3]
k3=v3
`
ini3 := New().Load([]byte(doc))
ini.Save("./save.ini")

```

file content
```ini

; 123
c11 = d12312312

# 434
[section]
k = v

; dsfads 
;123
#3452345
[section1]
k1 = v1

[section3]
k3 = v3

```



### Dump AST struct

```
INIDocNode {
    CommentNode {
        Comment: ; this is comment
        Line: 0
    }
    CommentNode {
        Comment: ; author levene
        Line: 1
    }
    CommentNode {
        Comment: ; date 2021-8-1
        Line: 2
    }
    KVNode {
        Key: a
        Value: '23'34?::'<>,.'
        Line: 5
    }
    KVNode {
        Key: c
        Value: d
        Line: 6
    }
    Section {
        Section: [s1]
        Line: 8
        KVNode {
            Value: 67676
            Line: 9
            Key: k
        }
        KVNode {
            Key: k1
            Value: fdasf
            Line: 10
        }
        KVNode {
            Value: 4w2
            Line: 11
            Key: k2
        }
    }
    CommentNode {
        Comment: # comment
        Line: 13
    }
    CommentNode {
        Line: 14
        Comment: # 12.0.0.1
    }
    Section {
        Section: [s2]
        Line: 15
        KVNode {
            Value: 3
            Line: 17
            Key: k
        }
        KVNode {
            Value: 945
            Line: 20
            Key: k2
        }
        KVNode {
            Key: k3
            Value: -435
            Line: 21
        }
        KVNode {
            Line: 22
            Key: k4
            Value: 0.0.0.0
        }
        KVNode {
            Line: 24
            Key: k5
            Value: 127.0.0.1
        }
        KVNode {
            Key: k6
            Value: levene@github.com
            Line: 25
        }
        KVNode {
            Key: k7
            Value: ~/.path.txt
            Line: 27
        }
        KVNode {
            Line: 28
            Key: k8
            Value: ./34/34/uh.txt
        }
        KVNode {
            Key: k9
            Value: 234@!@#$%^&*()324
            Line: 30
        }
        KVNode {
            Key: k10
            Value: '23'34?::'<>,.'
            Line: 31
        }
    }
}
```

## License

[MIT](https://github.com/RichardLitt/standard-readme/blob/master/LICENSE)  
