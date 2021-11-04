package googleplay

import (
   "bytes"
   "encoding/json"
   "fmt"
   "github.com/segmentio/encoding/proto"
   "net/http"
   "strconv"
)

const origin = "https://android.clients.google.com"

type checkinRequest struct {
   Checkin struct{} `json:"checkin"`
   Version int `json:"version"`
}

type checkin struct {
   Android_ID int64
}

func newCheckin() (*checkin, error) {
   cReq := checkinRequest{Version: 3}
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(cReq)
   if err != nil {
      return nil, err
   }
   res, err := http.Post(origin + "/checkin", "application/json", buf)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("status %q", res.Status)
   }
   check := new(checkin)
   if err := json.NewDecoder(res.Body).Decode(check); err != nil {
      return nil, err
   }
   return check, nil
}

func (c checkin) String() string {
   return strconv.FormatInt(c.Android_ID, 16)
}

type device struct {
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

func newDevice() device {
   var d device
   d.Configuration.GlEsVersion=196610
   d.Configuration.HasFiveWayNavigation = true
   d.Configuration.HasHardKeyboard = true
   d.Configuration.Keyboard = 1
   d.Configuration.NativePlatform = []string{
      "arm64-v8a","armeabi-v7a", "armeabi",
   }
   d.Configuration.Navigation = 1
   d.Configuration.ScreenDensity = 420
   d.Configuration.ScreenLayout = 2
   d.Configuration.SystemAvailableFeature = []string{
      "android.hardware.touchscreen","android.hardware.wifi",
   }
   d.Configuration.TouchScreen = 3
   return d
}

////////////////////////////////////////////////////////////////////////////////

func (d device) upload(auth, deviceID string) error {
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
      "Accept": {"*/*"},
      "Authorization": {"Bearer " + auth},
      "Content-Type": {"application/x-protobuf"},
      "User-Agent": {"Android-Finsky (versionCode=81031200,sdk=27)"},
      "X-DFE-Device-Id": {deviceID},
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
