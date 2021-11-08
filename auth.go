package googleplay

import (
   "bytes"
   "github.com/89z/parse/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "time"
)

const (
   Sleep = 16 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
)

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

type Auth struct {
   url.Values
}

func (a Auth) Delivery(dev *Device, app string, ver int) (*Delivery, error) {
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
   wrap := new(responseWrapper)
   if err := protobuf.NewDecoder(buf).Decode(wrap); err != nil {
      return nil, err
   }
   return &wrap.Payload.DeliveryResponse, nil
}

func (a Auth) Details(dev *Device, app string) (*Details, error) {
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
   wrap := new(responseWrapper)
   if err := protobuf.NewDecoder(buf).Decode(wrap); err != nil {
      return nil, err
   }
   return &wrap.Payload.DetailsResponse, nil
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

type Delivery struct {
   AppDeliveryData struct {
      DownloadURL string `json:"3"`
   } `json:"2"`
}

type Details struct {
   DocV2 struct {
      DocumentDetails struct {
         AppDetails struct {
            DeveloperName string `json:"1"`
            VersionCode int `json:"3"`
            Version string `json:"4"`
            InstallationSize int `json:"9"`
            Permission []string `json:"10"`
         } `json:"1"`
      } `json:"13"`
   } `json:"4"`
}

type responseWrapper struct {
   Payload struct {
      DetailsResponse Details `json:"2"`
      DeliveryResponse Delivery `json:"21"`
   } `json:"1"`
}
