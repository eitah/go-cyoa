package main

import (
	"flag"
	"fmt"
	"github.com/eitah/go-cyoa"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Printf("error!: %s", err)
		os.Exit(1)
	}
}

func mainErr() error {
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

	tpl := template.Must(template.New("").Parse(storyTemplate))
	h := cyoa.NewHandler(adventure, cyoa.WithPathFunction(pathFn), cyoa.WithTemplate(tpl))

	// We use a mux here to make sure that only traffic to /story hits this site sub-section.
	mux := http.NewServeMux()
	mux.Handle("/story/", h)

	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
	return nil
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	// "/story/intro" => "intro"
	return path[len("/story/"):]
}

var storyTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Choose your own adventure</title>
</head>
<body>
	<section class="page">
		<h1>{{.Title}}</h1>
		{{range .Paragraphs}}
			<p>{{.}}</p>
		{{end}}
		<ul>
			{{range .Options}}
				<li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
			{{end}}
		</ul>
	</section>
	<style>
        body {
            font-family: Helvetica, Arial;
        }

        h1 {
            text-align: center;
            position: relative;
        }

        .page {
            width: 80%;
            max-width: 500px;
            margin: auto;
            margin-top: 40px;
            margin-bottom: 40px;
            padding: 80px;
            background: #FFFCF6;
            border: 1px solid #eee;
            box-shadow: 0 10px 6px -6px #777;
        }

        ul {
            border-top: 1px dotted #ccc;
            padding: 10px 0 0 0;
            -webkit-padding-start: 0;
        }

        li {
            padding-top: 10px
        }

        a,
        a:visited {
            text-decoration: none;
            color: #6295b5;
        }

        a:active,
        a:hover {
            color: #7792a2;
        }

        p {
            text-indent: 1em;
        }
    </style>
</body>
</html>`
