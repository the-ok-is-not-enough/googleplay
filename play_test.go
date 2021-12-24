package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

const email = "srpen6@gmail.com"

type app struct {
   id string
   code int
}

var apps = map[int]app{
   0: {"com.google.android.youtube", 1524221376},
   1: {id: "com.axis.drawingdesk.v3"},
   2: {id: "com.instagram.android"},
   3: {id: "com.pinterest"},
   4: {id: "com.smarty.voomvoom"},
   5: {id: "com.vimeo.android.videoapp"},
   6: {id: "org.videolan.vlc"},
   7: {id: "org.thoughtcrime.securesms"},
   8: {id: "com.valvesoftware.android.steam.community"},
   9: {id: "com.miui.weather2"},
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
      time.Sleep(time.Second)
   }
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
   del, err := auth.Delivery(dev, apps[0].id, apps[0].code)
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

func TestDevice(t *testing.T) {
   dev, err := NewDevice()
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
   auth, err := tok.Auth()
   if err != nil {
      t.Fatal(err)
   }
   if err := auth.Upload(dev, DefaultConfig); err != nil {
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
