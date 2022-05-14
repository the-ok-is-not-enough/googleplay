package googleplay

import (
   "testing"
   "time"
)

var testApps = map[string][]app{
   "googleplay/arm64.json": {
      {"May 9, 2022", "com.kakaogames.twodin"},
   },
   "googleplay/armeabi.json": {
      {"Mar 14, 2022", "com.xiaomi.smarthome"},
      {"Mar 24, 2022", "com.miui.weather2"},
      {"Apr 8, 2022", "com.axis.drawingdesk.v3"},
   },
   "googleplay/x86.json": {
      {"Jun 1, 2021", "com.valvesoftware.android.steam.community"},
      {"Dec 8, 2021", "com.amctve.amcfullepisodes"},
      {"Feb 14, 2022", "org.videolan.vlc"},
      {"Mar 1, 2022", "kr.sira.metal"},
      {"Mar 17, 2022", "com.google.android.apps.walletnfcrel"},
      {"Apr 1, 2022", "com.clearchannel.iheartradio.controller"},
      {"Apr 6, 2022", "org.thoughtcrime.securesms"},
      {"Apr 7, 2022", "com.google.android.youtube"},
      {"Apr 11, 2022", "com.pinterest"},
      {"Apr 12, 2022", "br.com.rodrigokolb.realdrum"},
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
