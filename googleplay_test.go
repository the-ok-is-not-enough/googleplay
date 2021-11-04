package googleplay

import (
   "bytes"
   "fmt"
   "net/url"
   "testing"
   "time"
)

const email = "srpen6@gmail.com"

func TestCheckin(t *testing.T) {
   check, err := NewCheckin().Post()
   if err != nil {
      t.Fatal(err)
   }
   oauth := OAuth{
      url.Values{
         "Auth": {auth},
      },
   }
   if err := NewDevice().Upload(check, oauth); err != nil {
      t.Fatal(err)
   }
   fmt.Println(check)
   time.Sleep(16 * time.Second)
   det, err := oauth.Details(check.String(), "com.google.android.youtube")
   if err != nil {
      t.Fatal(err)
   }
   vers := []string{"16.", "16.4", "16.43.", "16.43.3", "16.43.34"}
   for _, ver := range vers {
      if bytes.Contains(det, []byte(ver)) {
         fmt.Println("pass", ver)
      } else {
         fmt.Println("fail", ver)
      }
   }
}

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
