package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (h *Handler) listStudent(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("./student.json")
	std := []Student{}
 
	if err := json.Unmarshal([]byte(file), &std); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	
	tmpl := h.lookupTemplate("list.html")
	if tmpl == nil {
		http.Error(w, "No Template Found", http.StatusNotFound)
	}

	data := tempData{
		URLs:       h.URLs(),
		Data:       std,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}