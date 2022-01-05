package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

const email = "srpen6@gmail.com"

type app struct {
   down, id string
   ver int64
}

var apps = []app{
   {down: "10.996 B", id: "com.google.android.youtube", ver: 1524221376},
   {down: "3.932 B", id: "com.instagram.android"},
   {down: "975.149 M", id: "com.miui.weather2"},
   {down: "689.574 M", id: "com.pinterest"},
   {down: "282.147 M", id: "org.videolan.vlc"},
   {down: "95.910 M", id: "org.thoughtcrime.securesms"},
   {down: "77.289 M", id: "com.valvesoftware.android.steam.community"},
   {down: "31.446 M", id: "com.xiaomi.smarthome"},
   {down: "30.702 M", id: "com.vimeo.android.videoapp"},
   {down: "9.419 M", id: "com.axis.drawingdesk.v3"},
   {down: "282.669 K", id: "com.smarty.voomvoom"},
}

func TestDetails(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev := new(Device)
   file, err := os.Open(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   if err := dev.Decode(file); err != nil {
      t.Fatal(err)
   }
   for _, app := range apps {
      det, err := auth.Details(dev, app.id)
      if err != nil {
         t.Fatal(err)
      }
      if det.VersionCode == 0 {
         t.Fatal(det)
      }
      if det.VersionString == "" {
         t.Fatal(det)
      }
      time.Sleep(time.Second)
   }
}

func TestDevice(t *testing.T) {
   dev, err := Checkin(DefaultConfig)
   if err != nil {
      t.Fatal(err)
   }
   tok := new(Token)
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   src, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      t.Fatal(err)
   }
   defer src.Close()
   if err := tok.Decode(src); err != nil {
      t.Fatal(err)
   }
   dst, err := os.Create(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   defer dst.Close()
   if err := dev.Encode(dst); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}

func TestDelivery(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev := new(Device)
   file, err := os.Open(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   if err := dev.Decode(file); err != nil {
      t.Fatal(err)
   }
   del, err := auth.Delivery(dev, apps[0].id, apps[0].ver)
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
   file, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      return nil, "", err
   }
   defer file.Close()
   if err := tok.Decode(file); err != nil {
      return nil, "", err
   }
   auth, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return auth, cache, nil
}

func TestTokenEncode(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   cache += "/googleplay"
   os.Mkdir(cache, os.ModeDir)
   file, err := os.Create(cache + "/token.json")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   if err := tok.Encode(file); err != nil {
      t.Fatal(err)
   }
}
