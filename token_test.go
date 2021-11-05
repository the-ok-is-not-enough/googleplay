package googleplay

import (
   "fmt"
   "net/url"
   "testing"
)

const email = "srpen6@gmail.com"

func TestOAuth(t *testing.T) {
   tok := Token{
      url.Values{
         "Token": {token},
      },
   }
   oauth, err := tok.OAuth()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(oauth)
}

func TestToken(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(tok)
}
