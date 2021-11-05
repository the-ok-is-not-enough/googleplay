package googleplay

import (
   "os"
   "testing"
   "time"
)

func TestCheckinEncode(t *testing.T) {
   // get Auth
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
   o, err := tok.Auth()
   if err != nil {
      t.Fatal(err)
   }
   // get Checkin
   c, err := NewCheckinRequest().Post()
   if err != nil {
      t.Fatal(err)
   }
   // Upload
   if err := o.Upload(c.String(), NewDevice()); err != nil {
      t.Fatal(err)
   }
   w, err := os.Create(cache + "/googleplay/checkin.json")
   if err != nil {
      t.Fatal(err)
   }
   defer w.Close()
   if err := c.Encode(w); err != nil {
      t.Fatal(err)
   }
   // make sure server processes request
   time.Sleep(16 * time.Second)
}
