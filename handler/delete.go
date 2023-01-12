package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) deleteStudent(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	file, _ := ioutil.ReadFile("./student.json")
	stds := []Student{}
	stdUpdate := []Student{}

	if err := json.Unmarshal([]byte(file), &stds); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	for i := range stds {
        if stds[i].ID == id {
            stdUpdate = append(stds[:i], stds[i+1:]...)
        }
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