package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

func TestCategory(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev, err := OpenDevice(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   head := auth.Header(dev)
   head.language("ko")
   docs, err := head.Category("GAME")
   if err != nil {
      t.Fatal(err)
   }
   for _, doc := range docs {
      fmt.Print(doc, "\n---\n")
   }
}

func TestReview(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev, err := OpenDevice(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   LogLevel = 1
   revs, err := auth.Header(dev).Reviews("com.comuto")
   if err != nil {
      t.Fatal(err)
   }
   for _, rev := range revs {
      fmt.Printf("%+v\n", rev)
   }
}

var apps = []app{
   {down: "10.996 B", id: "com.google.android.youtube", ver: 1524221376},
   {down: "3.932 B", id: "com.instagram.android"},
   {down: "975.149 M", id: "com.miui.weather2"},
   {down: "689.574 M", id: "com.pinterest"},
   {down: "422.289 M", id: "com.google.android.apps.walletnfcrel"},
   {down: "282.147 M", id: "org.videolan.vlc"},
   {down: "95.910 M", id: "org.thoughtcrime.securesms"},
   {down: "77.289 M", id: "com.valvesoftware.android.steam.community"},
   {down: "31.446 M", id: "com.xiaomi.smarthome"},
   {down: "30.702 M", id: "com.vimeo.android.videoapp"},
   {down: "13.832 M", id: "com.tgc.sky.android"},
   {down: "9.419 M", id: "com.axis.drawingdesk.v3"},
   {down: "371.263 K", id: "com.bbca.bbcafullepisodes"},
   {down: "282.669 K", id: "com.smarty.voomvoom"},
   {down: "83.801 K", id: "com.exnoa.misttraingirls"},
   {down: "58.860 K", id: "se.pax.calima"},
}

func TestToken(t *testing.T) {
   tok, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   if err := tok.Create(cache + "/googleplay/token.json"); err != nil {
      t.Fatal(err)
   }
}

const email = "srpen6@gmail.com"

type app struct {
   down, id string
   ver int64
}

func TestDetails(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev, err := OpenDevice(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   for _, app := range apps {
      det, err := auth.Header(dev).Details(app.id)
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
   dev, err := NewDevice(DefaultConfig)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   if err := dev.Create(cache + "/googleplay/device.json"); err != nil {
      t.Fatal(err)
   }
   time.Sleep(Sleep)
}

func TestDelivery(t *testing.T) {
   auth, cache, err := getAuth()
   if err != nil {
      t.Fatal(err)
   }
   dev, err := OpenDevice(cache + "/googleplay/device.json")
   if err != nil {
      t.Fatal(err)
   }
   del, err := auth.Header(dev).Delivery(apps[0].id, apps[0].ver)
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
   tok, err := OpenToken(cache + "/googleplay/token.json")
   if err != nil {
      return nil, "", err
   }
   auth, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return auth, cache, nil
}
