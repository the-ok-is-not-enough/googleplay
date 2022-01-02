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

const auth = "ya29.a0ARrdaM9WB7UHQMrsWAaBOm4-ndsaI6_qxJ0olxmNaWrCSlnr-1_uUg2HYeqXAZ8our4ZB_kFrS7zrkCahV6oE50V7kcqv6HXS0us1V16fRaEmAB__z3mGHRikES6g0yoi0fJZ1XbiPd7kAgWGQxkjZAICFj7dOJ4l-PSn2PbSbRN00EynjQ1FHqcViuRskgvKRwG2NZ7F9YC8LorYRSmdvogykHQSmWaLIj1DMMp-Ida6rsb4Ce2z1ob7Rgu8XZmwERWlTen7vKvCpqE9_2lQa216UJqslQT7ImGwVpEnEWcYOMKDYc"

var body1 = protobuf.Message{
   protobuf.Tag{Number:1, String:""}:protobuf.Message{
      protobuf.Tag{Number:1, String:""}:uint64(3),
      protobuf.Tag{Number:2, String:""}:uint64(0),
      protobuf.Tag{Number:3, String:""}:uint64(0),
      protobuf.Tag{Number:4, String:""}:uint64(0),
      protobuf.Tag{Number:5, String:""}:uint64(0),
      protobuf.Tag{Number:6, String:""}:uint64(0),
      protobuf.Tag{Number:7, String:""}:uint64(0),
      protobuf.Tag{Number:8, String:""}:uint64(196609),
      protobuf.Tag{Number:9, String:""}:[]string{
         "com.samsung.device", "global-miui11-empty.jar"
      },
      protobuf.Tag{Number:10, String:""}:[]string{
         "android.hardware.bluetooth",
         "android.hardware.bluetooth_le", "android.hardware.camera",
         "android.hardware.camera.autofocus", "android.hardware.location",
         "android.hardware.location.gps", "android.hardware.microphone",
         "android.hardware.screen.landscape", "android.hardware.screen.portrait",
         "android.hardware.sensor.accelerometer", "android.hardware.telephony",
         "android.hardware.touchscreen", "android.hardware.usb.host",
         "android.hardware.wifi", "com.samsung.android.api.version.2601",
         "com.samsung.feature.samsung_experience_mobile",
      },
      protobuf.Tag{Number:11, String:""}:[]string{"x86", "armeabi-v7a"},
      protobuf.Tag{Number:15, String:""}:"GL_OES_compressed_ETC1_RGB8_texture",
   },
}

func upload(id string) error {
   var req0 = &http.Request{
      Method:"POST",
      URL:&url.URL{
         Scheme:"https",
         Host:"android.clients.google.com",
         Path:"/fdfe/uploadDeviceConfig",
      },
      Header:http.Header{
         "Authorization":[]string{"Bearer " + auth},
         "Host":[]string{"android.clients.google.com"},
         "User-Agent":[]string{"Android-Finsky (sdk=99,versionCode=99999999)"},
         "X-Dfe-Device-Id":[]string{id},
      },
      Body:io.NopCloser(body1.Encode()),
   }
   res, err := new(http.Transport).RoundTrip(req0)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   time.Sleep(16 * time.Second)
   return nil
}

func checkin() (int64, error) {
   req0 := &http.Request{
      Method:"POST",
      URL:&url.URL{
         Scheme:"https",
         Host:"android.clients.google.com", Path:"/checkin",
      },
      Header:http.Header{"Host":[]string{"android.clients.google.com"}},
      Body:io.NopCloser(strings.NewReader("{\"checkin\":{},\"version\":3}\n")),
   }
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
