package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

var templates map[string]*template.Template

func init() {
	if templates != nil {
		panic("Templates map already initialized")
	}
	templates = make(map[string]*template.Template)
	LoadTemplates()
}

func LoadTemplates() {
	log.Info("Loading Templates")
	if templates == nil {
		panic("Templates map cannot be nil.")
	}
	templatesDir := "../dist/templates/"
	layouts, err := filepath.Glob(templatesDir + "layouts/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	basefiles, err := filepath.Glob(templatesDir + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	for _, layout := range layouts {
		files := append(basefiles, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s doesn't exist", name)
	}
	w.WriteHeader(http.StatusOK)

	return tmpl.ExecuteTemplate(w, "base.tmpl", data)
}
