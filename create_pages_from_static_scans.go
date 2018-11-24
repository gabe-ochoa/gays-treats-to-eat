// This go script creates pages for each picture in the static/Recipes folder
//
// usage: go run create_pages_from_scans.go

package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func contains(s []string, i string) bool {
	for _, a := range s {
		// split the .md off
		if strings.Contains(strings.Split(a, ".")[0], i) {
			return true
		}
	}
	return false
}

func loadArchetype(archetype string) string {
	// load default archetype
	file, err := ioutil.ReadFile("archetypes/" + archetype + ".md")
	if err != nil {
		log.Fatalln("Could not load archetype file", archetype, err)
	}
	return string(file)
}

func getAllArchetypes() []os.FileInfo {
	files, err := ioutil.ReadDir("archetypes/")
	if err != nil {
		log.Fatalln("Could not load archetype file", err)
	}
	return files
}

func currentPages() []string {
	// get current pages

	var pages []string

	walkPath := "content/"
	err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), ".DS_Store") {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		pages = append(pages, path)
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return pages
}

func newRecipes() []string {

	pages := currentPages()
	log.Println("Pages that already exist:", pages)

	// get all the Recipe categories that do not already have archetypes
	var files []string

	walkPath := "static/"
	err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {

		// don't include directories
		if info.IsDir() {
			return nil
		}
		// or if it's the top level dir
		// TODO: there's definitely a better way to do this
		if info.Name() == "static/Recipes/" {
			return nil
		}
		if contains(files, info.Name()) {
			return nil
		}
		// or add it to files we need to create
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return files[1:]
}

func main() {
	log.Println("Generating new pages for scanned Recipes")
	for _, r := range newRecipes() {
		r := strings.Split(r, "/")
		archetype, Recipe := r[1], strings.Split(r[2], ".")[0]

		// TODO: This does not work. Run the Hugo command here and then string replace in the generated file
		// load and replace Recipe name in template
		path := archetype + "/" + Recipe + ".md"
		log.Printf("Generating : %s \n", path)
		cmd := exec.Command("/usr/local/bin/hugo", "new", path)
		bytes, err := cmd.Output()
		log.Print(string(bytes))
		if err != nil {
			log.Fatalf("Error running hugo command: %s", err.Error())
		}
	}
	log.Println("done")
}
