package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getStories(toGet []string) map[string][][]string {
	resp, err := http.Get("https://magic.wizards.com/en/story")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	volumes := strings.Split(string(body), "<div class=\"gridder-list module_story-archive__grid-item\" ")[1:]
	for i := range volumes {
		volumes[i] = strings.Split(volumes[i], "<p>")[1]
		volumes[i] = strings.Split(volumes[i], "</p>")[0]
	}
	stories := map[string][][]string{}
	chapters := strings.Split(string(body), "\" class=\"gridder-content module_story-archive__grid-item-content\">")[1:]
	for i := range chapters {
		if !contains(toGet, volumes[i]) {
			continue
		}
		stories[volumes[i]] = [][]string{}
		volume_chapters := strings.Split(chapters[i], "<a href=\"")[1:]
		for j := range volume_chapters {
			parts := strings.Split(volume_chapters[j], "<h3>")
			parts[0] = strings.Split(parts[0], "\"")[0]
			parts[1] = strings.Split(parts[1], "</h3>")[0]
			stories[volumes[i]] = append(stories[volumes[i]], parts)
		}
	}

	return stories
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getContent(name string, article []string) string {
	content := "<h1>" + name + "</h1><h2>" + article[1] + "</h2>"
	resp, err := http.Get("https://magic.wizards.com" + article[0])
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	content += strings.Split(string(body), "<div id=\"content-detail-page-of-an-article\">")[1]
	content = strings.Split(content, "</body></html>")[0]
	content = strings.ReplaceAll(content, "<html><body>", "")
	return content
}

func writeFile(series, title, content string) {
	series = strings.ReplaceAll(series, ":", "")
	title = strings.ReplaceAll(title, ":", "")
	f, err := os.Create("html/" + series + " " + title + ".html")
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func convert() {
	files, err := ioutil.ReadDir("html")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		cmd := exec.Command("ebook-convert", "html/"+file.Name(), "epub/"+strings.ReplaceAll(file.Name(), ".html", ".epub"))
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%s\n", stdoutStderr)
		cmd = exec.Command("ebook-convert", "html/"+file.Name(), "azw3/"+strings.ReplaceAll(file.Name(), ".html", ".azw3"))
		stdoutStderr, err = cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%s\n", stdoutStderr)
		cmd = exec.Command("ebook-convert", "html/"+file.Name(), "mobi/"+strings.ReplaceAll(file.Name(), ".html", ".mobi"))
		stdoutStderr, err = cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%s\n", stdoutStderr)
	}
}

func main() {
	stories_to_get := []string{
		"Strixhaven: School of Mages",
		"Kaldheim",
		"Zendikar Rising",
	}
	stories := getStories(stories_to_get)

	for k := range stories {
		fmt.Println(k)
		for i := range stories[k] {
			content := getContent(k, stories[k][i])
			writeFile(k, strconv.Itoa(i)+". "+stories[k][i][1], content)
		}
	}

	convert()
}
