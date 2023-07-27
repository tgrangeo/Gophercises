package main

import (
	"errors"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.URL.Path
		if targetURL, ok := pathsToUrls[req]; ok {
			http.Redirect(w, r, targetURL, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	if yml == nil {
		return nil, errors.New("YAML data is not provided")
	}

	pathsToUrls := make(map[string]string)
	err := yaml.Unmarshal(yml, &pathsToUrls)
	if err != nil {
		return nil, err
	}

	return MapHandler(pathsToUrls, fallback), nil
}
