package googleplay

import (
   "os"
   "testing"
)

func TestHeader(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   token, err := OpenToken(home, "googleplay/token.json")
   if err != nil {
      t.Fatal(err)
   }
   LogLevel = 1
   if _, err := token.Header(0, false); err != nil {
      t.Fatal(err)
   }
}

func TestToken(t *testing.T) {
   token, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   if err := token.Create(home, "googleplay/token.json"); err != nil {
      t.Fatal(err)
   }
}
