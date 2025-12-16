package confluence

import (
	"bytes"
	"html/template"
)

// TemplateRenderer 提供基于 html/template 的内容渲染功能
// 用于将结构化数据转换为 Confluence Storage Format (XHTML)
type TemplateRenderer struct {
	tmpl *template.Template
}

// NewTemplateRenderer 创建一个新的渲染器
func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		tmpl: template.New("confluence").Funcs(template.FuncMap{
			// 这里可以添加自定义 helper 函数，例如生成宏
			"codeBlock": renderCodeBlock,
			"infoPanel": renderInfoPanel,
		}),
	}
}

// Parse 解析模板字符串
func (r *TemplateRenderer) Parse(name, text string) error {
	_, err := r.tmpl.New(name).Parse(text)
	return err
}

// Render 渲染指定模板
func (r *TemplateRenderer) Render(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := r.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Helper Functions for Confluence Macros

// renderCodeBlock 生成代码块宏
func renderCodeBlock(language, content string) template.HTML {
	// Confluence Storage Format for Code Block
	html := `
<ac:structured-macro ac:name="code" ac:schema-version="1">
  <ac:parameter ac:name="language">` + language + `</ac:parameter>
  <ac:plain-text-body><![CDATA[` + content + `]]></ac:plain-text-body>
</ac:structured-macro>`
	return template.HTML(html)
}

// renderInfoPanel 生成信息提示面板宏
func renderInfoPanel(title, content string) template.HTML {
	html := `
<ac:structured-macro ac:name="info" ac:schema-version="1">
  <ac:parameter ac:name="title">` + title + `</ac:parameter>
  <ac:rich-text-body><p>` + content + `</p></ac:rich-text-body>
</ac:structured-macro>`
	return template.HTML(html)
}
