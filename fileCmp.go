package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func getHash(filename, url string) (uint32, uint32, bool, error) {
	var hFileSum32 chan uint32 = make(chan uint32)
	var hURLSum32 chan uint32 = make(chan uint32)
	go func() {
		bs, err := os.ReadFile(filename)
		if err != nil {
			return
		}
		hFile := crc32.NewIEEE()
		hFile.Write(bs)
		hFileSum32 <- hFile.Sum32()
	}()
	h1 := <-hFileSum32
	go func() {
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		hUrl := crc32.NewIEEE()
		hUrl.Write(body)
		hURLSum32 <- hUrl.Sum32()
	}()
	h2 := <-hURLSum32
	return h1, h2, h1 == h2, nil
}

func main() {
	// initial values
	var urlPage string
	var cdn string
	var initGitPath string
	fmt.Printf("Page URL to check css and js resources: ")
	fmt.Scanln(&urlPage)
	fmt.Printf("CDN on the page: ")
	fmt.Scanln(&cdn)
	fmt.Printf("PATH to htdocs of a project: ")
	fmt.Scanln(&initGitPath)
	// urlPage = "https://sputniknews.lat/20230706/la-policia-ucraniana-usa-la-fuerza-contra-los-fieles-del-monasterio-de-las-cuevas-de-kiev--videos--1141291678.html"
	// cdn = "cdn(1|2).img.sputniknews.lat"
	// initGitPath = "/home/vboxuser/dev/sputnik-white/htdocs/"
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
	startTime := time.Now()
	for _, l := range links {
		reFile := regexp.MustCompile(`https://` + cdn + `(/[\w/.:]*(css|js))`)
		paths := reFile.FindStringSubmatch(l)
		file := paths[len(paths)-2]
		if changePath > 0 {
			file = strings.ReplaceAll(file, "/", "\\")
		}
		h1, h2, ver, err := getHash(initGitPath+file, l)
		if err != nil {
			return
		}
		switch ver {
		case true:
			fmt.Println(l, h1, h2, ver)
		case false:
			fmt.Println(l, h1, h2, ver, "!!!!!!!")
		}
	}
	fmt.Println("Elapsed time: ", time.Since(startTime))
}
