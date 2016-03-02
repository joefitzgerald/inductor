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
		tOutPath := filepath.Join(e.outDir, t.BaseFilename())
		tOutFile, err := os.Create(tOutPath)
		if err != nil {
			return err
		}
		defer tOutFile.Close()
		tOutWriter := bufio.NewWriter(tOutFile)
		e.renderTemplate(t, tOutWriter)
	}
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
