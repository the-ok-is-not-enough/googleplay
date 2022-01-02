package play

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

const auth = "ya29.A0ARrdaM8_M3iTgk8wQIA_h-ATs6C9Tmgq8J-GaT9TY6ZgyFwRiB4R6tmtV5iKBvd6K41jsrD-v_jwBvex2GoSByixifIe5UQ4C3ev1B7JT-xcjBJTy9G6IwULJv4rfV64h-hwbcVbQsdR0JKgg4kVtUy3zZkO4nMxf8z5bHDdYwH4iYyX3cp8WZ-uvbjPMsjNlapomWSA4dA0xxb-HK1WKebeuDdjjrPQ8FuoRd2aOu7QTyd1H_BYoiNA6qOz_7yj0ETtp2hkVowEoGRJjlAsHMs5X3mfyn-7unjUC84WkfOOJkUkdkw"

var body0 = strings.NewReader("\n\xf1\x04\b\x03\x18\x00(\x000\x008\x00J\x12com.samsung.deviceJ\x17global-miui11-empty.jarz#GL_OES_compressed_ETC1_RGB8_texture\x10\x00 \x00@\x81\x80\fR\x1aandroid.hardware.bluetoothR\x1dandroid.hardware.bluetooth_leR\x17android.hardware.cameraR!android.hardware.camera.autofocusR\x19android.hardware.locationR\x1dandroid.hardware.location.gpsR\x1bandroid.hardware.microphoneR!android.hardware.screen.landscapeR android.hardware.screen.portraitR%android.hardware.sensor.accelerometerR\x1aandroid.hardware.telephonyR\x1candroid.hardware.touchscreenR\x19android.hardware.usb.hostR\x15android.hardware.wifiR$com.samsung.android.api.version.2601R-com.samsung.feature.samsung_experience_mobileZ\x03x86Z\varmeabi-v7a")

func upload(id string) error {
   var req0 = &http.Request{Method:"POST", URL:&url.URL{Scheme:"https",
   Opaque:"", User:(*url.Userinfo)(nil), Host:"android.clients.google.com",
   Path:"/fdfe/uploadDeviceConfig", RawPath:"", ForceQuery:false, RawQuery:"",
   Fragment:"", RawFragment:""},
   Header:http.Header{"Authorization":[]string{"Bearer " + auth},
   "Host":[]string{"android.clients.google.com"},
   "User-Agent":[]string{"Android-Finsky (sdk=99,versionCode=99999999)"},
   "X-Dfe-Device-Id":[]string{id},
   }, Body:io.NopCloser(body0)}
   res, err := new(http.Transport).RoundTrip(req0)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   time.Sleep(16 * time.Second)
   return nil
}

var body1 = protobuf.Message{
   protobuf.Tag{Number:4}:protobuf.Message{
      protobuf.Tag{Number:1}:protobuf.Message{
         protobuf.Tag{Number:10}:uint64(29),
      },
   },
   protobuf.Tag{Number:14}:uint64(3),
   protobuf.Tag{Number:18}:protobuf.Message{
      protobuf.Tag{Number:1, String:"touchScreen"}:uint64(0),
      protobuf.Tag{Number:2, String:"keyboard"}:uint64(0),
      protobuf.Tag{Number:3, String:"navigation"}:uint64(0),
      protobuf.Tag{Number:4, String:"screenLayout"}:uint64(0),
      protobuf.Tag{Number:5, String:"hasHardKeyboard"}:uint64(0),
      protobuf.Tag{Number:6, String:"hasFiveWayNavigation"}:uint64(0),
      protobuf.Tag{Number:7, String:"screenDensity"}:uint64(0),
      protobuf.Tag{Number:8, String:"glEsVersion"}:uint64(0x3_0000),
      protobuf.Tag{Number:12, String:"screenWidth"}:uint64(1),
      protobuf.Tag{Number:26}:[]protobuf.Message{
         protobuf.Message{
            protobuf.Tag{Number:1}:"android.hardware.touchscreen",
         },
         protobuf.Message{
            protobuf.Tag{Number:1}:"android.hardware.screen.portrait",
         },
         protobuf.Message{
            protobuf.Tag{Number:1}:"android.hardware.wifi",
         },
      },
   },
}

func checkin() (uint64, error) {
   var req0 = &http.Request{
      Method:"POST",
      URL:&url.URL{Scheme:"https",
         Host:"android.clients.google.com",
         Path:"/checkin", 
      },
      Header:http.Header{
         "Content-Type":[]string{"application/x-protobuffer"},
      },
      Body: io.NopCloser(body1.Encode()),
   }
   res, err := new(http.Transport).RoundTrip(req0)
   if err != nil {
      return 0, err
   }
   defer res.Body.Close()
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(7), nil
}

func details(app string) (uint64, error) {
   id, err := checkin()
   if err != nil {
      return 0, err
   }
   sID := strconv.FormatUint(id, 16)
   fmt.Println(sID)
   var req5 = &http.Request{Method:"GET", URL:&url.URL{Scheme:"https",
      Host:"android.clients.google.com",
      Path:"/fdfe/details", RawQuery:"doc=" + app,
      },
      Header:http.Header{
         "Authorization":[]string{"Bearer " + auth},
         "X-Dfe-Device-Id":[]string{sID},
      },
   }
   res, err := new(http.Transport).RoundTrip(req5)
   if err != nil {
      return 0, err
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(1,2,4,13,1,3), nil
}
