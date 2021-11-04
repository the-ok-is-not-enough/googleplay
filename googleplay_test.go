package googleplay

import (
   "fmt"
   "net/url"
   "testing"
)

const email = "srpen6@gmail.com"

func TestAuth(t *testing.T) {
   tok := Token{
      url.Values{
         "Token": {token},
      },
   }
   auth, err := tok.Auth()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(auth)
}

func TestToken(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(tok)
}
