package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Roll  string `json:"roll"`
	Class string `json:"class"`
}

type tempData struct {
	URLs       map[string]string
	FormError  map[string]string
	FormAction string
	Student
	Data []Student
}

func (h *Handler) createStudent(w http.ResponseWriter, r *http.Request) {
	tmpl := h.lookupTemplate("form.html")
	if tmpl == nil {
		http.Error(w, "No Template Found", http.StatusNotFound)
	}

	data := tempData{
		FormAction: createStudentPath,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (h *Handler) storeStudent(w http.ResponseWriter, r *http.Request) {
	var form Student

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err := h.decoder.Decode(&form, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err := form.Validate(); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}

		data := tempData{
			Student:    form,
			FormError:  vErrs,
			FormAction: createStudentPath,
		}
		tmpl := h.lookupTemplate("form.html")
		if tmpl == nil {
			http.Error(w, "No Template Found", http.StatusNotFound)
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	file, _ := ioutil.ReadFile("./student.json")

	data := []Student{}

	if err := json.Unmarshal([]byte(file), &data); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	if len(data) > 0 {
		form.ID = data[len(data)-1].ID + 1
	} else {
		form.ID = 1
	}

	data = append(data, form)

	content, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	err = ioutil.WriteFile("./student.json", content, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	http.Redirect(w, r, listStudentPath, http.StatusSeeOther)
}
