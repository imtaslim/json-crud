package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

func (h *Handler) editStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	tmpl := h.lookupTemplate("form.html")
	if tmpl == nil {
		http.Error(w, "No Template Found", http.StatusNotFound)
	}

	file, _ := ioutil.ReadFile("./student.json")
	stds := []Student{}
	std := Student{}

	if err := json.Unmarshal([]byte(file), &stds); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	for _, v := range stds {
		if v.ID == id {
			std = Student{
				ID:    id,
				Name:  v.Name,
				Roll:  v.Roll,
				Class: v.Class,
			}
			break
		}
	}

	data := tempData{
		FormAction: fmt.Sprintf("/student/%d/edit", id),
		Student:    std,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (h *Handler) updateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

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
	stds := []Student{}
	stdUpdate := []Student{}

	if err := json.Unmarshal([]byte(file), &stds); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	for _, v := range stds {
		var std Student
		if v.ID == id {
			std = Student{
				ID:    v.ID,
				Name:  form.Name,
				Roll:  form.Roll,
				Class: form.Class,
			}
		} else {
			std = Student{
				ID:    v.ID,
				Name:  v.Name,
				Roll:  v.Roll,
				Class: v.Class,
			}
		}
		stdUpdate = append(stdUpdate, std)
	}

	content, err := json.Marshal(stdUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err = ioutil.WriteFile("./student.json", content, 0644); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	http.Redirect(w, r, listStudentPath, http.StatusSeeOther)
}
