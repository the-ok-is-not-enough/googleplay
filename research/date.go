package main

import (
   "fmt"
   "time"
)

var bads = []string{
   "Квітень",
   "квіт",
   "квіт.",
}

func main() {
   for _, bad := range bads {
      _, err := time.Parse("Jan", bad)
      fmt.Print(err, "\n\n")
   }
}
