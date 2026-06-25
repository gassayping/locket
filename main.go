package main

import (
	"encoding/xml"
	"io/fs"
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
	PostsFile string
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
	f, err := os.ReadFile(posts)
	if err == fs.ErrNotExist {
		os.Create(posts)
		f = []byte{}
	} else if err != nil {
		log.Fatal(err)
	}
	l.Posts = string(f)
	l.PostsFile = posts
	return l, nil
}

func (l *Locket) RequestAddPost(w http.ResponseWriter, r *http.Request) {
	l.Posts += "<item>\n" +
		r.FormValue("newItem") +
		"\n</item>"
	os.WriteFile(l.PostsFile, []byte(l.Posts), 0x664)
}

func (l *Locket) RequestFeed(w http.ResponseWriter, r *http.Request) {
	err := l.Template.Execute(w, l.Posts)
	if err != nil {
		log.Fatal(err)
	}
}
