package googleplay

import (
   "fmt"
   "os"
   "testing"
)

func TestAuth(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   // read token
   tok := new(Token)
   {
      r, err := os.Open(cache + "/googleplay/token.json")
      if err != nil {
         t.Fatal(err)
      }
      defer r.Close()
      if err := tok.Decode(r); err != nil {
         t.Fatal(err)
      }
   }
   o, err := tok.Auth()
   if err != nil {
      t.Fatal(err)
   }
   // read checkin
   check := new(Checkin)
   {
      r, err := os.Open(cache + "/googleplay/checkin.json")
      if err != nil {
         t.Fatal(err)
      }
      defer r.Close()
      if err := check.Decode(r); err != nil {
         t.Fatal(err)
      }
   }
   // details
   det, err := o.Details(check.String(), "com.google.android.youtube")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", det)
}
