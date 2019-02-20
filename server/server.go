package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"log"
	"regexp"
	"errors"
	_ "fmt"
    "webserver/mongo"
    "os"
)
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

type Page struct {
    Title string
    Body  []byte //https://blog.golang.org/go-slices-usage-and-internals
}

func (p *Page) save() error {
    filename := p.Title + ".jpg"
    var percorso = "sensori/" 
    //The octal integer literal 0600, passed as the third parameter to WriteFile, indicates that the file should be created with read-write permissions for the current user only
    os.MkdirAll(percorso+p.Title, os.FileMode(0522))
    return ioutil.WriteFile(percorso+p.Title+"/"+filename, p.Body, 0600)
}

func loadPage(title string) (*Page,error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
    	return nil, err
    }
    return &Page{Title: title, Body: body},nil
}

var Client = mongo.ConnectToMongo()

func main() {

    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}


func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
      renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
     err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    r.ParseForm()

    println(r.Form.Get("Device Id"))
    mongo.PostTemperature(r.Form.Get("Device Id"), r.Form.Get("timestamp"),r.Form.Get("temperatura") , Client)

    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
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