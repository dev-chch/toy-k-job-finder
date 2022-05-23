package main

import (
	"errors"
	"html/template"
	"io"
	"main/package/handler"

	"github.com/labstack/echo"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, name, data)
}

const Port string = "8000"

func main() {
	e := echo.New()

	templates := make(map[string]*template.Template)
	// templates["파일명.html"] = template.Must(template.ParseFiles("./view/템플릿 파일1.html", "./view/template/템플릿 파일2.html", "./view/template/템플릿 파일3.html" ...))
	templates["index.html"] = template.Must(template.ParseFiles("./view/index.html", "./view/template/header.html", "./view/template/footer.html"))
	templates["search.html"] = template.Must(template.ParseFiles("./view/search.html", "./view/template/header.html", "./view/template/footer.html"))
	templates["about.html"] = template.Must(template.ParseFiles("./view/about.html", "./view/template/header.html", "./view/template/footer.html"))
	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.Static("/assets", "assets")
	e.Static("/images", "images")

	e.GET("/", handler.Index)
	e.GET("/search", handler.Search)
	e.GET("/about", handler.About)

	e.Logger.Fatal(e.Start(":" + Port))
}
