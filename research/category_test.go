package googleplay

import (
   "fmt"
   "testing"
   gp "github.com/89z/googleplay"
)

func TestCategory(t *testing.T) {
   tok, err := gp.OpenToken("ignore/token.json")
   if err != nil {
      t.Fatal(err)
   }
   dev, err := gp.OpenDevice("ignore/device.json")
   if err != nil {
      t.Fatal(err)
   }
   head, err := tok.Header(dev)
   if err != nil {
      t.Fatal(err)
   }
   docs, err := Header{head.Header}.Category("FINANCE", 9)
   if err != nil {
      t.Fatal(err)
   }
   for _, doc := range docs {
      fmt.Printf("%+v\n", doc)
   }
}
