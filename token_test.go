package googleplay

import (
   "os"
   "testing"
)

const email = "srpen6@gmail.com"

func TestTokenEncode(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   c, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   c += "/googleplay"
   os.Mkdir(c, os.ModeDir)
   f, err := os.Create(c + "/token.json")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   if err := tok.Encode(f); err != nil {
      t.Fatal(err)
   }
}
