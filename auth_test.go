package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

var apps = map[string]int{
   "com.axis.drawingdesk.v3": 190,
   "com.google.android.youtube": 1524221376,
   "com.instagram.android": 321403734,
   "com.pinterest": 9398020,
   "com.smarty.voomvoom": 369,
   "com.vimeo.android.videoapp": 3510004,
   "org.videolan.vlc": 13040207,
}

func TestDetails(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev := new(Device)
   r, err := os.Open(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   defer r.Close()
   if err := dev.Decode(r); err != nil {
      t.Fatal(err)
   }
   for app := range apps {
      det, err := auth.Details(dev, app)
      if err != nil {
         t.Fatal(err)
      }
      if det.DocV2.Details.AppDetails.VersionCode == 0 {
         t.Fatal(app)
      }
      fmt.Println(det.DocV2.Details.AppDetails.VersionCode)
      fmt.Println(det.DocV2.Details.AppDetails.InstallationSize)
      time.Sleep(time.Second)
   }
}

func TestDelivery(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev := new(Device)
   r, err := os.Open(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   defer r.Close()
   if err := dev.Decode(r); err != nil {
      t.Fatal(err)
   }
   app := "org.videolan.vlc"
   del, err := auth.Delivery(dev, app, apps[app])
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", del)
}

func getAuth() (*Auth, string, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, "", err
   }
   tok := new(Token)
   r, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      return nil, "", err
   }
   defer r.Close()
   if err := tok.Decode(r); err != nil {
      return nil, "", err
   }
   auth, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return auth, cache, nil
}


