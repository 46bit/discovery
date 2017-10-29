package main

import (
  "log"
  "time"
  "fmt"
  "net/http"
)

func main() {
  log.Println("Long running!")

  go func() {
    for i := 0; i < 3600; i++ {
      log.Println(i)
      time.Sleep(time.Second)
    }
  }()

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
  })
  http.ListenAndServe(":8080", nil)
}
