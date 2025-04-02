package main

import (
	"fmt"
	"net/http"

	"currencyservice/internal/repo"
)

func main() {
	repository, err := repo.NewDB()
	defer repository.Close()

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world")
	})

	http.ListenAndServe(":80", nil)
}
