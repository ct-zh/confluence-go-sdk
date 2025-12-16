package confluence

import (
	"bytes"
	"html/template"
	"sync"

	tpl "confluence-go-sdk/confluence/template"
)

// TemplateRenderer 提供基于 html/template 的内容渲染功能
// 用于将结构化数据转换为 Confluence Storage Format (XHTML)
type TemplateRenderer struct {
	tmpl   *template.Template
	loader tpl.Loader
	mu     sync.RWMutex                  // 保护 cache
	cache  map[string]*template.Template // 缓存已解析的模板
}

// NewTemplateRenderer 创建一个新的渲染器
// 如果 loader 为 nil，则默认使用内置的 BuiltinLoader
func NewTemplateRenderer(loader tpl.Loader) *TemplateRenderer {
	if loader == nil {
		loader = tpl.NewBuiltinLoader()
	}
	return &TemplateRenderer{
		tmpl: template.New("confluence").Funcs(template.FuncMap{
			// 这里可以添加自定义 helper 函数，例如生成宏
			"codeBlock": renderCodeBlock,
			"infoPanel": renderInfoPanel,
		}),
		loader: loader,
		cache:  make(map[string]*template.Template),
	}
}

// Parse 解析模板字符串
// 注意：Parse 会直接修改基础模板，非并发安全，建议初始化时完成
func (r *TemplateRenderer) Parse(name, text string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.tmpl.New(name).Parse(text)
	return err
}

// LoadAndRender 加载并渲染模板
// 线程安全，支持缓存
func (r *TemplateRenderer) LoadAndRender(templateName string, data interface{}) (string, error) {
	// 1. 尝试从缓存获取
	r.mu.RLock()
	t, ok := r.cache[templateName]
	r.mu.RUnlock()

	if !ok {
		// 2. 缓存未命中，加载并解析
		content, err := r.loader.Load(templateName)
		if err != nil {
			return "", err
		}

		// 克隆基础模板以继承 FuncMap
		// 注意：template.Clone() 并非完全线程安全，需要在初始化后不再修改基础模板
		base, err := r.tmpl.Clone()
		if err != nil {
			return "", err
		}

		t, err = base.New(templateName).Parse(content)
		if err != nil {
			return "", err
		}

		// 3. 写入缓存
		r.mu.Lock()
		r.cache[templateName] = t
		r.mu.Unlock()
	}

	// 4. 渲染
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Render 渲染指定已解析的模板 (保留旧接口兼容性)
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
