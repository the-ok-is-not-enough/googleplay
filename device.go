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

var DefaultConfig = Object{
   1: Object{
      1: uint64(1),
      2: uint64(1),
      3: uint64(1),
      4: uint64(1),
      5: true,
      6: true,
      7: uint64(1),
      // developer.android.com/guide/topics/manifest/uses-feature-element
      8: uint64(0x0009_0000),
      // developer.android.com/guide/topics/manifest/uses-feature-element
      10: Array{
         // com.pinterest
         "android.hardware.camera",
         // com.pinterest
         "android.hardware.faketouch",
         // com.pinterest
         "android.hardware.location",
         // com.pinterest
         "android.hardware.screen.portrait",
         // com.google.android.youtube
         "android.hardware.touchscreen",
         // com.google.android.youtube
         "android.hardware.wifi",
      },
      // developer.android.com/ndk/guides/abis
      11: Array{
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
