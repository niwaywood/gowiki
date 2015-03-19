package main

import (
	"gowiki/vendor/_nuts/gopkg.in/mgo.v2"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

// A struct to represent the contents of a wiki page
type Page struct {
	Title string `bson:"title,omitempty"`
	Body  []byte `bson:"body,omitempty"`
}

var (
	templates     = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
	validPath     = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
	globalSession *mgo.Session
	dbName        string = "gowiki"
)

// Saves the contents of the Page struct into mongoDB, inserting a new entry
// if its a new page and updating an existing entry if the page already exists
func (p *Page) save() error {
	session := globalSession.Copy()
	defer session.Close()

	// store saved page in mongo
	c := session.DB(dbName).C("pages")

	_, err := c.Upsert(Page{Title: p.Title}, p)
	if err != nil {
		return err
	}
	return nil
}

// Loads an existing page from mongo. Returns `not found` error if page doesn't exist
func loadPage(title string) (*Page, error) {
	session := globalSession.Copy()
	defer session.Close()

	c := session.DB(dbName).C("pages")

	page := Page{}
	err := c.Find(Page{Title: title}).One(&page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// renders template to be sent to web client
func renderTemplate(w http.ResponseWriter, teml string, p *Page) {
	err := templates.ExecuteTemplate(w, teml+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	log.Println("Starting server...")

	// connect to mongodb
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	globalSession = session
	defer session.Close()

    // setup handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
	})
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

    log.Println("Started server, listening on post 8080")
	http.ListenAndServe(":8080", nil)

}
