package main

import (
   "encoding/json"
   "fmt"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/http/httputil"
   "os"
   "net/url"
   "strconv"
   "strings"
   "time"
)

const auth = "ya29.a0ARrdaM_u42HehDpxYQih_9NlnYwP7Zl6xK6MYVwXzThh91zi4FRmSjx-6KBdmmQPnSVipoUE5Rl8-jEJCfCujaIVAyhPkA8ZrAebGD-uV7OQRdNyH3GTS8ZlnUx0T_Q8xMHTmTGXUrjz5Hzhdbq327RkyI_EJr5_QwYsC8TXOt1Yggh_VEvXFV95lu8PjQtQIEkktNQtY7yCDX5CDaVz_C1SKNn7QrW0deI5p4eZ6Rr8qj4W6_pUg3bUcPkfoTwcrKYf_ICaamOz7hlR5g7gRVYGUZvPsUFJU-1imB1VNaYBInjV9gg"

var body0 = strings.NewReader("\n\xf1\x04\b\x03\x18\x00(\x000\x008\x00J\x12com.samsung.deviceJ\x17global-miui11-empty.jarz#GL_OES_compressed_ETC1_RGB8_texture\x10\x00 \x00@\x81\x80\fR\x1aandroid.hardware.bluetoothR\x1dandroid.hardware.bluetooth_leR\x17android.hardware.cameraR!android.hardware.camera.autofocusR\x19android.hardware.locationR\x1dandroid.hardware.location.gpsR\x1bandroid.hardware.microphoneR!android.hardware.screen.landscapeR android.hardware.screen.portraitR%android.hardware.sensor.accelerometerR\x1aandroid.hardware.telephonyR\x1candroid.hardware.touchscreenR\x19android.hardware.usb.hostR\x15android.hardware.wifiR$com.samsung.android.api.version.2601R-com.samsung.feature.samsung_experience_mobileZ\x03x86Z\varmeabi-v7a")

func checkin() (int64, error) {
   req0 := &http.Request{Method:"POST", URL:&url.URL{Scheme:"https", Opaque:"",
   User:(*url.Userinfo)(nil), Host:"android.clients.google.com", Path:"/checkin",
   RawPath:"", ForceQuery:false, RawQuery:"", Fragment:"", RawFragment:""},
   Header:http.Header{"Host":[]string{"android.clients.google.com"}},
   Body:io.NopCloser(strings.NewReader("{\"checkin\":{},\"version\":3}\n"))}
   res, err := new(http.Transport).RoundTrip(req0)
   if err != nil {
      return 0, err
   }
   defer res.Body.Close()
   var check struct {
      Android_ID int64
   }
   if err := json.NewDecoder(res.Body).Decode(&check); err != nil {
      return 0, err
   }
   return check.Android_ID, nil
}

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

func main() {
   id, err := checkin()
   if err != nil {
      panic(err)
   }
   sID := strconv.FormatInt(id, 16)
   if err := upload(sID); err != nil {
      panic(err)
   }
   req5 := &http.Request{Method:"GET", URL:&url.URL{Scheme:"https", Opaque:"",
User:(*url.Userinfo)(nil), Host:"android.clients.google.com",
Path:"/fdfe/details", RawPath:"", ForceQuery:false,
RawQuery:"doc=com.google.android.youtube", Fragment:"", RawFragment:""},
Header:http.Header{"Authorization":[]string{"Bearer " + auth},
"Host":[]string{"android.clients.google.com"},
"X-Dfe-Device-Id":[]string{sID}}}
   buf, err := httputil.DumpRequestOut(req5, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
   res, err := new(http.Transport).RoundTrip(req5)
   if err != nil {
      panic(err)
   }
   if res.StatusCode != http.StatusOK {
      panic(res.Status)
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      panic(err)
   }
   ver := mes.GetUint64(1,2,4,13,1,3)
   fmt.Println(ver)
}
