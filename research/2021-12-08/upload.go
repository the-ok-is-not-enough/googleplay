package main

import (
   "encoding/json"
   "fmt"
   "github.com/89z/parse/protobuf"
   "io"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
   "strings"
   "time"
)

var config = protobuf.Message{
   protobuf.Tag{Number:1, String:""}:protobuf.Message{
      protobuf.Tag{Number:1, String:""}:uint64(0),
      protobuf.Tag{Number:2, String:""}:uint64(0),
      protobuf.Tag{Number:3, String:""}:uint64(0),
      protobuf.Tag{Number:4, String:""}:uint64(0),
      protobuf.Tag{Number:5, String:""}:uint64(0),
      protobuf.Tag{Number:6, String:""}:uint64(0),
      protobuf.Tag{Number:7, String:""}:uint64(0),
      protobuf.Tag{Number:15, String:""}:[]string{
         "GL_OES_compressed_ETC1_RGB8_texture",
      },
      protobuf.Tag{Number:8, String:""}:uint64(0x3_0001),
      protobuf.Tag{Number:11, String:""}:[]string{
         "x86",
         "armeabi-v7a",
      },
      protobuf.Tag{Number:10, String:""}:[]string{
         "android.hardware.bluetooth",
         "android.hardware.camera",
         "android.hardware.location",
         "android.hardware.location.gps",
         "android.hardware.screen.landscape",
         "android.hardware.screen.portrait",
         "android.hardware.sensor.accelerometer",
         "android.hardware.touchscreen",
         "android.hardware.wifi",
         "android.hardware.microphone",
         // FIXME
         //"android.software.leanback",
      },
   },
}

func upload(dev *Device) error {
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", config.Encode(),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {auth},
      "User-Agent": {"Android-Finsky (sdk=99,versionCode=99999999)"},
      "X-DFE-Device-ID": {dev.String()},
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return status(res.Status)
   }
   return nil
}

func main() {
   dev, err := NewDevice()
   if err != nil {
      panic(err)
   }
   if err := upload(dev); err != nil {
      panic(err)
   }
   fmt.Println("Sleep")
   time.Sleep(16 * time.Second)
   req, err := http.NewRequest(
      "GET",
      origin + "/fdfe/delivery?doc=com.google.android.youtube&vc=1524493760",
      nil,
   )
   if err != nil {
      panic(err)
   }
   req.Header = http.Header{
      "Authorization": {auth},
      "User-Agent": {"Android-Finsky (sdk=99,versionCode=99999999)"},
      "X-DFE-Device-ID": {dev.String()},
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err := httputil.DumpResponse(res, false)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
   buf, err = io.ReadAll(res.Body)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%q\n", buf)
}

const origin = "https://android.clients.google.com"

type Device struct {
   Android_ID int64
}

func NewDevice() (*Device, error) {
   req, err := http.NewRequest(
      "POST", origin + "/checkin",
      strings.NewReader(`{"checkin":{},"version":3}`),
   )
   if err != nil {
      return nil, err
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   dev := new(Device)
   if err := json.NewDecoder(res.Body).Decode(dev); err != nil {
      return nil, err
   }
   return dev, nil
}

func (d Device) String() string {
   return strconv.FormatInt(d.Android_ID, 16)
}

type status string

func (s status) Error() string {
   return string(s)
}


