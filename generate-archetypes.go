// This go script generates new default architypes for the recipe sub-directory cateorgies.
//
// usage: go run generate-archetypes.go

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
		if strings.Contains(a, i) {
			return true
		}
	}
	return false
}

func loadDefaultArchetypes() string {
	// load default archetype
	defaultArchetype, err := ioutil.ReadFile("archetypes/default.md")
	if err != nil {
		log.Fatalln("Could not load default archetype file [archetypes/default.md]", err)
	}
	return string(defaultArchetype)
}

func currentArchetypes() []string {
	// get current archetypes

	var archetypes []string

	walkPath := "archetypes"
	err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), ".DS_Store") {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		archetypes = append(archetypes, path)
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return archetypes
}

func archetypesToCreate() []string {

	archetypes := currentArchetypes()
	log.Println("Archetypes that already exist:", archetypes)

	// get all the recipe categories that do not already have archetypes
	var files []string

	walkPath := "content/recipes/"
	err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {

		// only make new archetypes for directories
		if !info.IsDir() {
			return nil
		}
		// or if it's the top level dir
		// TODO: there's definitely a better way to do this
		if info.Name() == "content/recipes/" {
			return nil
		}
		if contains(archetypes, info.Name()) {
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
	log.Println("Generating archetypes")

	newArchetypes := archetypesToCreate()

	defaultArchetype := loadDefaultArchetypes()

	// for each file, create a new archetype file
	for _, newArchetype := range newArchetypes {
		log.Println("Creating archetype file for", newArchetype)
		newArchetype := strings.Split(newArchetype, "/")[2]

		// Replace "default" with the name of the new archetype
		file := strings.Replace(defaultArchetype, "default", newArchetype, -1)

		err := ioutil.WriteFile("archetypes/"+newArchetype+".md", []byte(file), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("done")
}
