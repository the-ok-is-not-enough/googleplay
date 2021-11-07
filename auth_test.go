package googleplay

import (
   "fmt"
   "os"
   "testing"
)

const (
   app = "com.google.android.youtube"
   ver = 1524094400
)

func TestDelivery(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev := new(Device)
   r, err := os.Open(cache + "/googleplay/checkin.json")
   if err != nil {
      t.Fatal(err)
   }
   defer r.Close()
   if err := dev.Decode(r); err != nil {
      t.Fatal(err)
   }
   del, err := auth.Delivery(dev, app, ver)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", del)
}

func TestDetails(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev := new(Device)
   r, err := os.Open(cache + "/googleplay/checkin.json")
   if err != nil {
      t.Fatal(err)
   }
   defer r.Close()
   if err := dev.Decode(r); err != nil {
      t.Fatal(err)
   }
   det, err := auth.Details(dev, app)
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
   auth, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return auth, cache, nil
}


