package main

import (
    "net/http"
    "text/template"
    "regexp"
    "gowiki/vendor/_nuts/gopkg.in/mgo.v2"
    "gowiki/vendor/_nuts/gopkg.in/mgo.v2/bson"
    "log"
)

/*
Web server built from this guide
https://golang.org/doc/articles/wiki/
*/

type Page struct {
	Title string
	Body []byte
}

var (
    templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
    validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
    globalSession *mgo.Session
    dbName string = "gowiki"
)

func (p *Page) save() error {

    session := globalSession.Copy()
    defer session.Close()

    // store saved page in mongo
    c := session.DB(dbName).C("pages")

    page := Page{}
    err := c.Find(bson.M{"title": p.Title}).One(&page)
    if err != nil {
        return c.Insert(p)
    } else {
        return c.Update(bson.M{"title": p.Title}, p)
    }
}

func loadPage(title string) (*Page, error) {
	//filename := title + ".txt"
	//body, err := ioutil.ReadFile("data/" + filename)
    session := globalSession.Copy()
    defer session.Close()

    c := session.DB(dbName).C("pages")

    page := Page{}
    err := c.Find(bson.M{"title": title}).One(&page)
	if err != nil {
        return nil, err
    }

    return &page, nil
}

func renderTemplate(w http.ResponseWriter, teml string, p *Page) {
    err := templates.ExecuteTemplate(w, teml + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    log.Println("in viewHander...")

    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+ title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    log.Println("in editHandler...")

    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

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
