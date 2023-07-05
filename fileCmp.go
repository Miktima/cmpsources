package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func getHash(filename string) (uint32, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	h := crc32.NewIEEE()
	h.Write(bs)
	return h.Sum32(), nil
}

func main() {
	// initial values
	urlPage := "https://sputnikglobe.com/20230703/sweden-belatedly-denounces-quran-burning-after-international-backlash-1111632481.html"
	cdn := "cdn1.img.sputnikglobe.com"
	initGitPath := "C:\\dev\\sputnik-white\\htdocs\\"
	// Send an HTTP GET request to the urlPage web page
	resp, err := http.Get(urlPage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	// find all matched links
	re := regexp.MustCompile(`(https://` + cdn + `[\w/.:]*(css|js))`)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	content := string(body)
	links := re.FindAllString(content, -1)
	// loops through the links slice and find corresponded files in git
	// Check if paths must be changed (for windows)
	changePath := strings.Count(initGitPath, "\\")
	for _, l := range links {
		fmt.Println("link - ", l)
		before, file, found := strings.Cut(l, "https://"+cdn+"/")
		if !found {
			fmt.Println("Error: pattern not found!")
			// Just to avoid error
			fmt.Println(before)
		}
		if changePath > 0 {
			file = strings.ReplaceAll(file, "/", "\\")
		}
		fmt.Println("File: ", initGitPath+file)
		h1, err := getHash(initGitPath + file)
		if err != nil {
			return
		}
		fmt.Println(initGitPath+file, h1)
	}

	// Find and print all links on the web page
	/*	var links []string
		var link func(*html.Node)
		link = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						// adds a new link entry when the attribute matches
						links = append(links, a.Val)
					}
				}
			}

			// traverses the HTML of the webpage from the first child node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				link(c)
			}
		}
		link(doc)

		// loops through the links slice
		for _, l := range links {
			fmt.Println("Link:", l)
		}*/
	/*h1, err := getHash("test1.txt")
	if err != nil {
		return
	}
	h2, err := getHash("test2.txt")
	if err != nil {
		return
	}
	fmt.Println(h1, h2, h1 == h2)*/
}
