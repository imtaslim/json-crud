package handler

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Handler struct {
	templates *template.Template
	assets fs.FS
	decoder *schema.Decoder

}

const (
	createStudentPath = "/student/create"
	listStudentPath = "/students"
	updateStudentPath = "/student/{id}/edit"
	deleteStudentPath = "/student/{id}/delete"
)

func New(assets fs.FS, decoder *schema.Decoder) *mux.Router {
	h := &Handler{
		assets: assets,
		decoder: decoder,
	}

	r := mux.NewRouter()
	r.HandleFunc(createStudentPath, h.createStudent).Methods(http.MethodGet)
	r.HandleFunc(createStudentPath, h.storeStudent).Methods(http.MethodPost)
	r.HandleFunc(listStudentPath, h.listStudent).Methods(http.MethodGet)
	r.HandleFunc(updateStudentPath, h.editStudent).Methods(http.MethodGet)
	r.HandleFunc(updateStudentPath, h.updateStudent).Methods(http.MethodPost)
	r.HandleFunc(deleteStudentPath, h.deleteStudent).Methods(http.MethodGet)

	return r
}

func (s Student) Validate() error {
	vre := validation.Required.Error
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, vre("The name  is required")),
		validation.Field(&s.Class, vre("The class id is required")),
		validation.Field(&s.Roll, vre("The roll no id is required")),
	)
}

func (s *Handler) parseTemplates() error {
	templates := template.New("cms-templates")

	tmpl, err := templates.ParseFS(s.assets, "template/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}

func (s *Handler) lookupTemplate(name string) *template.Template {
	if err := s.parseTemplates(); err != nil {
		return nil
	}
	return s.templates.Lookup(name)
}

func (s *Handler) URLs() map[string]string {
	return map[string]string{
		"create": createStudentPath,
		"list": listStudentPath,
		"update": updateStudentPath,
		"delete": deleteStudentPath,
	}
}
