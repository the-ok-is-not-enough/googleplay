package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

func TestDetails(t *testing.T) {
   head, err := newHeader()
   if err != nil {
      t.Fatal(err)
   }
   for _, app := range apps {
      det, err := head.Details(app.id)
      if err != nil {
         t.Fatal(err)
      }
      if det.CurrencyCode == "" {
         t.Fatal(det)
      }
      if det.NumDownloads == 0 {
         t.Fatal(det)
      }
      if det.Size == 0 {
         t.Fatal(det)
      }
      if det.Title == "" {
         t.Fatal(det)
      }
      if det.UploadDate == "" {
         t.Fatal(det)
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

type app struct {
   down, id string
   ver uint64
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
   {down: "282.669 K", id: "com.smarty.voomvoom"},
   {down: "83.801 K", id: "com.exnoa.misttraingirls"},
   {down: "58.860 K", id: "se.pax.calima"},
   {down: "1", id: "com.google.android.GoogleCamera"},
}

func TestDelivery(t *testing.T) {
   head, err := newHeader()
   if err != nil {
      t.Fatal(err)
   }
   del, err := head.Delivery(apps[0].id, apps[0].ver)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", del)
}

func newHeader() (*Header, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, err
   }
   tok, err := OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   dev, err := OpenDevice(cache, "googleplay/device.json")
   if err != nil {
      return nil, err
   }
   return tok.Header(dev)
}
