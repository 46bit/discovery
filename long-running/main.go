package main

import (
  "log"
  "time"
)

func main() {
  log.Println("Long running!")
  for i := 0; i < 3600; i++ {
    fmt.Println(i)
    time.Sleep(time.Second)
  }
}
