package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/feeds"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Author struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type FrontMatter struct {
	Title string `json:"title"`
}

type Entry struct {
	EntryId     int         `json:"entryId"`
	FrontMatter FrontMatter `json:"frontMatter"`
	Created     Author      `json:"created"`
	Updated     Author      `json:"updated"`
}

type Page struct {
	TotalElements    int  `json:"totalElements"`
	NumberOfElements int  `json:"numberOfElements"`
	FirstPage        bool `json:"firstPage"`
	LastPage         bool `json:"lastPage"`
	TotalPages       int  `json:"totalPages"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
}

type Entries struct {
	Page
	Content []Entry `json:"content"`
}

func apiUrl() string {
	return os.Getenv("API_URL")
}

func Feed(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "BLOG.IK.AM",
		Link:        &feeds.Link{Href: "https://blog.ik.am"},
		Description: "maki's memo",
		Author:      &feeds.Author{Name: "Toshiaki Maki", Email: "makingx@gmail.com"},
		Created:     now,
	}

	url := apiUrl()
	log.Printf("Calling %s", url)
	client := http.Client{
		Timeout: time.Second * 10, // Maximum of 10 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "blog-feed")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	entries := Entries{}
	err = json.Unmarshal(body, &entries)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries.Content {
		f := &feeds.Item{
			Id:      strconv.Itoa(entry.EntryId),
			Title:   entry.FrontMatter.Title,
			Link:    &feeds.Link{Href: "https://blog.ik.am/entries/" + strconv.Itoa(entry.EntryId)},
			Author:  &feeds.Author{Name: entry.Created.Name},
			Created: entry.Created.Date,
			Updated: entry.Updated.Date,
		}
		feed.Add(f)
	}

	rss, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, rss)
}
