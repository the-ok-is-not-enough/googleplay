package googleplay

import (
   "fmt"
   "testing"
   "time"
)

func TestCategory(t *testing.T) {
   doc, err := list("FINANCE")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(doc)
   time.Sleep(time.Second)
   next, err := doc.list()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(next)
}
