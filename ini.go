package ini

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/wlevene/ini/ast"
	"github.com/wlevene/ini/lexer"
	"github.com/wlevene/ini/parser"
	"github.com/wlevene/ini/token"

	"github.com/fsnotify/fsnotify"
)

type (
	Ini struct {
		currectSection *ast.SetcionNode
		src            []byte // source
		l              *lexer.Lexer
		p              *parser.Parser
		doc            *ast.Doc
		err            error

		watcher       *fsnotify.Watcher
		exitWatchChan chan bool
	}
)

func New() *Ini {

	in := &Ini{}
	return in
}

func (in *Ini) Err() error {
	return in.err
}

func (in *Ini) LoadFile(file string) *Ini {

	// read file content
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		in.err = err
		return in
	}

	in.Load(bts)

	return in
}

func (in *Ini) WatchFile(file string) *Ini {

	in.LoadFile(file)
	in.watch(file)

	return in
}

func (in *Ini) watch(file string) {
	if file == "" {
		return
	}

	in.watcher, in.err = fsnotify.NewWatcher()
	in.exitWatchChan = make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-in.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					in.LoadFile(file)
				}
			case _, ok := <-in.watcher.Errors:
				if !ok {
					return
				}
			case <-in.exitWatchChan:
				return
			}

		}
	}()

	in.err = in.watcher.Add(file)
}

func (in *Ini) StopWatch() *Ini {

	in.watcher.Close()
	in.exitWatchChan <- true
	return in
}

func (in *Ini) Load(doc []byte) *Ini {

	if len(doc) <= 0 {
		return in
	}

	in.src = doc
	in.l = lexer.New(string(in.src))
	in.p = parser.New(in.l)

	in.doc, in.err = in.p.ParseDocument()
	return in
}

func (in *Ini) Dump() {

	if in.doc == nil {
		return
	}

	in.doc.Dump()
}

func (this *Ini) Marshal2Map() map[string]interface{} {

	if this.doc == nil {
		return nil
	}

	if this.err != nil {
		return nil
	}

	kvMaps := make(map[string]interface{})
	for _, kv := range this.doc.KVs {
		kvMaps[kv.Key.Literal] = kv.Value.Literal
	}

	for _, sec := range this.doc.SectS {

		secMap := make(map[string]interface{})

		for kv := sec.FirstChild(); kv != nil; kv = kv.NextSibling() {
			if kvnode, ok := kv.(*ast.KVNode); ok {
				secMap[kvnode.Key.Literal] = kvnode.Value.Literal
			}
		}

		kvMaps[sec.Name.Literal] = secMap
	}

	return kvMaps
}

func (this *Ini) Marshal2Json() []byte {

	kvMaps := this.Marshal2Map()

	if kvMaps == nil {
		return nil
	}

	result, err := json.Marshal(kvMaps)
	this.err = err

	return result
}

func (this *Ini) Section(section string) *Ini {

	if this.doc == nil {
		return this
	}

	if this.err != nil {
		return this
	}

	this.sectionForAstDoc(section)
	return this
}

func (this *Ini) Get(key string) string {
	return this.GetDef(key, "")
}

func (this *Ini) GetDef(key string, def string) string {

	if this.doc == nil ||
		this.err != nil {
		return def
	}

	if key == "" {
		return def
	}

	tok := this.getToken(key)
	if tok.Type != token.TokenType_VALUE {
		return def
	}

	return tok.Literal
}

func (this *Ini) GetInt(key string) int {

	return this.GetIntDef(key, 0)
}

func (this *Ini) GetIntDef(key string, def int) int {

	val := this.Get(key)
	if val == "" {
		return def
	}

	ival, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return ival
}

func (this *Ini) GetInt64(key string) int64 {

	return this.GetInt64Def(key, 0)
}

func (this *Ini) GetInt64Def(key string, def int64) int64 {

	val := this.Get(key)
	if val == "" {
		return def
	}

	ival, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return def
	}

	return ival
}

func (this *Ini) GetFloat64(key string) float64 {

	return this.GetFloat64Def(key, 0)
}

func (this *Ini) GetFloat64Def(key string, def float64) float64 {

	val := this.Get(key)
	if val == "" {
		return def
	}

	fval, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return def
	}
	return fval
}

// TODO: implement Write
func (this *Ini) Set(key string, val interface{}) error {
	return nil
}

func (this *Ini) Del(key string) error {

	if this.doc == nil ||
		this.err != nil {
		return nil
	}

	if key == "" {
		return nil
	}

	if this.currectSection == nil {

		for index, kvnodev2 := range this.doc.KVs {

			if kvnodev2.Key.Literal == key {
				this.doc.KVs = append(this.doc.KVs[:index], this.doc.KVs[(index+1):]...)
				break
			}
		}
	} else {
		for c := this.currectSection.FirstChild(); c != nil; c = c.NextSibling() {
			kvnodev2 := c.(*ast.KVNode)
			if kvnodev2.Key.Literal == key {

				this.currectSection.RemoveChild(this.currectSection, c)
				break
			}
		}
	}

	return nil
}

func (this *Ini) DelSection(section string) error {
	return nil
}

func (this *Ini) Save(filename string) error {
	return nil
}

func (this *Ini) SaveFile(filename string) error {
	return nil
}

// ----------------------------------------------------------------

func (this *Ini) sectionForAstDoc(section string) {

	if this.doc == nil ||
		this.err != nil {
		return
	}

	this.currectSection = nil
	for _, v := range this.doc.SectS {
		if v.Name.Literal == section {
			this.currectSection = v
			return
		}
	}
}

func (this *Ini) getToken(key string) token.Token {

	var tok token.Token

	if key == "" {
		return tok
	}

	if this.currectSection == nil {

		for _, kvnodev2 := range this.doc.KVs {

			if kvnodev2.Key.Literal == key {

				tok = kvnodev2.Value

				return tok
			}
		}

	} else {
		for c := this.currectSection.FirstChild(); c != nil; c = c.NextSibling() {

			kvnodev2 := c.(*ast.KVNode)
			if kvnodev2.Key.Literal == key {
				tok = kvnodev2.Value
				return tok
			}
		}
	}

	return tok
}
