package googleplay

import (
   "bytes"
   "fmt"
   "github.com/89z/googleplay"
   "net/url"
   "testing"
   "time"
)

const device = "3cac75f5da7c75d3"

var _ = time.Second

func TestUpload(t *testing.T) {
   tok := googleplay.Token{
      url.Values{
         "Token": {token},
      },
   }
   auth, err := tok.Auth()
   if err != nil {
      t.Fatal(err)
   }
   /*
   device, err := newCheckin()
   if err != nil {
      t.Fatal(err)
   }
   if err := upload(auth.Get("Auth"), device.String()); err != nil {
      t.Fatal(err)
   }
   fmt.Println(device)
   time.Sleep(9 * time.Second)
   det, err := auth.Details(device.String(), "com.google.android.youtube")
   */
   det, err := auth.Details(device, "com.google.android.youtube")
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
