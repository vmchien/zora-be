package format

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"sync"
)

// Renderer caches compiled html/templates and executes them efficiently.
// It is safe for concurrent use across goroutines.
type Renderer struct {
	cache   sync.Map // map[string]*template.Template
	funcMap template.FuncMap
	bufPool sync.Pool // reuse buffers to reduce allocations
}

// NewRenderer returns a Renderer with a default FuncMap.
func NewRenderer() *Renderer {
	return &Renderer{
		funcMap: defaultFuncMap(),
		bufPool: sync.Pool{New: func() any { return new(bytes.Buffer) }},
	}
}

// Render renders a template string with given data and caches the compiled template by key.
// Example key: "vi:order.created:title"
func (r *Renderer) Render(key, tmplStr string, data any) (string, error) {
	if tmplStr == "" {
		return "", nil
	}

	// Try cache
	t, ok := r.cache.Load(key)
	if !ok {
		parsed, err := template.New(key).
			Option("missingkey=zero").
			Funcs(r.funcMap).
			Parse(tmplStr)
		if err != nil {
			return "", fmt.Errorf("parse template %s: %w", key, err)
		}
		r.cache.Store(key, parsed)
		t = parsed
	}

	buf := r.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer r.bufPool.Put(buf)

	if err := t.(*template.Template).Execute(buf, data); err != nil {
		return "", fmt.Errorf("execute template %s: %w", key, err)
	}

	return buf.String(), nil
}

// RenderPair renders both title and description templates with shared data.
func (r *Renderer) RenderPair(locale, typeKey, titleTmpl, descTmpl string, data any) (title, desc string, err error) {
	titleKey := fmt.Sprintf("%s:%s:title", locale, typeKey)
	descKey := fmt.Sprintf("%s:%s:desc", locale, typeKey)

	title, err = r.Render(titleKey, titleTmpl, data)
	if err != nil {
		return "", "", err
	}
	desc, err = r.Render(descKey, descTmpl, data)
	return title, desc, err
}

// ClearCache removes all cached templates (useful in tests or template reloads).
func (r *Renderer) ClearCache() {
	r.cache = sync.Map{}
}

// PathGet retrieves a nested value from a map[string]any using dot notation.
// Example:
//
//	payload := map[string]any{
//	  "user": map[string]any{"name": "Duong"},
//	  "client": map[string]any{"ip": "127.0.0.1"},
//	}
//	PathGet(payload, "client.ip") -> "127.0.0.1"
//	PathGet(payload, "user.name") -> "Duong"
//	PathGet(payload, "user.age")  -> nil
func PathGet(payload map[string]any, path string) any {
	if payload == nil || path == "" {
		return nil
	}
	cur := any(payload)
	for _, seg := range strings.Split(path, ".") {
		switch node := cur.(type) {
		case map[string]any:
			cur = node[seg]
		default:
			return nil
		}
		if cur == nil {
			return nil
		}
	}
	return cur
}
