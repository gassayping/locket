package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Post struct {
	XMLName   xml.Name `xml:"item"`
	Title     string   `xml:"title"`
	Link      string   `xml:"link"`
	Published string   `xml:"pubDate"`
}

type Locket struct {
	Template  *template.Template
	Posts     string
	PostsFile *os.File
}

func main() {
	feed, err := FeedFromFile("./feed.tmpl", "./posts.xml")
	if err != nil {
		log.Fatal("Failed to read save file: ", err)
	}
	http.HandleFunc("POST /add", feed.RequestAddPost)
	http.HandleFunc("GET /feed.xml", feed.RequestFeed)
	http.HandleFunc("GET /", feed.RequestHomePage)
	log.Fatal(http.ListenAndServe(":7777", nil))
}

func FeedFromFile(tmpl, posts string) (Locket, error) {
	var l Locket
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return l, err
	}
	l.Template = t
	f, err := os.OpenFile(posts, os.O_RDWR|os.O_CREATE, os.FileMode(0o644))
	if err != nil {
		log.Fatalf("Could not open %v: %v", posts, err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	l.Posts = string(b)
	l.PostsFile = f
	return l, nil
}

func (l Locket) RequestHomePage(w http.ResponseWriter, r *http.Request) {
	f, err := os.ReadFile("./index.html")
	if err != nil {
		w.Write([]byte("Error reading file: " + err.Error()))
		return
	}
	w.Write(f)
}

func (l *Locket) RequestAddPost(w http.ResponseWriter, r *http.Request) {
	l.Posts += r.FormValue("newItem")
	if _, err := l.PostsFile.Write([]byte(l.Posts)); err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (l *Locket) RequestFeed(w http.ResponseWriter, r *http.Request) {
	err := l.Template.Execute(w, l.Posts)
	if err != nil {
		log.Fatal(err)
	}
}
