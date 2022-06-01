package search

import (
   "testing"
)

func TestSearch(t *testing.T) {
   err := search("Taxi")
   if err != nil {
      t.Fatal(err)
   }
}
