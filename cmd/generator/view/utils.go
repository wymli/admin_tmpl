package main

import (
	"bytes"
	"strings"
	"text/template"
)

type TmplRender struct {
	t    *template.Template
	opts TmplRenderOpts
}

type TmplRenderOpts struct {
	Tmpl                string
	RemovePrefixNewline bool
}

func NewTmplRenderOrDie(opts TmplRenderOpts) *TmplRender {
	t, err := template.New("").Parse(opts.Tmpl)
	if err != nil {
		panic(err)
	}
	return &TmplRender{t: t, opts: opts}
}

func (r *TmplRender) Funcs(v template.FuncMap) *TmplRender {
	r.t = r.t.Funcs(v)
	return r
}

func (r *TmplRender) Render(v any) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := r.t.Execute(buf, v); err != nil {
		return "", err
	}

	content := buf.String()

	if r.opts.RemovePrefixNewline {
		content = strings.TrimLeft(content, "\n")
	}

	return content, nil
}
