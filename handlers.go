package main

import (
	"log"
	"net/http"
    "gowiki/vendor/_nuts/github.com/gorilla/mux"
)

// handler to display wiki page contents, redirecting to edit template if page doesnt exist
func viewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("in viewHander...")

    title := mux.Vars(r)["title"]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/wiki/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// handler to edit the contents of a wiki page
func editHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("in editHandler...")

    title := mux.Vars(r)["title"]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// handler to save a new page or update an existing page
func saveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("in saveHandler...")

    title := mux.Vars(r)["title"]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}

	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/wiki/view/"+title, http.StatusFound)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("My 404"))
}