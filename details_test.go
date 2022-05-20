package googleplay

import (
   "testing"
   "time"
)

var testApps = map[string][]app{
   "googleplay/x86.json": {
      {"2021-06-01", "com.valvesoftware.android.steam.community"},
      {"2021-12-08", "com.amctve.amcfullepisodes"},
      {"2022-02-14", "org.videolan.vlc"},
      {"2022-03-01", "kr.sira.metal"},
      {"2022-03-17", "com.google.android.apps.walletnfcrel"},
      {"2022-03-24", "app.source.getcontact"},
      {"2022-05-05", "com.clearchannel.iheartradio.controller"},
      {"2022-05-12", "com.google.android.youtube"},
      {"2022-05-13", "com.app.xt"},
      {"2022-05-16", "com.binance.dev"},
      {"2022-05-16", "com.instagram.android"},
      {"2022-05-17", "br.com.rodrigokolb.realdrum"},
      {"2022-05-17", "com.pinterest"},
      {"2022-05-17", "org.thoughtcrime.securesms"},
   },
   "googleplay/armeabi-v7a.json": {
      {"2022-03-24", "com.miui.weather2"},
      {"2022-04-08", "com.axis.drawingdesk.v3"},
      {"2022-04-27", "com.sygic.aura"},
      {"2022-04-29", "com.xiaomi.smarthome"},
      {"2022-05-18", "com.madhead.tos.zh"},
   },
   "googleplay/arm64-v8a.json": {
      {"2022-02-02", "com.illumix.fnafar"},
      {"2022-03-16", "com.miHoYo.GenshinImpact"},
      {"2022-05-15", "com.kakaogames.twodin"},
   },
}

func TestDetails(t *testing.T) {
   for elem, apps := range testApps {
      head, err := newHeader(elem)
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
         time.Sleep(99 * time.Millisecond)
      }
   }
}

type app struct {
   date string
   id string
}
