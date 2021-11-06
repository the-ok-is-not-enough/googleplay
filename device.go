package googleplay

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

var DefaultCheckin = Checkin{Version: 3}

var DefaultConfig = Config{
   DeviceConfiguration{
      TouchScreen: 1,
      Keyboard: 1,
      Navigation: 1,
      ScreenLayout: 1,
      HasFiveWayNavigation: true,
      HasHardKeyboard: true,
      ScreenDensity: 1,
      // developer.android.com/guide/topics/manifest/uses-feature-element
      GlEsVersion: 0x0009_0000,
      // developer.android.com/guide/topics/manifest/uses-feature-element
      SystemAvailableFeature: []string{
         "android.hardware.touchscreen",
         "android.hardware.wifi",
      },
      // developer.android.com/ndk/guides/abis
      NativePlatform: []string{
         "armeabi-v7a",
      },
   },
}

func roundTrip(req *http.Request) (*http.Response, error) {
   dum, err := httputil.DumpRequest(req, false)
   if err != nil {
      return nil, err
   }
   os.Stdout.Write(dum)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != http.StatusOK {
      dum, err := httputil.DumpResponse(res, false)
      if err != nil {
         return nil, err
      }
      return nil, fmt.Errorf("%s", dum)
   }
   return res, nil
}

type Checkin struct {
   Checkin struct{} `json:"checkin"`
   Version int `json:"version"`
}

type Config struct {
   DeviceConfiguration DeviceConfiguration `protobuf:"bytes,1"`
}

type Device struct {
   Android_ID int64
}

func NewDevice(check Checkin) (*Device, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(check)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", origin + "/checkin", buf)
   if err != nil {
      return nil, err
   }
   res, err := roundTrip(req)
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

// Read Device from file.
func (d *Device) Decode(r io.Reader) error {
   return json.NewDecoder(r).Decode(d)
}

// Write Device to file.
func (d Device) Encode(w io.Writer) error {
   enc := json.NewEncoder(w)
   enc.SetIndent("", " ")
   return enc.Encode(d)
}

func (d Device) String() string {
   return strconv.FormatInt(d.Android_ID, 16)
}

type DeviceConfiguration struct {
   TouchScreen            int32   `protobuf:"varint,1"`
   Keyboard               int32   `protobuf:"varint,2"`
   Navigation             int32   `protobuf:"varint,3"`
   ScreenLayout           int32   `protobuf:"varint,4"`
   HasFiveWayNavigation   bool    `protobuf:"varint,6"`
   HasHardKeyboard        bool    `protobuf:"varint,5"`
   ScreenDensity          int32   `protobuf:"varint,7"`
   GlEsVersion            int32   `protobuf:"varint,8"`
   SystemAvailableFeature []string `protobuf:"bytes,10"`
   NativePlatform         []string `protobuf:"bytes,11"`
}
