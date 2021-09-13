package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var adventure Story
	err := d.Decode(&adventure)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return adventure, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title   string       `json:"title"`
	Paragraphs   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Chapter  string `json:"arc"`
}

func NewHandler(s Story) http.Handler {
	return handler{
		s: s,
	}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	// "/intro" => "intro"
	path = path[1:]

	//					"["intro"]"
	if chapter, ok := h.s[path]; ok {
		if err := tpl.Execute(w, chapter); err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

var defaultHandlerTemplate = `<!DOCTYPE html>
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
				<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
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