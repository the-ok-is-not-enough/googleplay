package googleplay

import (
   "bytes"
   "github.com/89z/parse/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

func (a Auth) Delivery(dev *Device, app string, ver int) (protobuf.Decoder, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/delivery", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Get("Auth")},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = url.Values{
      "doc": {app},
      "vc": {strconv.Itoa(ver)},
   }.Encode()
   res, err := roundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return protobuf.NewDecoder(buf), nil
}

func (a Auth) Details(dev *Device, app string) (protobuf.Decoder, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Get("Auth")},
      "X-DFE-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = url.Values{
      "doc": {app},
   }.Encode()
   res, err := roundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return protobuf.NewDecoder(buf), nil
}

type Config struct {
   DeviceConfiguration DeviceConfiguration `json:"1"`
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

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (a Auth) Upload(dev *Device, con Config) error {
   enc, err := protobuf.NewEncoder(con)
   if err != nil {
      return err
   }
   buf, err := enc.Encode()
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
      "Authorization": {"Bearer " + a.Get("Auth")},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   res, err := roundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

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
      NativePlatform: []string{
         "armeabi-v7a",
      },
   },
}
