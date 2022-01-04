package googleplay

import (
   "fmt"
   "testing"
   "time"
)

var apps = []string{
   "com.google.android.projection.gearhead.phonescreen",
   "com.google.android.youtube",
   "com.instagram.android",
   "com.miui.weather2",
   "com.pinterest",
   "com.valvesoftware.android.steam.community",
}

func TestDelivery(t *testing.T) {
   dev, err := Checkin(DefaultConfig)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(dev)
   time.Sleep(4 * time.Second)
   del, err := Auth{auth}.Delivery(dev, "com.pinterest", 9428030)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", del)
}

func TestDetails(t *testing.T) {
   dev, err := Checkin(DefaultConfig)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(dev)
   time.Sleep(4 * time.Second)
   for _, app := range apps {
      det, err := Auth{auth}.Details(dev, app)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", det)
      time.Sleep(time.Second)
   }
}
