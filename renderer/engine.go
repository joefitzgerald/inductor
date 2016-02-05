package renderer

import (
	"io"
	"text/template"
)

// Render generates the packer.json and Autounattend.xml files used by Packer
func (pt *PackerTemplate) Render(ro *RenderOptions, packerJSON io.Writer, autounattendXML io.Writer, vagrantfile io.Writer) error {
	err := pt.renderPackerJSON(ro, packerJSON)
	if err == nil {
		err = pt.renderAutounattendXML(ro, autounattendXML)
		if err == nil {
			err = pt.renderVagrantfile(ro, vagrantfile)
		}
	}
	return err
}

func (pt *PackerTemplate) renderPackerJSON(ro *RenderOptions, packerJSON io.Writer) error {
	tmpl, err := template.New("packer.json").Funcs(templateFuncs).Parse(pt.PackerTpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(packerJSON, ro)
}

func (pt *PackerTemplate) renderAutounattendXML(ro *RenderOptions, autounattendXML io.Writer) error {
	tmpl, err := template.New("Autounattend.xml").Funcs(templateFuncs).Parse(pt.AutounattendTpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(autounattendXML, ro)
}

func (pt *PackerTemplate) renderVagrantfile(ro *RenderOptions, vagrantfile io.Writer) error {
	tmpl, err := template.New("Vagrantfile").Funcs(templateFuncs).Parse(pt.VagrantfileTpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(vagrantfile, ro)
}
