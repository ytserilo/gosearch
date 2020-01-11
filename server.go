package main

import (
  "net/http"
  "./indexing"
  "./search"
)

func main() {
  http.Handle("/add-index/", http.HandlerFunc(indexing.AddIndex))
  http.Handle("/remove-index/", http.HandlerFunc(indexing.RemoveIndex))

  http.Handle("/search/", http.HandlerFunc(search.Search))

  http.ListenAndServe("0.0.0.0:9200", nil)
}
