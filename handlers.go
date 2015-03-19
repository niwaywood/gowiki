package main

import (
	"log"
	"net/http"
)

// handler middleware which ensures that the path is valid
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we will extract the page title from the Request,
		// and call the provided handler 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// handler to display wiki page contents, redirecting to edit template if page doesnt exist
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Println("in viewHander...")

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// handler to edit the contents of a wiki page
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Println("in editHandler...")

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// handler to save a new page or update an existing page
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Println("in saveHandler...")

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}

	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
