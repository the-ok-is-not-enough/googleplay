package googleplay

import (
   "github.com/89z/parse/protobuf"
   "net/http"
   "net/url"
   "strconv"
)

type AppDeliveryData struct {
   protobuf.Message
}

func (a AppDeliveryData) DownloadURL() string {
   return a.GetString(3)
}

func (a AppDeliveryData) SplitDeliveryData() []SplitDeliveryData {
   var splits []SplitDeliveryData
   for _, mes := range a.GetMessages(15) {
      splits = append(splits, SplitDeliveryData{mes})
   }
   return splits
}

type AppDetails struct {
   protobuf.Message
}

func (a AppDetails) VersionCode() uint64 {
   return a.GetUint64(3)
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
   del := responseWrapper{mes}.payload().deliveryResponse()
   if del.Status() == purchaseRequired.statusCode {
      return nil, purchaseRequired
   }
   return &del, nil
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
   return responseWrapper{mes}.payload().detailsResponse(), nil
}

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (a Auth) Upload(dev *Device, con Config) error {
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", con.Encode(),
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

type Config struct {
   protobuf.Message
}

func NewConfig() Config {
   dev := NewDeviceConfiguration()
   // this can be 0, but it must be included:
   dev.SetTouchScreen(0)
   // this can be 0, but it must be included:
   dev.SetKeyboard(0)
   // this can be 0, but it must be included:
   dev.SetNavigation(0)
   // this can be 0, but it must be included:
   dev.SetScreenLayout(0)
   // this can be false, but it must be included:
   dev.SetHasHardKeyboard(false)
   // this can be false, but it must be included:
   dev.SetHasFiveWayNavigation(false)
   // this can be 0, but it must be included:
   dev.SetScreenDensity(0)
   // com.axis.drawingdesk.v3
   dev.SetGLESversion(0x0003_0001)
   dev.SetSystemAvailableFeature([]string{
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
   })
   dev.SetNativePlatform([]string{
      // com.vimeo.android.videoapp
      "x86",
      // com.axis.drawingdesk.v3
      "armeabi-v7a",
   })
   dev.SetGLextension([]string{
      // com.instagram.android
      "GL_OES_compressed_ETC1_RGB8_texture",
   })
   con := Config{
      make(protobuf.Message),
   }
   con.SetDeviceConfiguration(dev)
   return con
}

func (c Config) SetDeviceConfiguration(v DeviceConfiguration) bool {
   return c.Set(1, v)
}

type DeliveryResponse struct {
   protobuf.Message
}

func (d DeliveryResponse) AppDeliveryData() AppDeliveryData {
   return AppDeliveryData{d.Get(2)}
}

func (d DeliveryResponse) Status() uint64 {
   return d.GetUint64(1)
}

type Details struct {
   protobuf.Message
}

func (d Details) AppDetails() AppDetails {
   return AppDetails{d.Get(1)}
}

type DetailsResponse struct {
   protobuf.Message
}

func (d DetailsResponse) DocV2() DocV2 {
   return DocV2{d.Get(4)}
}

type DeviceConfiguration struct {
   protobuf.Message
}

func NewDeviceConfiguration() DeviceConfiguration {
   return DeviceConfiguration{
      make(protobuf.Message),
   }
}

func (d DeviceConfiguration) SetGLESversion(v int) bool {
   return d.Set(8, v)
}

func (d DeviceConfiguration) SetGLextension(v []string) bool {
   return d.Set(15, v)
}

func (d DeviceConfiguration) SetHasFiveWayNavigation(v bool) bool {
   return d.Set(6, v)
}

func (d DeviceConfiguration) SetHasHardKeyboard(v bool) bool {
   return d.Set(5, v)
}

func (d DeviceConfiguration) SetKeyboard(v int) bool {
   return d.Set(2, v)
}

func (d DeviceConfiguration) SetNativePlatform(v []string) bool {
   return d.Set(11, v)
}

func (d DeviceConfiguration) SetNavigation(v int) bool {
   return d.Set(3, v)
}

func (d DeviceConfiguration) SetScreenDensity(v int) bool {
   return d.Set(7, v)
}

func (d DeviceConfiguration) SetScreenLayout(v int) bool {
   return d.Set(4, v)
}

func (d DeviceConfiguration) SetSystemAvailableFeature(v []string) bool {
   return d.Set(10, v)
}

func (d DeviceConfiguration) SetTouchScreen(v int) bool {
   return d.Set(1, v)
}

type DocV2 struct {
   protobuf.Message
}

func (d DocV2) Details() Details {
   return Details{d.Get(13)}
}

type SplitDeliveryData struct {
   protobuf.Message
}

func (s SplitDeliveryData) DownloadURL() string {
   return s.GetString(5)
}

func (s SplitDeliveryData) ID() string {
   return s.GetString(1)
}

type payload struct {
   protobuf.Message
}

func (p payload) deliveryResponse() DeliveryResponse {
   return DeliveryResponse{p.Get(21)}
}

func (p payload) detailsResponse() *DetailsResponse {
   return &DetailsResponse{p.Get(2)}
}

type responseWrapper struct {
   protobuf.Message
}

func (r responseWrapper) payload() payload {
   return payload{r.Get(1)}
}
