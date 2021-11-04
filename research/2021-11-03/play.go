package googleplay

import (
   "bytes"
   "encoding/json"
   "fmt"
   //"google.golang.org/protobuf/proto"
   //"github.com/segmentio/encoding/proto"
   "net/http"
   "net/http/httputil"
   "strconv"
   "strings"
)

const origin = "https://android.clients.google.com"

type uploadDeviceConfigRequest struct {
   DeviceConfiguration struct {
      GlEsVersion            int32   `protobuf:"varint,8"`
      HasFiveWayNavigation   *bool    `protobuf:"varint,6"`
      HasHardKeyboard        *bool    `protobuf:"varint,5"`
      Keyboard               int32   `protobuf:"varint,2"`
      NativePlatform         []string `protobuf:"bytes,11"`
      Navigation             int32   `protobuf:"varint,3"`
      ScreenDensity          int32   `protobuf:"varint,7"`
      ScreenLayout           int32   `protobuf:"varint,4"`
      SystemAvailableFeature []string `protobuf:"bytes,10"`
      TouchScreen            int32   `protobuf:"varint,1"`
   } `protobuf:"bytes,1"`
}

var buf = []byte("\ni\x08\x03\x10\x01\x18\x01 \x02(\x000\x008\xa4\x03@\x82\x80\x0cR\x1candroid.hardware.touchscreenR\x15android.hardware.wifiZ\tarm64-v8aZ\x0barmeabi-v7aZ\x07armeabi")

func upload(auth, deviceID string) error {
   /*
   var (
      version int32 = 196610
      fivewayNavigation bool
      hardKeyboard bool
      keyboard int32 = 1
      navigation int32 = 1
      touchscreen int32 = 3
      screenLayout int32 = 2
      screenDensity int32 = 420
   )
   upload := &UploadDeviceConfigRequest{
      DeviceConfiguration: &DeviceConfigurationProto{
         GlEsVersion            : &version,
         HasFiveWayNavigation   : &fivewayNavigation,
         HasHardKeyboard        : &hardKeyboard,
         Keyboard               : &keyboard,
         Navigation             : &navigation,
         NativePlatform         : []string{"arm64-v8a,armeabi-v7a", "armeabi"},
         TouchScreen            : &touchscreen,
         ScreenLayout           : &screenLayout,
         ScreenDensity          : &screenDensity,
         SystemAvailableFeature : []string{
            "android.hardware.touchscreen", "android.hardware.wifi",
         },
      },
   }
   buf, err := proto.Marshal(upload)
   var fal bool
   var u uploadDeviceConfigRequest
   u.DeviceConfiguration.HasFiveWayNavigation = &fal
   u.DeviceConfiguration.HasHardKeyboard = &fal
   u.DeviceConfiguration.GlEsVersion=196610
   u.DeviceConfiguration.Keyboard = 1
   u.DeviceConfiguration.Navigation = 1
   u.DeviceConfiguration.ScreenDensity = 420
   u.DeviceConfiguration.ScreenLayout = 2
   u.DeviceConfiguration.TouchScreen = 3
   u.DeviceConfiguration.NativePlatform = []string{
      "arm64-v8a,armeabi-v7a", "armeabi",
   }
   u.DeviceConfiguration.SystemAvailableFeature = []string{
      "android.hardware.touchscreen","android.hardware.wifi",
   }
   buf, err := proto.Marshal(u)
   if err != nil {
      return err
   }
   */
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", bytes.NewReader(buf),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + auth},
      "Content-Type": {"application/x-protobuf"},
      "User-Agent": {"Android-Finsky (versionCode=81031200,sdk=27)"},
      "X-DFE-Device-ID": {deviceID},
   }
   dum, err := httputil.DumpRequest(req, false)
   if err != nil {
      return err
   }
   fmt.Printf("%s\n", dum)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   dum, err = httputil.DumpResponse(res, true)
   if err != nil {
      return err
   }
   fmt.Printf("%q\n", dum)
   if res.StatusCode != http.StatusOK {
      return fmt.Errorf("status %q", res.Status)
   }
   return nil
}

type checkin struct {
   Android_ID int64
}

func newCheckin() (*checkin, error) {
   req, err := http.NewRequest(
      "POST", origin + "/checkin",
      strings.NewReader(`{"checkin":{}, "version":3}`),
   )
   if err != nil {
      return nil, err
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   dum, err := httputil.DumpResponse(res, true)
   if err != nil {
      return nil, err
   }
   fmt.Printf("%s\n", dum)
   check := new(checkin)
   if err := json.NewDecoder(res.Body).Decode(check); err != nil {
      return nil, err
   }
   return check, nil
}

func (c checkin) String() string {
   return strconv.FormatInt(c.Android_ID, 16)
}
