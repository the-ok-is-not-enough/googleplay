package googleplay

import (
   "bytes"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

var DefaultConfig = Config{
   DeviceFeature: []string{
      // com.instagram.android
      "android.hardware.bluetooth",
      // com.xiaomi.smarthome
      "android.hardware.bluetooth_le",
      // com.pinterest
      "android.hardware.camera",
      // com.xiaomi.smarthome
      "android.hardware.camera.autofocus",
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
      // org.thoughtcrime.securesms
      "android.hardware.telephony",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      // com.xiaomi.smarthome
      "android.hardware.usb.host",
      // com.google.android.youtube
      "android.hardware.wifi",
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
   },
   ScreenWidth: 1,
   SystemSharedLibrary: []string{
      // com.miui.weather2
      "global-miui11-empty.jar",
   },
   // com.valvesoftware.android.steam.community
   TouchScreen: 3,
}

type Config struct {
   DeviceFeature []string
   GLESversion uint64
   GLextension []string
   HasFiveWayNavigation uint64
   HasHardKeyboard uint64
   Keyboard uint64
   NativePlatform []string
   Navigation uint64
   ScreenDensity uint64
   ScreenLayout uint64
   ScreenWidth uint64
   SystemSharedLibrary []string
   TouchScreen uint64
}

type Details struct {
   VersionCode uint64
   VersionString string
}

type Device struct {
   AndroidID uint64
}

// A Sleep is needed after this.
func Checkin(con Config) (*Device, error) {
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
         {12, "screenWidth"}: con.ScreenWidth,
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
      "POST", "https://android.clients.google.com/checkin",
      bytes.NewReader(checkinRequest.Marshal()),
   )
   req.Header.Set("Content-Type", "application/x-protobuffer")
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, response{res}
   }
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   checkinResponse, err := protobuf.Unmarshal(buf)
   if err != nil {
      return nil, err
   }
   var dev Device
   dev.AndroidID = checkinResponse.GetUint64(7, "androidId")
   return &dev, nil
}

func (d Device) String() string {
   return strconv.FormatUint(d.AndroidID, 16)
}

type response struct {
   *http.Response
}

func (r response) Error() string {
   return r.Status
}

type Auth struct {
   Auth string
}

func (a Auth) Details(dev *Device, app string) (*Details, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/details", nil,
   )
   req.Header = http.Header{
      "Authorization": []string{"Bearer " + a.Auth},
      "X-Dfe-Device-ID": []string{dev.String()},
   }
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != http.StatusOK {
      return nil, response{res}
   }
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   responseWrapper, err := protobuf.Unmarshal(buf)
   if err != nil {
      return nil, err
   }
   var det Details
   docV2 := responseWrapper.Get(1, "payload").
      Get(2, "detailsResponse").
      Get(4, "docV2")
   det.VersionCode = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetUint64(3, "versionCode")
   det.VersionString = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(4, "versionString")
   return &det, nil
}

// FIXME
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
   LogLevel.Dump(req)
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
   if deliveryResponse.GetUint64(1) == uint64(purchaseRequired.StatusCode) {
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
