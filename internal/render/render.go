package render

import (
	"fmt"
	"io"
	"io/fs"
	"maps"
	"text/template"
)

func New(fsys fs.FS, patterns ...string) (*FS, error) {
	var templ, errParse = template.New("").
		Funcs(Funcs).
		ParseFS(fsys, patterns...)
	if errParse != nil {
		return nil, fmt.Errorf("parsing templates: %w", errParse)
	}

	return &FS{
		back:     fsys,
		patterns: patterns,
		templ:    templ,
		values:   map[string]any{},
		funcs:    maps.Clone(Funcs),
	}, nil
}

type FS struct {
	back     fs.FS
	patterns []string
	templ    *template.Template
	values   map[string]any
	funcs    template.FuncMap
}

func (fsys *FS) SetValue(key string, value any) {
	fsys.values[key] = value
}

func (fsys *FS) Open(name string) (fs.File, error) {
	var stat, errStat = fs.Stat(fsys.back, name)
	if errStat != nil {
		return nil, errStat
	}

	var templ = fsys.templ.Lookup(name)
	if templ == nil {
		return fsys.back.Open(name)
	}

	var out, in = io.Pipe()

	go func() {
		var err = templ.Execute(in, fsys.values)
		in.CloseWithError(err)
	}()

	return templFile{
		stat:       stat,
		ReadCloser: out,
	}, nil
}

type templFile struct {
	stat fs.FileInfo
	io.ReadCloser
}

func (file templFile) Stat() (fs.FileInfo, error) {
	return file.stat, nil
}
