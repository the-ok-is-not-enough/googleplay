package main

import (
   "fmt"
   "testing"
)

func TestFilename(t *testing.T) {
   dirs := []string{"", "one"}
   for _, dir := range dirs {
      name := filename(dir, "two", "three", 4)
      fmt.Println(name)
   }
}
