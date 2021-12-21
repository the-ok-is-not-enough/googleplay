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

func numberFormat(val float64, metric []string) string {
   var key int
   for val >= 1000 {
      val /= 1000
      key++
   }
   if key >= len(metric) {
      return ""
   }
   return strconv.FormatFloat(val, 'f', 3, 64) + " " + metric[key]
}

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
   req.URL.RawQuery = url.Values{
      "doc": {app},
      "vc": {strconv.Itoa(ver)},
   }.Encode()
   dumpRequest(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var del Delivery
   appDeliveryData := mes.Get(1, 21, 2)
   del.DownloadURL = appDeliveryData.GetString(3)
   for _, mes := range appDeliveryData.GetMessages(15) {
      split := SplitDeliveryData{
         mes.GetString(1), mes.GetString(5),
      }
      del.SplitDeliveryData = append(del.SplitDeliveryData, split)
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
   req.URL.RawQuery = url.Values{
      "doc": {app},
   }.Encode()
   dumpRequest(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var det Details
   det.Title = mes.GetString(1, 2, 4, 5)
   det.VersionCode = mes.GetUint64(1, 2, 4, 13, 1, 3)
   det.VersionString = mes.GetString(1, 2, 4, 13, 1, 4)
   det.InstallationSize.Size = mes.GetUint64(1, 2, 4, 13, 1, 9)
   det.Offer.Micros = mes.GetUint64(1, 2, 4, 8, 1)
   det.Offer.CurrencyCode = mes.GetString(1, 2, 4, 8, 2)
   return &det, nil
}

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (a Auth) Upload(dev *Device, con Config) error {
   mes := protobuf.Message{
      1: protobuf.Message{
         1: con.TouchScreen,
         2: con.Keyboard,
         3: con.Navigation,
         4: con.ScreenLayout,
         5: con.HasHardKeyboard,
         6: con.HasFiveWayNavigation,
         7: con.ScreenDensity,
         8: con.GLESversion,
         10: con.SystemAvailableFeature,
         11: con.NativePlatform,
         15: con.GLextension,
      },
   }
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", mes.Encode(),
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
   Title string
   VersionCode uint64
   VersionString string
   InstallationSize InstallationSize
   Offer Offer
}

type InstallationSize struct {
   Size uint64
}

func (i InstallationSize) String() string {
   val := float64(i.Size)
   metric := []string{"B", "kB", "MB"}
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
