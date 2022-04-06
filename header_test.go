package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

var tabletApps = []app{
   {id: "com.google.android.apps.youtube.music.pwa"},
}

var tvApps = []app{
   {down: "148.435 M", id: "com.google.android.youtube.tv"},
   {down: "3.934 M", id: "com.google.android.youtube.googletv"},
}

var phoneApps = []app{
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
   {down: "10.683 M", id: "com.google.android.apps.youtube.vr"},
   {down: "9.419 M", id: "com.axis.drawingdesk.v3"},
   {down: "282.669 K", id: "com.smarty.voomvoom"},
   {down: "83.801 K", id: "com.exnoa.misttraingirls"},
   {down: "58.860 K", id: "se.pax.calima"},
}

func TestTabletDetails(t *testing.T) {
   err := testDetails("googleplay/tablet.json", tabletApps)
   if err != nil {
      t.Fatal(err)
   }
}

func TestTvDetails(t *testing.T) {
   err := testDetails("googleplay/tv.json", tvApps)
   if err != nil {
      t.Fatal(err)
   }
}

func TestPhoneDetails(t *testing.T) {
   err := testDetails("googleplay/phone.json", phoneApps)
   if err != nil {
      t.Fatal(err)
   }
}

func (a app) Error() string {
   return a.id
}

type app struct {
   down string
   id string
   ver uint64
}

func testDetails(device string, apps []app) error {
   head, err := newHeader(device)
   if err != nil {
      return err
   }
   for _, app := range apps {
      det, err := head.Details(app.id)
      if err != nil {
         return err
      }
      if det.CurrencyCode == "" {
         return app
      }
      if det.NumDownloads == 0 {
         return app
      }
      if det.Size == 0 {
         return app
      }
      if det.Title == "" {
         return app
      }
      if det.UploadDate == "" {
         return app
      }
      if det.VersionCode == 0 {
         return app
      }
      if det.VersionString == "" {
         return app
      }
      time.Sleep(99 * time.Millisecond)
   }
   return nil
}

func TestDelivery(t *testing.T) {
   head, err := newHeader("googleplay/phone.json")
   if err != nil {
      t.Fatal(err)
   }
   del, err := head.Delivery(phoneApps[0].id, phoneApps[0].ver)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", del)
}

func newHeader(device string) (*Header, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, err
   }
   tok, err := OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   phone, err := OpenDevice(cache, device)
   if err != nil {
      return nil, err
   }
   return tok.Header(phone)
}
