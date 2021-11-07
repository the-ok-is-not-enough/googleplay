package googleplay

import (
   "os"
   "testing"
   "time"
)

func TestDevice(t *testing.T) {
   tok := new(Token)
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   r, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      t.Fatal(err)
   }
   defer r.Close()
   if err := tok.Decode(r); err != nil {
      t.Fatal(err)
   }
   auth, err := tok.Auth()
   if err != nil {
      t.Fatal(err)
   }
   dev, err := NewDevice(DefaultCheckin)
   if err != nil {
      t.Fatal(err)
   }
   if err := auth.Upload(dev, DefaultConfig); err != nil {
      t.Fatal(err)
   }
   w, err := os.Create(cache + "/googleplay/checkin.json")
   if err != nil {
      t.Fatal(err)
   }
   defer w.Close()
   if err := dev.Encode(w); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}
