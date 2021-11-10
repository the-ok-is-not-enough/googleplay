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
   DeviceConfiguration: DeviceConfiguration{
      TouchScreen: 1,
      Keyboard: 1,
      Navigation: 1,
      ScreenLayout: 1,
      HasHardKeyboard: true,
      HasFiveWayNavigation: true,
      ScreenDensity: 1,
      // developer.android.com/guide/topics/manifest/uses-feature-element
      GlEsVersion: 0x0009_0000,
      // developer.android.com/guide/topics/manifest/uses-feature-element
      SystemAvailableFeature: []string{
         // com.pinterest
         "android.hardware.camera",
         // com.pinterest
         "android.hardware.location",
         // org.videolan.vlc
         "android.hardware.screen.landscape",
         // com.pinterest
         "android.hardware.screen.portrait",
         // com.google.android.youtube
         "android.hardware.touchscreen",
         // com.google.android.youtube
         "android.hardware.wifi",
      },
      // developer.android.com/ndk/guides/abis
      NativePlatform: []string{
         "x86",
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
   DeviceConfiguration DeviceConfiguration `json:"1"`
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
   TouchScreen int32 `json:"1"`
   Keyboard int32 `json:"2"`
   Navigation int32 `json:"3"`
   ScreenLayout int32 `json:"4"`
   HasHardKeyboard bool `json:"5"`
   HasFiveWayNavigation bool `json:"6"`
   ScreenDensity int32 `json:"7"`
   GlEsVersion int32 `json:"8"`
   SystemAvailableFeature []string `json:"10"`
   NativePlatform []string `json:"11"`
}
