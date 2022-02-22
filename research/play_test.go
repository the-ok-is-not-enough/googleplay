package googleplay

import (
   "fmt"
   "testing"
)

func TestCategory(t *testing.T) {
   docs, err := list("FINANCE")
   if err != nil {
      t.Fatal(err)
   }
   for _, doc := range docs {
      fmt.Print(doc, "\n---\n")
   }
}
