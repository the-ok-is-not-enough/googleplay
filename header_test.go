package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

type app struct {
   id string
   nativeCode []string
}

var phoneApps = []app{
   {"br.com.rodrigokolb.realdrum", []string{
      "arm64-v8a", "armeabi-v7a", "x86", "x86_64",
   }},
   {"com.google.android.apps.walletnfcrel", nil},
   {"com.google.android.youtube", nil},
   {"com.instagram.android", []string{
      "x86",
   }},
   {"com.miui.weather2", []string{
      "arm64-v8a", "armeabi", "armeabi-v7a",
   }},
   {"com.pinterest", nil},
   {"com.valvesoftware.android.steam.community", nil},
   {"com.xiaomi.smarthome", []string{
      "arm64-v8a",
   }},
   {"org.thoughtcrime.securesms", []string{
      "x86",
   }},
   {"org.videolan.vlc", []string{
      "arm64-v8a",
   }},
   /////////////////////////////////////////////////////////////////////////////
   {"com.vimeo.android.videoapp", nil},
   {"kr.sira.metal", nil},
   {"com.tgc.sky.android", nil},
   {"com.google.android.apps.youtube.vr", nil},
   {"com.axis.drawingdesk.v3", nil},
   {"com.hamsterbeat.wallpapers.fx.panorama", nil},
   {"com.jackpocket", nil},
   {"com.smarty.voomvoom", nil},
   {"com.exnoa.misttraingirls", nil},
   {"se.pax.calima", nil},
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
   del, err := head.Delivery("com.google.android.youtube", 1524221376)
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
