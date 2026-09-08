package main

import (
	"io"
	"locket/rss"
	"log"
	"net/http"
	"os"
)

type Locket struct {
	Feed rss.RSS
	File *os.File
}

func main() {
	feed := FeedFromFile("./feed.xml")
	http.HandleFunc("POST /add", feed.RequestAddPost)
	http.HandleFunc("GET /feed.xml", feed.RequestFeed)
	http.HandleFunc("GET /", feed.RequestHomePage)
	log.Fatal(http.ListenAndServe(":7777", nil))
}

func FeedFromFile(filePath string) Locket {
	var l Locket
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.FileMode(0o644))
	if err != nil {
		log.Fatalf("Could not open %v: %v", filePath, err)
	}
	l.File = f
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	if len(b) == 0 {
		l.Feed = rss.NewFeed()
	} else {
		l.Feed, err = rss.FeedFromXML(b)
		if err != nil {
			log.Fatalf("Error parsing feed from file %s: %s", filePath, err)
		}
	}
	return l
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
	err := l.Feed.AddItem([]byte(r.FormValue("newItem")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed XML received: " + err.Error()))
		return
	}
	b, err := l.Feed.GetXML()
	if err != nil {
		log.Fatal(err)
	}
	l.File.Truncate(0)
	l.File.Seek(0, 0)
	if _, err := l.File.Write(b); err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (l *Locket) RequestFeed(w http.ResponseWriter, r *http.Request) {
	b, err := l.Feed.GetXML()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "application/xml")
	w.Write(b)
}
