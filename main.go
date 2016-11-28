package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const versionInfo = "0.0.2"

type track struct {
	ID       string    `json:"id"`
	Language string    `json:"language"`
	Active   bool      `json:"active"`
	Problems []problem `json:"problems"`
}

func (t track) has(exercise string) bool {
	if exercise == "" {
		return true
	}
	for _, p := range t.Problems {
		if p.Slug == exercise {
			return true
		}
	}
	return false
}

type problem struct {
	Slug string `json:"slug"`
}

type ticket struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

func (t *ticket) addLabels(s string) {
	if s == "" {
		return
	}
	t.Labels = strings.Split(s, ",")
	for i, label := range t.Labels {
		t.Labels[i] = strings.Trim(label, " ")
	}
}

var (
	title    = flag.String("title", "", "The title of your issue.")
	body     = flag.String("body", "", "The body of your issue.")
	file     = flag.String("file", "", "The file containing the title and body of the issue (overrides -title and -body).")
	labels   = flag.String("labels", "", "A comma-separated list of labels to add.")
	exercise = flag.String("exercise", "", "The slug of the relevant exercise (optional). If no exercise is passed, the issue will be submitted to all active tracks.")
	yes      = flag.Bool("yes", false, "Actually submit issues. Without this flag, a dry run is performed.")
	version  = flag.Bool("version", false, "Version of binary.")
)

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("Blazon version: %s\n", versionInfo)
		os.Exit(0)
	}
	if *file == "" && (*title == "" || *body == "") {
		flag.Usage()
		os.Exit(1)
	}

	t := &ticket{
		Title:  *title,
		Body:   *body,
		Labels: []string{}, // api doesn't handle null values here
	}

	if *file != "" {
		b, err := ioutil.ReadFile(*file)
		if err != nil {
			log.Fatal(err)
		}

		i := bytes.Index(b, []byte("\n"))
		if i == -1 {
			log.Fatal("First line of file must contain the subject line of the issue.")
		}

		t.Title = string(b[:i])
		t.Body = string(b[i+1:])
	}
	t.Body = strings.Trim(t.Body, "\n")

	t.addLabels(*labels)

	postBody, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Get("http://x.exercism.io/v3/tracks")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var pld struct {
		Tracks []track `json:"tracks"`
	}
	if err := json.NewDecoder(res.Body).Decode(&pld); err != nil {
		log.Fatal(err)
	}

	if !*yes {
		fmt.Printf("Title:\n%s\n\nBody:\n%s\n\n", t.Title, t.Body)
		fmt.Println("This issue would be submitted to the following tracks:")
	}

	for _, track := range pld.Tracks {
		if !track.has(*exercise) {
			continue
		}

		if !*yes {
			var status string
			if !track.Active {
				status = " (inactive)"
			}

			fmt.Printf("- %s%s\n", track.ID, status)
			continue
		}
		time.Sleep(1200 * time.Millisecond)

		r := bytes.NewReader(postBody)

		url := fmt.Sprintf("https://api.github.com/repos/exercism/x%s/issues", track.ID)
		req, err := http.NewRequest("POST", url, r)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "exercism/blazon")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("BLAZON_GITHUB_API_TOKEN")))

		c := &http.Client{}
		res, err := c.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
			var pld struct {
				URL string `json:"html_url"`
			}
			if err := json.NewDecoder(res.Body).Decode(&pld); err != nil {
				log.Printf("%s %s", track.Language, err)
				continue
			}

			fmt.Println(pld.URL)
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("%s %s", track.Language, err)
			continue
		}
		log.Printf("ERROR (%d) %s\n  %s", res.StatusCode, track.Language, body)
	}

	if !*yes {
		fmt.Println("\nTo submit, rerun the command with the -yes flag.")
	}
}
