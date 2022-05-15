package googleplay

import (
   "fmt"
   "os"
   "testing"
)

/*
1 APK:
kr.sira.metal

4 APK:
com.pinterest

1 OBB:
com.PirateBayGames.ZombieDefense2

2 OBB:
com.sigmateam.alienshootermobile.free
*/
func TestDelivery(t *testing.T) {
   head, err := newHeader("googleplay/x86.json")
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
   token, err := OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   dev, err := OpenDevice(cache, device)
   if err != nil {
      return nil, err
   }
   return token.Header(dev)
}
