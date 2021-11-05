package googleplay

import (
   "bytes"
   "fmt"
   "os"
   "testing"
   "time"
)

func TestCheckinDecode(t *testing.T) {
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
   a, err := tok.OAuth()
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
   det, err := a.Details(check.String(), "com.google.android.youtube")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", det)
   vers := []string{"16.", "16.4", "16.43.", "16.43.3", "16.43.34"}
   for _, ver := range vers {
      if bytes.Contains(det, []byte(ver)) {
         fmt.Println("pass", ver)
      } else {
         fmt.Println("fail", ver)
      }
   }
}

func TestCheckinEncode(t *testing.T) {
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
   a, err := tok.OAuth()
   if err != nil {
      t.Fatal(err)
   }
   c, err := NewCheckinRequest().Post()
   if err != nil {
      t.Fatal(err)
   }
   if err := NewDevice().Upload(c.String(), a.Get("Auth")); err != nil {
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
