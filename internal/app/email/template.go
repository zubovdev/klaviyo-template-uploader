package email

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// Template ...
type Template struct {
	images     []string
	data       TemplateData
	footerPath string
}

// NewEmailTemplate ...
func NewEmailTemplate(footerPath string, images []string) *Template {
	return &Template{footerPath: footerPath, images: images}
}

// Render ...
func (e *Template) Render() string {
	e.renderFooter()
	e.renderImages()

	wr := &bytes.Buffer{}
	t, _ := template.New("name").Parse(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title></title>
</head>
<body>
<table border="0" cellpadding="0" cellspacing="0" style="min-width:100%" width="100%">
    {{ .Images }}
</table>
{{ .Footer }}
</body>
</html>
`)
	_ = t.Execute(wr, e.data)
	return html.UnescapeString(wr.String())
}

// renderFooter ...
func (e *Template) renderFooter() {
	b, _ := os.ReadFile(filepath.Join(e.footerPath, "footer.html"))
	e.data.Footer = string(b)
}

// renderImages ...
func (e *Template) renderImages() {
	var images []string

	for range e.images {
		val := `
<tr>
    <td align="center">
        <img src="https://d3k81ch9hvuctc.cloudfront.net/company/S763nY/images/63edb50f-384d-48c6-8a8a-7a99131d74da.png" alt="img" width="576" style="max-width: 700px">
    </td>
</tr>`
		images = append(images, fmt.Sprintf("%s\n", val))
	}

	e.data.Images = strings.Join(images, "\n")
}
