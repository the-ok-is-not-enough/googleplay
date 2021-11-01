package googleplay

import (
   "fmt"
   "testing"
)

const (
   app = "com.google.android.youtube"
   device = "38B5418D8683ADBB"
   email = "srpen6@gmail.com"
)

func TestAuth(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(tok)
}
