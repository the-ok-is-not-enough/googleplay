package googleplay

import (
   "os"
   "testing"
   "time"
)

func TestTabletCheckin(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   tablet, err := Tablet.Checkin()
   if err != nil {
      t.Fatal(err)
   }
   if err := tablet.Create(cache, "googleplay/tablet.json"); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}

func TestTvCheckin(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   tv, err := TV.Checkin()
   if err != nil {
      t.Fatal(err)
   }
   if err := tv.Create(cache, "googleplay/tv.json"); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}

func TestPhoneCheckin(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   phone, err := Phone.Checkin()
   if err != nil {
      t.Fatal(err)
   }
   if err := phone.Create(cache, "googleplay/phone.json"); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}
