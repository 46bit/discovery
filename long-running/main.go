package main

import (
  "log"
  "time"
)

func main() {
  log.Println("Long running!")
  for i := 0; i < 3600; i++ {
    log.Println(i)
    time.Sleep(time.Second)
  }
}
