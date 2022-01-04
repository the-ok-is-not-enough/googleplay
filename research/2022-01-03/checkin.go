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
      // com.pinterest
      "android.hardware.camera",
      // com.pinterest
      "android.hardware.location",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      // com.pinterest
      "android.hardware.screen.portrait",
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
      "x86",
   },
   ScreenWidth: 1,
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
   TouchScreen uint64
}

type Details struct {
   VersionCode uint64
   VersionString string
}

func NewDetails(dev *Device, app string) (*Details, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/details", nil,
   )
   req.Header = http.Header{
      "Authorization": []string{"Bearer " + auth},
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
