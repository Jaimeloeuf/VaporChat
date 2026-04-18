package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, Gopher!")
	})

	fmt.Println("Server starting on :3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
