package googleplay

import (
   "fmt"
   "testing"
)

func TestCategory(t *testing.T) {
   docs, err := topChartItems("FINANCE")
   if err != nil {
      t.Fatal(err)
   }
   for _, doc := range docs {
      fmt.Print(doc, "\n---\n")
   }
}
