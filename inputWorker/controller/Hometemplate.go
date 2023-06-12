package controller

import (
	"html/template"
	"net/http"
)

func HomeTemplateHandler(w http.ResponseWriter, r *http.Request) {
	// Load and parse the template file
	tmpl, err := template.ParseFiles("./template/Home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := struct {
		Title string
	}{
		Title: "My Template",
	}

	// Render the template with the provided data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
