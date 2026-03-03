package main

import (
	"Continuity/ui"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// PageData holds the commonly used data passed to the templates.
type PageData struct {
	CreatedAt       string
	CSRFToken       string
	Form            any
	IsAuthenticated bool
	Flash           string
	UserView        *UserView
}

// NewPageData returns an initialized PageData struct.
func (b *backend) NewPageData(c echo.Context) *PageData {
	return &PageData{
		CreatedAt: time.Now().Format("03 Jan 2006"),
		CSRFToken: c.Get("csrf").(string),
		// UserView: b.currentUserView(c),
		// IsAuthenticated: b.isAuthenticated(c),
	}
}

// TemplateRenderer caches templates in memory and implements Renderer interface.
type TemplateRenderer struct {
	cache map[string]*template.Template
	md    goldmark.Markdown
}

// NewCacheRenderer returns new initialized TemplateRenderer.
func NewCacheRenderer() (*TemplateRenderer, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("evergarden"),
			),
		),
	)
	cr := &TemplateRenderer{
		md: md, cache: make(map[string]*template.Template),
	}
	if err := cr.buildCache(); err != nil {
		return nil, fmt.Errorf("couldn't build template cache: %w", err)
	}
	return cr, nil
}

// HasError returns true if the given field has an error in the given errors map.
func HasError(field string, errs map[string]string) bool {
	_, ok := errs[field]
	return ok
}

// Render executes the template with the given pageName and data to the writer.
// If the template hasn't already been cached, it builds the cache and then
// executes the template.
func (tr *TemplateRenderer) Render(w io.Writer, pageName string, data any, _ echo.Context) error {
	ct, ok := tr.cache[pageName]
	if !ok {
		if err := tr.buildCache(); err != nil {
			return fmt.Errorf("couldn't build template cache: %w", err)
		}
	}
	return ct.ExecuteTemplate(w, pageName, data)
}

// Markdown converts a raw Markdown string into safe HTML.
func (tr *TemplateRenderer) Markdown(body string) (template.HTML, error) {
	buf := new(bytes.Buffer)

	err := tr.md.Convert([]byte(body), buf)
	if err != nil {
		return "", fmt.Errorf("markdown conversion failed: %w", err)
	}
	return template.HTML(buf.String()), nil
}

// CreateHighlightCSS generates the chroma CSS file for the chosen style.
func CreateHighlightCSS(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("couldn't create highlight.css: %w", err)
	}
	defer func() { _ = file.Close() }()

	formatter := chromahtml.New(chromahtml.WithClasses(true))
	return formatter.WriteCSS(file, styles.Get("evergarden"))
}

// buildCache parses every page template and caches them by their file name.
func (tr *TemplateRenderer) buildCache() error {
	errFunc := template.FuncMap{
		"hasError": HasError,
	}
	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return fmt.Errorf("couldn't read pages: %w", err)
	}
	for _, page := range pages {
		pageName := filepath.Base(page)
		patterns := []string{
			"html/base.gohtml", "html/partial/*.gohtml", page,
		}
		templ, err := template.New(pageName).
			Funcs(errFunc).
			ParseFS(ui.Files, patterns...)

		if err != nil {
			return fmt.Errorf("couldn't parse template: %w", err)
		}
		tr.cache[pageName] = templ
	}
	return nil
}
