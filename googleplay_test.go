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
   sig, err := Signature(email, password)
   if err != nil {
      t.Fatal(err)
   }
   val, err := Token(email, sig)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(val)
}
