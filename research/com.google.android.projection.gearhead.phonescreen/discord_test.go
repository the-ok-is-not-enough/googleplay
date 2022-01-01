package discord

import (
   "fmt"
   "testing"
)

func TestDiscord(t *testing.T) {
   det, err := details()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(det)
}
