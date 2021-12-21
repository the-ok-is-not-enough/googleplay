package googleplay

import (
   "github.com/89z/parse/protobuf"
   "net/http"
   "net/url"
   "strconv"
)

type Config struct {
   DeviceConfiguration struct {
      GLESversion uint64
      GLextension []string
      NativePlatform []string
      SystemAvailableFeature []string
      // this can be 0, but it must be included:
      ScreenDensity uint64
      // this can be 0, but it must be included:
      HasFiveWayNavigation uint64
      // this can be 0, but it must be included:
      HasHardKeyboard uint64
      // this can be 0, but it must be included:
      ScreenLayout uint64
      // this can be 0, but it must be included:
      Keyboard uint64
      // this can be 0, but it must be included:
      Navigation uint64
      // this can be 0, but it must be included:
      TouchScreen uint64
   }
}

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (a Auth) Upload(dev *Device, con Config) error {
   /*
   mes := make(protobuf.Message)
   mes.SetUint64(0, 1, 1)
   */
   req, err := http.NewRequest("POST", origin + "/fdfe/uploadDeviceConfig", nil)
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

func NewConfig() Config {
   var c Config
   // com.axis.drawingdesk.v3
   c.DeviceConfiguration.GLESversion = 0x0003_0001
   c.DeviceConfiguration.SystemAvailableFeature = []string{
      // com.smarty.voomvoom
      "android.hardware.bluetooth",
      // com.pinterest
      "android.hardware.camera",
      // com.pinterest
      "android.hardware.location",
      // com.smarty.voomvoom
      "android.hardware.location.gps",
      // com.vimeo.android.videoapp
      "android.hardware.microphone",
      // org.videolan.vlc
      "android.hardware.screen.landscape",
      // com.pinterest
      "android.hardware.screen.portrait",
      // com.smarty.voomvoom
      "android.hardware.sensor.accelerometer",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      // com.google.android.youtube
      "android.hardware.wifi",
   }
   c.DeviceConfiguration.NativePlatform = []string{
      // com.vimeo.android.videoapp
      "x86",
      // com.axis.drawingdesk.v3
      "armeabi-v7a",
   }
   c.DeviceConfiguration.GLextension = []string{
      // com.instagram.android
      "GL_OES_compressed_ETC1_RGB8_texture",
   }
   return c
}

type DetailsResponse struct {
   docV2 struct {
      details struct {
         appDetails struct {
            versionCode uint64
         }
      }
   }
}

func (a Auth) DeliveryResponse(dev *Device, app string, ver int) (*DeliveryResponse, error) {
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
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var del DeliveryResponse
   del.AppDeliveryData.DownloadURL = mes.GetString(1, 21, 2, 3)
   return &del, nil
}

type DeliveryResponse struct {
   AppDeliveryData struct {
      DownloadURL string
   }
}

func (a Auth) DetailsResponse(dev *Device, app string) (*DetailsResponse, error) {
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
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var det DetailsResponse
   det.docV2.details.appDetails.versionCode = mes.GetUint64(1, 2, 4, 13, 1, 3)
   return &det, nil
}
