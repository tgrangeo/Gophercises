package main

import (
	"Adventure/story"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type templateHandler struct {
	book map[string]story.StoryArc
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	arc := strings.TrimPrefix(r.URL.Path, "/")
	if arc == "" {
		arc = "intro"
	}

	arcData, ok := th.book[arc]
	if !ok {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("template/template.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, arcData)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func main() {
	data, error := ioutil.ReadFile("story.json")
	if error != nil {
		panic(error)
	}

	var book map[string]story.StoryArc
	json.Unmarshal(data, &book)

	bookHandler := &templateHandler{book: book}
	http.Handle("/", bookHandler)
	fmt.Println("listenning on port 8080...")
	http.ListenAndServe(":8080", nil)
}
