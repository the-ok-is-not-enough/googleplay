package googleplay

import (
   "os"
   "testing"
   "time"
)

func TestDevice(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   device, err := Phone.Checkin()
   if err != nil {
      t.Fatal(err)
   }
   if err := device.Create(cache, "googleplay/device.json"); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}
