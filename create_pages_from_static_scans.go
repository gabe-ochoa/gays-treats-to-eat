// This go script creates pages for each picture in the static/recipes folder
//
// usage: go run create_pages_from_scans.go

package main

import (
	"io/ioutil"
	"log"
	"os"
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

func currentPages() []string {
	// get current pages

	var pages []string

	walkPath := "content/recipes/"
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

	// get all the recipe categories that do not already have archetypes
	var files []string

	walkPath := "static/recipes/"
	err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {

		// don't include directories
		if info.IsDir() {
			return nil
		}
		// or if it's the top level dir
		// TODO: there's definitely a better way to do this
		if info.Name() == "static/recipes/" {
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
	log.Println("Generating new pages for scanned recipes")

	for _, r := range newRecipes() {
		r := strings.Split(r, "/")
		archetype, recipe := r[2], strings.Split(r[3], ".")[0]

		// load and replace recipe name in template
		template := loadArchetype(archetype)
		file := strings.Replace(template, "recipe", recipe, -1)

		// write file
		err := ioutil.WriteFile("content/recipes/"+archetype+"/"+recipe+".md", []byte(file), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("done")
}
