package renderer

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/joefitzgerald/inductor/tpl"
)

type engine struct {
	renderOptions *RenderOptions
	outDir        string
}

// New creates a new Renderer instance
func New(opts *RenderOptions, outDir string) Renderer {
	return &engine{
		renderOptions: opts,
		outDir:        outDir,
	}
}

// Render generates the packer.json and Autounattend.xml files used by Packer
func (e *engine) Render(tc tpl.TemplateContainer) error {
	for _, t := range tc.ListTemplates() {
		if err := e.writeTemplate(t); err != nil {
			return err
		}
	}
	return nil
}

func (e *engine) writeTemplate(t tpl.Templater) error {
	f, err := os.Create(filepath.Join(e.outDir, t.BaseFilename()))
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	e.renderTemplate(t, w)
	return nil
}

func (e *engine) renderTemplate(tpl tpl.Templater, outWriter io.Writer) error {
	var buffer bytes.Buffer
	tpl.Content(&buffer)
	tmpl, err := template.New("tpl").Funcs(templateFuncs).Parse(buffer.String())
	if err != nil {
		return err
	}
	return tmpl.Execute(outWriter, e.renderOptions)
}
