package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var DEBUG bool = true

var ActEdit string = "edit"
var ActView string = "view"
var ActSave string = "save"
var ActError string = "error"

var ExtHTML string = ".html"
var ExtPage string = ".txt"

var HTMLPath string = "html/"
var PagesPath string = "pages/"

var EditHTML string = ActEdit + ExtHTML
var ViewHTML string = ActView + ExtHTML
var ErrorHTML string = ActError + ExtHTML

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := PagesPath + p.Title + ExtPage
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := PagesPath + title + ExtPage
	if DEBUG {
		fmt.Println("[loadPage]: ", filename)
	}
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет! Мне нравится %s!", r.URL.Path[1:])
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ExtHTML)
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	actView := "/" + ActView + "/"
	title := r.URL.Path[len(actView):]
	if DEBUG {
		fmt.Println(actView, title)
	}
	p, err := loadPage(title)
	templ := HTMLPath + ActView
	if err != nil {
		p = &Page{Title: title}
		templ = HTMLPath + ActError
	}
	renderTemplate(w, templ, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	actEdit := "/" + ActEdit + "/"
	title := r.URL.Path[len(actEdit):]
	if DEBUG {
		fmt.Println(actEdit, title)
	}
	p, err := loadPage(title)
	templ := HTMLPath + ActEdit
	if err != nil {
		p = &Page{Title: title, Body: []byte("начните редактировать...")}
	}
	renderTemplate(w, templ, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	actSave := "/" + ActSave + "/"
	actView := "/" + ActView + "/"
	title := r.URL.Path[len(actSave):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, actView+title, http.StatusFound)
}

func main() {
	p1 := &Page{Title: "Test Page", Body: []byte("This is a test page - тестовая страница")}
	p1.save()
	p2 , _ := loadPage("Test Page")
	fmt.Println(p2.Title, string(p2.Body))
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
