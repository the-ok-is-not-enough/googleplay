package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

var appsX86 = []app{
   {"Apr 11, 2022", "com.instagram.android"},
   {"Apr 11, 2022", "com.pinterest"},
   {"Apr 12, 2022", "br.com.rodrigokolb.realdrum"},
   {"Apr 4, 2022", "com.vimeo.android.videoapp"},
   {"Apr 6, 2022", "org.thoughtcrime.securesms"},
   {"Apr 7, 2022", "com.google.android.youtube"},
   {"Feb 14, 2022", "org.videolan.vlc"},
   {"Jun 1, 2021", "com.valvesoftware.android.steam.community"},
   {"Mar 1, 2022", "kr.sira.metal"},
   {"Mar 17, 2022", "com.google.android.apps.walletnfcrel"},
   {"Mar 21, 2022", "com.jackpocket"},
   {"May 22, 2021", "com.smarty.voomvoom"},
}

var appsArmeabi = []app{
   {"Apr 8, 2022", "com.axis.drawingdesk.v3"},
   {"Mar 14, 2022", "com.xiaomi.smarthome"},
   {"Mar 24, 2022", "com.miui.weather2"},
}

var appsArm64 = []app{
   {"Mar 30, 2022", "com.exnoa.misttraingirls"},
}

func TestDetailsArm64(t *testing.T) {
   err := testDetails("googleplay/arm64.json", appsArm64)
   if err != nil {
      t.Fatal(err)
   }
}

func TestDetailsArmeabi(t *testing.T) {
   err := testDetails("googleplay/armeabi.json", appsArmeabi)
   if err != nil {
      t.Fatal(err)
   }
}

func TestDetailsX86(t *testing.T) {
   err := testDetails("googleplay/x86.json", appsX86)
   if err != nil {
      t.Fatal(err)
   }
}

type app struct {
   date string
   id string
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
