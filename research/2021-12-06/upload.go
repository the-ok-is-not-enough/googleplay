package main

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
   "strings"
   "time"
)

func main() {
   dev, err := NewDevice()
   if err != nil {
      panic(err)
   }
   if err := upload(dev); err != nil {
      panic(err)
   }
   time.Sleep(9 * time.Second)
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
   buf, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
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

func upload(dev *Device) error {
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", bytes.NewReader(body),
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


