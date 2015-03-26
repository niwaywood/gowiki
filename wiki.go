package main

import (
	"gowiki/vendor/_nuts/gopkg.in/mgo.v2"
	"log"
	"net/http"
	"regexp"
	"text/template"
    "gowiki/vendor/_nuts/github.com/gorilla/mux"
    "gowiki/vendor/_nuts/github.com/codegangsta/negroni"
)

// A struct to represent the contents of a wiki page
type Page struct {
	Title string `bson:"title,omitempty"`
	Body  []byte `bson:"body,omitempty"`
}

var (
	templates     = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
	validPath     = regexp.MustCompile("^/wiki/(edit|save|view)/([a-zA-Z0-9]+)$")
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
    // connect to mongodb
    session, err := mgo.Dial("localhost")
    if err != nil {
        log.Println("Unable to connect to MongoDB")
        return
    }
    globalSession = session
    defer session.Close()

    // setup mux routers
    r := mux.NewRouter()
    mr := mux.NewRouter()

    // setup handlers
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/wiki/view/FrontPage", http.StatusFound)
    })
    // route specific middleware for /wiki paths
    r.PathPrefix("/wiki").Handler(negroni.New(
    negroni.HandlerFunc(validateURL),
    negroni.Wrap(mr),
    ))
    sub := mr.PathPrefix("/wiki").Subrouter()
    sub.HandleFunc("/view/{title}", viewHandler).Methods("GET")
    sub.HandleFunc("/edit/{title}", editHandler).Methods("GET")
    sub.HandleFunc("/save/{title}", saveHandler).Methods("POST")

    // setup negroni middleware for all routes
    n := negroni.New(negroni.HandlerFunc(myMiddleware), negroni.NewLogger())
    n.UseHandler(r)

    n.Run(":8080")
}
