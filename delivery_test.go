package googleplay

import (
   "fmt"
   "os"
   "testing"
)

func TestDelivery(t *testing.T) {
   head, err := newHeader("googleplay/phone.json")
   if err != nil {
      t.Fatal(err)
   }
   del, err := head.Delivery("com.google.android.youtube", 1524221376)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", del)
}

func newHeader(device string) (*Header, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, err
   }
   tok, err := OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   phone, err := OpenDevice(cache, device)
   if err != nil {
      return nil, err
   }
   return tok.Header(phone)
}
