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
   Android_ID int64
}

// Read Checkin from file.
func (c *Checkin) Decode(r io.Reader) error {
   return json.NewDecoder(r).Decode(c)
}

// Write Checkin to file.
func (c Checkin) Encode(w io.Writer) error {
   enc := json.NewEncoder(w)
   enc.SetIndent("", " ")
   return enc.Encode(c)
}

func (c Checkin) String() string {
   return strconv.FormatInt(c.Android_ID, 16)
}

type CheckinRequest struct {
   Checkin struct{} `json:"checkin"`
   Version int `json:"version"`
}

func NewCheckinRequest() CheckinRequest {
   return CheckinRequest{Version: 3}
}

func (c CheckinRequest) Post() (*Checkin, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(c)
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
   check := new(Checkin)
   if err := json.NewDecoder(res.Body).Decode(check); err != nil {
      return nil, err
   }
   return check, nil
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
