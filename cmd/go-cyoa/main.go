package main

import (
	"flag"
	"fmt"
	"github.com/eitah/go-cyoa"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Printf("error!: %s", err)
		os.Exit(1)
	}
}

func mainErr() error{
port := flag.Int("port", 3000, "the port to start the CYOA Web Application on.")
	filename := flag.String("file", "example-story.json", "the json file with the current adventure")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		return fmt.Errorf("opening story: %w", err)
	}

	adventure, err := cyoa.JsonStory(f)
	if err != nil {
		return fmt.Errorf("json story: %w", err)
	}

	storyHandler := cyoa.NewHandler(adventure, nil)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), storyHandler))
	return nil
}


