package main

import (
  "net/http"
  "io"
)

func main() {
  http.HandleFunc("/", servePage)
	http.ListenAndServe(":8080", nil)
}

func servePage(writer http.ResponseWriter, reqest *http.Request) {
  io.WriteString(writer, "Welcome to my Fedora CoreOS Demo. :)<")
}