package googleplay

import (
   "fmt"
   "os"
   "testing"
)

const app = "com.google.android.youtube"

func TestPurchase(t *testing.T) {
   a, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   check := new(Checkin)
   r, err := os.Open(cache + "/googleplay/checkin.json")
   if err != nil {
      t.Fatal(err)
   }
   defer r.Close()
   if err := check.Decode(r); err != nil {
      t.Fatal(err)
   }
   p, err := a.Purchase(check.String(), app)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", p)
}

func TestDetails(t *testing.T) {
   a, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   check := new(Checkin)
   r, err := os.Open(cache + "/googleplay/checkin.json")
   if err != nil {
      t.Fatal(err)
   }
   defer r.Close()
   if err := check.Decode(r); err != nil {
      t.Fatal(err)
   }
   det, err := a.Details(check.String(), app)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", det)
}

func getAuth() (*Auth, string, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, "", err
   }
   tok := new(Token)
   r, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      return nil, "", err
   }
   defer r.Close()
   if err := tok.Decode(r); err != nil {
      return nil, "", err
   }
   a, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return a, cache, nil
}
