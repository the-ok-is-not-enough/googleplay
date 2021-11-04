package googleplay

import (
   "bytes"
   "fmt"
   "net/url"
   "testing"
)

const (
   app = "com.google.android.youtube"
   email = "srpen6@gmail.com"
)

func TestDetails(t *testing.T) {
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
   det, err := auth.Details(device, app)
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

func TestToken(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(tok)
}
