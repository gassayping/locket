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
	PostsFile os.File
}

func main() {
	feed, err := FeedFromFile("./feed.tmpl", "./posts.xml")
	if err != nil {
		log.Fatal("Failed to read save file: ", err)
	}
	http.HandleFunc("POST /add", feed.RequestAddPost)
	http.HandleFunc("GET /feed.xml", feed.RequestFeed)
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
	return l, nil
}

func (l *Locket) RequestAddPost(w http.ResponseWriter, r *http.Request) {
	l.Posts += "<item>\n" +
		r.FormValue("newItem") +
		"\n</item>"
	l.PostsFile.Write([]byte(l.Posts))
}

func (l *Locket) RequestFeed(w http.ResponseWriter, r *http.Request) {
	err := l.Template.Execute(w, l.Posts)
	if err != nil {
		log.Fatal(err)
	}
}
