package ini_test

import (
	"io/ioutil"
	"testing"

	"github.com/ncpa0cpl/ini"
	"github.com/stretchr/testify/assert"
)

func TestIni0(t *testing.T) {
	doc := `
[section]
k =v
`
	i_doc := ini.New().Load([]byte(doc)).Section("section")

	v := i_doc.Get("k")
	if v != "v" {
		t.Errorf("error %s", v)
	}
}

func TestIni2(t *testing.T) {
	doc := `
[section]
k=v
k1 = 1
`
	ini := ini.New().Load([]byte(doc))

	v := ini.Section("section").Get("k")
	if v != "v" {
		t.Errorf("error %s", v)
	}

	iv := ini.GetInt("k1")
	if iv != 1 {
		t.Errorf("error %d", iv)
	}

	iv = ini.GetInt("k2")
	if iv != 0 {
		t.Errorf("error %d", iv)
	}

	iv = ini.GetIntDef("k2", 12)
	if iv != 12 {
		t.Errorf("error %d", iv)
	}

}

func TestIni3(t *testing.T) {
	doc := `
a =b
c=d
[section]
k =v
`
	ini := ini.New().Load([]byte(doc))

	v := ini.Section("section").Get("k")
	if v != "v" {
		t.Errorf("error %s", v)
	}
	v = ini.Section("").Get("a")
	if v != "b" {
		t.Errorf("error %s", v)
	}

}

func TestIni4(t *testing.T) {
	doc := `
a =b
c= d

a1 = 2.1
`
	ini := ini.New().Section("").Load([]byte(doc))
	v := ini.Get("a")
	if v != "b" {
		t.Errorf("error %s:%d", v, len(v))
	}

	iv := ini.GetInt("a")
	if iv != 0 {
		t.Errorf("error %d", iv)
	}

	iv = ini.GetIntDef("a", 10)
	if iv != 10 {
		t.Errorf("error %d", iv)
	}

	iv = ini.GetIntDef("a1", 10)
	if iv != 10 {
		t.Errorf("error %d", iv)
	}

	fv := ini.GetFloat64Def("a1", 10)
	if fv != 2.1 {
		t.Errorf("error %d", iv)
	}

}

func TestIni5(t *testing.T) {
	doc := `
a =b
[s1]
k=v
k1 = v12

[s2]
k2=v2
k2= v22

[s3]
k =v
a= b
`
	ini := ini.New().Load([]byte(doc))
	json_str := string(ini.Marshal2Json())
	if json_str != `{"a":"b","s1":{"k":"v","k1":"v12"},"s2":{"k2":"v22"},"s3":{"a":"b","k":"v"}}` {
		t.Errorf("error %v", json_str)
	}

}

func TestIniFile(t *testing.T) {
	file := "./test.ini"
	ini := ini.New().LoadFile(file)

	a := assert.New(t)
	a.Equal("'23'34?::'<>,.'", ini.Get("a"), "'a' is incorrect")
	a.Equal("d", ini.Get("c"), "'c' is incorrect")
	a.Equal(67676, ini.Section("s1").GetInt("k"), "'s1.k' is incorrect")
	a.Equal("fdasf", ini.Section("s1").Get("k1"), "'s1.k1' is incorrect")
}

func TestIniDelete(t *testing.T) {
	doc := `
k=v
a=b
c=d
[section]

`
	ini := ini.New().Load([]byte(doc))

	ini.Del("a")

	ini.Del("c")

	ini.Del("k")

}

func TestIniSet(t *testing.T) {
	doc := `
k =v
[section]
a=b
c=d
`
	ini := ini.New().Load([]byte(doc)).Section("section")

	ini.Set("a", 11).Set("c", 12.3).Section("").Set("k", "SET")

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

	ini.Set("a1", 1).Section("section").Set("k1", 11.11)

}

func TestIniSave(t *testing.T) {
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
	ini := ini.New().Load([]byte(doc))

	ini.Save("./save.ini")
}

func TestIniSave2(t *testing.T) {

	filename := "./save.ini"
	ini := ini.New().Set("a1", 1)
	ini.Save(filename)

	bts, _ := ioutil.ReadFile(filename)

	if string(bts) != "a1 = 1\n" {
		t.Errorf("Error: %v", string(bts))
	}
}

func TestIniSave3(t *testing.T) {

	filename := "./save.ini"
	ini := ini.New().Set("a1", 1).Section("s1").Set("a2", "v2")
	ini.Save(filename)

	bts, _ := ioutil.ReadFile(filename)

	if string(bts) != "a1 = 1\n\n[s1]\na2 = v2\n" {
		t.Errorf("Error: %v", string(bts))
	}
}
