package googleplay

import (
   "os"
   "testing"
   "time"
)

func TestDevice(t *testing.T) {
   dev, err := DefaultConfig.Device()
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   if err := dev.Create(cache, "googleplay/device.json"); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}
