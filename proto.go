package googleplay

import (
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

var DefaultConfig = Config{
   DeviceFeature: []string{
      // com.google.android.apps.walletnfcrel
      "android.software.device_admin",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      "android.hardware.wifi",
      // com.instagram.android
      "android.hardware.bluetooth",
      // com.pinterest
      "android.hardware.camera",
      "android.hardware.location",
      "android.hardware.screen.portrait",
      // com.smarty.voomvoom
      "android.hardware.location.gps",
      "android.hardware.sensor.accelerometer",
      // com.tgc.sky.android
      "android.hardware.touchscreen.multitouch",
      "android.hardware.touchscreen.multitouch.distinct",
      "android.hardware.vulkan.level",
      "android.hardware.vulkan.version",
      // org.videolan.vlc
      "android.hardware.screen.landscape",
      // com.vimeo.android.videoapp
      "android.hardware.microphone",
      // com.xiaomi.smarthome
      "android.hardware.bluetooth_le",
      "android.hardware.camera.autofocus",
      "android.hardware.usb.host",
      // org.thoughtcrime.securesms
      "android.hardware.telephony",
      // se.pax.calima
      "android.hardware.location.network",
   },
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
      // com.exnoa.misttraingirls
      "arm64-v8a",
   },
   SystemSharedLibrary: []string{
      // com.miui.weather2
      "global-miui11-empty.jar",
   },
   // com.valvesoftware.android.steam.community
   TouchScreen: 3,
}

func (a Auth) Delivery(dev *Device, app string, ver int64) (*Delivery, error) {
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
      "vc": {strconv.FormatInt(ver, 10)},
   }.Encode()
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   status := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      GetUint64(1, "status")
   if int(status) == purchaseRequired.StatusCode {
      return nil, purchaseRequired
   }
   var del Delivery
   deliveryData := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      Get(2, "appDeliveryData")
   del.DownloadURL = deliveryData.GetString(3, "downloadUrl")
   for _, split := range deliveryData.GetMessages(15, "splitDeliveryData") {
      var dSplit SplitDeliveryData
      dSplit.ID = split.GetString(1, "id")
      dSplit.DownloadURL = split.GetString(5, "downloadUrl")
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
      // This is needed for some apps, for example:
      // com.xiaomi.smarthome
      "User-Agent": {agent},
      "X-Dfe-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var det Details
   docV2 := responseWrapper.Get(1, "payload").
      Get(2, "detailsResponse").
      Get(4, "docV2")
   det.CurrencyCode = docV2.Get(8, "offer").GetString(2, "currencyCode")
   det.Micros = docV2.Get(8, "offer").GetUint64(1, "micros")
   det.NumDownloads = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetUint64(70, "numDownloads")
   // The shorter path 13,1,9 returns wrong size for some packages:
   // com.riotgames.league.wildriftvn
   det.Size = docV2.Get(13, "details").
      Get(1, "appDetails").
      Get(34, "installDetails").
      GetUint64(2, "size")
   det.Title = docV2.GetString(5, "title")
   det.UploadDate = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(16, "uploadDate")
   det.VersionCode = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetUint64(3, "versionCode")
   det.VersionString = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(4, "versionString")
   return &det, nil
}

type Config struct {
   DeviceFeature []string
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
   SystemSharedLibrary []string
   // this can be 0, but it must be included:
   TouchScreen uint64
}

type Details struct {
   Title string
   UploadDate string
   VersionString string
   VersionCode uint64
   NumDownloads uint64
   Size uint64
   Micros uint64
   CurrencyCode string
}

func (d Details) String() string {
   str := new(strings.Builder)
   fmt.Fprintln(str, "Title:", d.Title)
   fmt.Fprintln(str, "UploadDate:", d.UploadDate)
   fmt.Fprintln(str, "VersionString:", d.VersionString)
   fmt.Fprintln(str, "VersionCode:", d.VersionCode)
   fmt.Fprint(str, "NumDownloads: ")
   format.Number.Uint64(str, d.NumDownloads)
   fmt.Fprintln(str)
   fmt.Fprint(str, "Size: ")
   format.Size.Uint64(str, d.Size)
   fmt.Fprintln(str)
   fmt.Fprintf(str, "Offer: %.2f ", float64(d.Micros)/1_000_000)
   fmt.Fprint(str, d.CurrencyCode)
   return str.String()
}

// A Sleep is needed after this.
func NewDevice(con Config) (*Device, error) {
   checkinRequest := protobuf.Message{
      {4, "checkin"}: protobuf.Message{
         {1, "build"}: protobuf.Message{
            {10, "sdkVersion"}: uint64(29),
         },
      },
      {14, "version"}: uint64(3),
      {18, "deviceConfiguration"}: protobuf.Message{
         {1, "touchScreen"}: con.TouchScreen,
         {2, "keyboard"}: con.Keyboard,
         {3, "navigation"}: con.Navigation,
         {4, "screenLayout"}: con.ScreenLayout,
         {5, "hasHardKeyboard"}: con.HasHardKeyboard,
         {6, "hasFiveWayNavigation"}: con.HasFiveWayNavigation,
         {7, "screenDensity"}: con.ScreenDensity,
         {8, "glEsVersion"}: con.GLESversion,
         {9, "systemSharedLibrary"}: con.SystemSharedLibrary,
         {11, "nativePlatform"}: con.NativePlatform,
         {15, "glExtension"}: con.GLextension,
      },
   }
   for _, feature := range con.DeviceFeature {
      checkinRequest.Get(18, "deviceConfiguration").
      Add(26, "deviceFeature", protobuf.Message{
         {1, "name"}: feature,
      })
   }
   req, err := http.NewRequest(
      "POST", origin + "/checkin", checkinRequest.Encode(),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-protobuffer")
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   checkinResponse, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var dev Device
   dev.AndroidID = checkinResponse.GetUint64(7, "androidId")
   return &dev, nil
}
