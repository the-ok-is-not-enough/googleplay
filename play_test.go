package googleplay

import (
   "os"
   "testing"
)

func TestToken(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   if err := tok.Create(cache, "googleplay/token.json"); err != nil {
      t.Fatal(err)
   }
}
