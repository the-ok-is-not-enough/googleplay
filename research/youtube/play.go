package play

import (
   "encoding/json"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

const auth = "ya29.a0ARrdaM8WY37NU2AwiJuhmKBCVee_VF0wBYHL1FtWCpZ-_2y-dJgyaEij0GW7Cnvh-NM0GNWpgLFxrFD97wjfSyrkw_f7uiyHnveUKjEJtVQrNobuqXayGZFJWDrXgFIjZdPRf9KXSTWqXtLEDgyMoFsztF1Jw_wma4onRymLnK6fyaBbCqyAFHh8Jo73nVu5ovVOzDQSOR1-dN8A4qdYpcfE7J9drzSUj5QMGmKeh9MhR9qM9EhDgf5ZTn-adEursZYymjRpRxVC8a2EQbbbV6O-IhNav37ptVgK3Ko-yBWxW44UDFY"

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

func details(app string) (uint64, error) {
   id, err := checkin()
   if err != nil {
      return 0, err
   }
   sID := strconv.FormatInt(id, 16)
   if err := upload(sID); err != nil {
      return 0, err
   }
   req5 := &http.Request{
      Method:"GET",
      URL:&url.URL{
         Scheme:"https", Opaque:"",
         User:(*url.Userinfo)(nil), Host:"android.clients.google.com",
         Path:"/fdfe/details", RawPath:"", ForceQuery:false,
         RawQuery:"doc=" + app,
      },
      Header:http.Header{"Authorization":[]string{"Bearer " + auth},
      "Host":[]string{"android.clients.google.com"},
      "X-Dfe-Device-Id":[]string{sID}},
   }
   res, err := new(http.Transport).RoundTrip(req5)
   if err != nil {
      return 0, err
   }
   if res.StatusCode != http.StatusOK {
      return 0, response{res}
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(1,2,4,13,1,3), nil
}

type response struct {
   *http.Response
}

func (r response) Error() string {
   return r.Status
}
