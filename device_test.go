package googleplay

import (
   "os"
   "testing"
   "time"
)

func TestCheckinArmeabi(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   device, err := Phone.Checkin(Armeabi)
   if err != nil {
      t.Fatal(err)
   }
   if err := device.Create(cache, "googleplay/armeabi.json"); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}

func TestCheckinX86(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   device, err := Phone.Checkin(X86)
   if err != nil {
      t.Fatal(err)
   }
   if err := device.Create(cache, "googleplay/x86.json"); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}
