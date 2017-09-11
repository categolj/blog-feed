package main

import (
	"encoding/json"
	"fmt"
	"github.com/categolj/blog-feed/spring_cloud_services"
	"github.com/gorilla/feeds"
	"github.com/ryanmoran/viron"

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

type Environment struct {
	VCAPServices struct {
		ServiceRegistry []struct {
			Credentials struct {
				RegistryURI    string `json:"uri"`
				ClientSecret   string `json:"client_secret"`
				ClientID       string `json:"client_id"`
				AccessTokenURI string `json:"access_token_uri"`
			} `json:"credentials"`
		} `json:"p-service-registry"`
	} `env:"VCAP_SERVICES" env-required:"true"`
}

func main() {
	http.HandleFunc("/", feed)

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "4000"
	}
	log.Printf(fmt.Sprintf("Listening at %s", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func apiUrl() string {
	var env Environment
	err := viron.Parse(&env)
	if err != nil {
		log.Printf("Unable to parse environment (%s). Falling back to default url", err)
		var url string
		if url = os.Getenv("API_URL"); len(url) == 0 {
			return "https://blog-api.cfapps.pez.pivotal.io/api/entries"
		}
		return url
	}
	serviceCredentials := env.VCAPServices.ServiceRegistry[0].Credentials
	uaaClient := &spring_cloud_services.UAAClient{
		BaseURL: serviceCredentials.AccessTokenURI,
		Name:    serviceCredentials.ClientID,
		Secret:  serviceCredentials.ClientSecret,
	}
	eurekaClient := &spring_cloud_services.EurekaClient{
		BaseURL:    serviceCredentials.RegistryURI,
		HttpClient: http.DefaultClient,
		UAAClient:  uaaClient,
	}
	destination, err := eurekaClient.GetAppByName("blog-api")
	if err != nil {
		log.Fatalf("Unable to get application by name: %s", err)
	}
	return fmt.Sprintf("https://%s:%d/api/entries", destination.HostName, destination.SecurePort)
}

func feed(w http.ResponseWriter, r *http.Request) {
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

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	entries := Entries{}
	jsonErr := json.Unmarshal(body, &entries)
	if jsonErr != nil {
		log.Fatal(jsonErr)
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
