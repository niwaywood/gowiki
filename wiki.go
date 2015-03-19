package main

import (
	"gowiki/vendor/_nuts/gopkg.in/mgo.v2"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

/*
Web server built from this guide
https://golang.org/doc/articles/wiki/
*/

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

func renderTemplate(w http.ResponseWriter, teml string, p *Page) {
	err := templates.ExecuteTemplate(w, teml+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	log.Println("Started server")

	// connect to mongodb
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	globalSession = session
	defer session.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
	})
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.ListenAndServe(":8080", nil)

}
