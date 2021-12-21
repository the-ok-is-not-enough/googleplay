package googleplay

import (
   "github.com/89z/parse/protobuf"
   "net/http"
   "net/url"
   "strconv"
)

var DefaultConfig = Config{
   // com.axis.drawingdesk.v3
   GLESversion: 0x0003_0001,
   GLextension: []string{
      // com.instagram.android
      "GL_OES_compressed_ETC1_RGB8_texture",
   },
   NativePlatform: []string{
      // com.vimeo.android.videoapp
      "x86",
      // com.axis.drawingdesk.v3
      "armeabi-v7a",
   },
   SystemAvailableFeature: []string{
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
   },
}

var purchaseRequired = response{3, "purchase required"}

func (a Auth) Delivery(dev *Device, app string, ver int) (*Delivery, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/delivery", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Auth},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = "vc=" + strconv.Itoa(ver) + "&doc=" + url.QueryEscape(app)
   dumpRequest(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   deliveryResponse := responseWrapper.Get(1, 21)
   if deliveryResponse.GetUint64(1) == purchaseRequired.code {
      return nil, purchaseRequired
   }
   appDeliveryData := deliveryResponse.Get(2)
   var del Delivery
   del.DownloadURL = appDeliveryData.GetString(3)
   for _, split := range appDeliveryData.GetMessages(15) {
      var dSplit SplitDeliveryData
      dSplit.ID = split.GetString(1)
      dSplit.DownloadURL = split.GetString(5)
      del.SplitDeliveryData = append(del.SplitDeliveryData, dSplit)
   }
   return &del, nil
}

func (a Auth) Details(dev *Device, app string) (*Details, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Auth},
      "X-DFE-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   dumpRequest(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   docV2 := responseWrapper.Get(1, 2, 4)
   var det Details
   det.InstallationSize.Size = docV2.GetUint64(13, 1, 9)
   det.NumDownloads.Num = docV2.GetUint64(13, 1, 70)
   det.Offer.CurrencyCode = docV2.GetString(8, 2)
   det.Offer.Micros = docV2.GetUint64(8, 1)
   det.Title = docV2.GetString(5)
   det.VersionCode = docV2.GetUint64(13, 1, 3)
   det.VersionString = docV2.GetString(13, 1, 4)
   return &det, nil
}

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (a Auth) Upload(dev *Device, con Config) error {
   uploadDeviceConfigRequest := protobuf.Message{
      {1, "deviceConfiguration"}: protobuf.Message{
         {1, "touchScreen"}: con.TouchScreen,
         {2, "keyboard"}: con.Keyboard,
         {3, "navigation"}: con.Navigation,
         {4, "screenLayout"}: con.ScreenLayout,
         {5, "hasHardKeyboard"}: con.HasHardKeyboard,
         {6, "hasFiveWayNavigation"}: con.HasFiveWayNavigation,
         {7, "screenDensity"}: con.ScreenDensity,
         {8, "glEsVersion"}: con.GLESversion,
         {10, "systemAvailableFeature"}: con.SystemAvailableFeature,
         {11, "nativePlatform"}: con.NativePlatform,
         {15, "glExtension"}: con.GLextension,
      },
   }
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig",
      uploadDeviceConfigRequest.Encode(),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Auth},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   dumpRequest(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

type Config struct {
   GLESversion uint64
   GLextension []string
   // this can be 0, but it must be included:
   HasFiveWayNavigation uint64
   // this can be 0, but it must be included:
   HasHardKeyboard uint64
   // this can be 0, but it must be included:
   Keyboard uint64
   NativePlatform []string
   // this can be 0, but it must be included:
   Navigation uint64
   // this can be 0, but it must be included:
   ScreenDensity uint64
   // this can be 0, but it must be included:
   ScreenLayout uint64
   SystemAvailableFeature []string
   // this can be 0, but it must be included:
   TouchScreen uint64
}

type Delivery struct {
   DownloadURL string
   SplitDeliveryData []SplitDeliveryData
}

type Details struct {
   InstallationSize InstallationSize
   NumDownloads NumDownloads
   Offer Offer
   Title string
   VersionCode uint64
   VersionString string
}

type InstallationSize struct {
   Size uint64
}

func (i InstallationSize) String() string {
   val := float64(i.Size)
   metric := []string{" B", " kB", " MB"}
   return numberFormat(val, metric)
}

type NumDownloads struct {
   Num uint64
}

func (n NumDownloads) String() string {
   val := float64(n.Num)
   metric := []string{"", " k", " M", " B"}
   return numberFormat(val, metric)
}

type Offer struct {
   Micros uint64
   CurrencyCode string
}

func (o Offer) String() string {
   val := float64(o.Micros) / 1_000_000
   return strconv.FormatFloat(val, 'f', 2, 64) + " " + o.CurrencyCode
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}
