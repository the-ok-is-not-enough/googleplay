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
      // com.google.android.youtube
      GlEsVersion: 0x0002_0000,
      SystemAvailableFeature: []string{
         // com.pinterest
         "android.hardware.camera",
         // com.pinterest
         "android.hardware.location",
         // com.vimeo.android.videoapp
         "android.hardware.microphone",
         // org.videolan.vlc
         "android.hardware.screen.landscape",
         // com.pinterest
         "android.hardware.screen.portrait",
         // com.google.android.youtube
         "android.hardware.touchscreen",
         // com.google.android.youtube
         "android.hardware.wifi",
      },
      NativePlatform: []string{
         // com.google.android.youtube
         "x86",
      },
      GlExtension: []string{
         // com.instagram.android
         "GL_OES_compressed_ETC1_RGB8_texture",
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
   // this can be 0, but it must be included:
   TouchScreen int32 `json:"1"`
   // this can be 0, but it must be included:
   Keyboard int32 `json:"2"`
   // this can be 0, but it must be included:
   Navigation int32 `json:"3"`
   // this can be 0, but it must be included:
   ScreenLayout int32 `json:"4"`
   // this can be false, but it must be included:
   HasHardKeyboard bool `json:"5"`
   // this can be false, but it must be included:
   HasFiveWayNavigation bool `json:"6"`
   // this can be 0, but it must be included:
   ScreenDensity int32 `json:"7"`
   GlEsVersion int32 `json:"8"`
   SystemAvailableFeature []string `json:"10"`
   NativePlatform []string `json:"11"`
   GlExtension []string `json:"15"`
}
