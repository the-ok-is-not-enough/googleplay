package discord

import (
   "fmt"
   "testing"
)

const app = "com.discord"

func TestDiscord(t *testing.T) {
   det, err := details(app)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(det)
}
