package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

type app struct {
   id string
   code int
}

var apps = []app{
   0: {"com.google.android.youtube", 1524221376},
   1: {"com.axis.drawingdesk.v3", 190},
   2: {"com.instagram.android", 321403734},
   3: {"com.pinterest", 9398020},
   4: {"com.smarty.voomvoom", 369},
   5: {"com.vimeo.android.videoapp", 3510004},
   6: {"org.videolan.vlc", 13040205},
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
   del, err := auth.DeliveryResponse(dev, apps[0].id, apps[0].code)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", del)
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
   det, err := auth.DetailsResponse(dev, apps[0].id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", det)
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
   r, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      t.Fatal(err)
   }
   defer r.Close()
   if err := tok.Decode(r); err != nil {
      t.Fatal(err)
   }
   auth, err := tok.Auth()
   if err != nil {
      t.Fatal(err)
   }
   if err := auth.Upload(dev, NewConfig()); err != nil {
      t.Fatal(err)
   }
   w, err := os.Create(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   defer w.Close()
   if err := dev.Encode(w); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}

const email = "srpen6@gmail.com"

func TestTokenEncode(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   c, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   c += "/googleplay"
   os.Mkdir(c, os.ModeDir)
   f, err := os.Create(c + "/token.json")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   if err := tok.Encode(f); err != nil {
      t.Fatal(err)
   }
}
