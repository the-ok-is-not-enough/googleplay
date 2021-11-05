package googleplay

import (
   "bytes"
   "encoding/json"
   "fmt"
   "github.com/segmentio/encoding/proto"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

type Checkin struct {
   Checkin struct{} `json:"checkin"`
   Version int `json:"version"`
}

func NewCheckin() Checkin {
   return Checkin{Version: 3}
}

func (c Checkin) Post() (*CheckinResponse, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(c)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", origin + "/checkin", buf)
   if err != nil {
      return nil, err
   }
   dum, err := httputil.DumpRequest(req, true)
   if err != nil {
      return nil, err
   }
   os.Stdout.Write(dum)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   check := new(CheckinResponse)
   if err := json.NewDecoder(res.Body).Decode(check); err != nil {
      return nil, err
   }
   return check, nil
}

type CheckinResponse struct {
   Android_ID int64
}

func (c CheckinResponse) String() string {
   return strconv.FormatInt(c.Android_ID, 16)
}

type Device struct {
   Configuration struct {
      GlEsVersion            int32   `protobuf:"varint,8"`
      HasFiveWayNavigation   bool    `protobuf:"varint,6"`
      HasHardKeyboard        bool    `protobuf:"varint,5"`
      Keyboard               int32   `protobuf:"varint,2"`
      NativePlatform         []string `protobuf:"bytes,11"`
      Navigation             int32   `protobuf:"varint,3"`
      ScreenDensity          int32   `protobuf:"varint,7"`
      ScreenLayout           int32   `protobuf:"varint,4"`
      SystemAvailableFeature []string `protobuf:"bytes,10"`
      TouchScreen            int32   `protobuf:"varint,1"`
   } `protobuf:"bytes,1"`
}

func NewDevice() Device {
   var d Device
   // developer.android.com/guide/topics/manifest/uses-feature-element
   d.Configuration.GlEsVersion = 0x0009_0000
   d.Configuration.HasFiveWayNavigation = true
   d.Configuration.HasHardKeyboard = true
   d.Configuration.Keyboard = 1
   // developer.android.com/ndk/guides/abis
   d.Configuration.NativePlatform = []string{
      "armeabi-v7a",
   }
   d.Configuration.Navigation = 1
   d.Configuration.ScreenDensity = 1
   d.Configuration.ScreenLayout = 1
   // developer.android.com/guide/topics/manifest/uses-feature-element
   d.Configuration.SystemAvailableFeature = []string{
      "android.hardware.touchscreen",
      "android.hardware.wifi",
   }
   d.Configuration.TouchScreen = 1
   return d
}

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (d Device) Upload(deviceID, auth string) error {
   buf, err := proto.Marshal(d)
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", bytes.NewReader(buf),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + auth},
      "User-Agent": {"Android-Finsky (sdk=99,versionCode=99999999)"},
      "X-DFE-Device-ID": {deviceID},
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return fmt.Errorf("status %q", res.Status)
   }
   return nil
}
