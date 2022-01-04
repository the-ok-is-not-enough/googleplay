package googleplay

import (
   "bytes"
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "os"
   "strconv"
   "strings"
   "time"
)

const (
   Sleep = 4 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
   origin = "https://android.clients.google.com"
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

var Log = format.Log{Writer: os.Stdout}

var purchaseRequired = response{
   &http.Response{StatusCode: 3, Status: "purchase required"},
}

type Auth struct {
   Auth string
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
   Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   responseWrapper, err := protobuf.Unmarshal(buf)
   if err != nil {
      return nil, err
   }
   status := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      GetUint(1, "status")
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
   req.Header = http.Header{
      "Authorization": []string{"Bearer " + a.Auth},
      "X-Dfe-Device-ID": []string{dev.String()},
   }
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   Log.Dump(req)
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
   det.CurrencyCode = docV2.Get(8, "offer").GetString(2, "currencyCode")
   det.Micros = docV2.Get(8, "offer").GetUint(1, "micros")
   det.NumDownloads = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetUint(70, "numDownloads")
   // The shorter path 13,1,9 returns wrong size for some packages:
   // com.riotgames.league.wildriftvn
   det.Size = docV2.Get(13, "details").
      Get(1, "appDetails").
      Get(34, "installDetails").
      GetUint(2, "size")
   det.Title = docV2.GetString(5, "title")
   det.VersionCode = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetUint(3, "versionCode")
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
   ScreenWidth uint64
   SystemSharedLibrary []string
   // this can be 0, but it must be included:
   TouchScreen uint64
}

type Delivery struct {
   DownloadURL string
   SplitDeliveryData []SplitDeliveryData
}

type Details struct {
   Title string
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
   fmt.Fprintln(str, "VersionString:", d.VersionString)
   fmt.Fprintln(str, "VersionCode:", d.VersionCode)
   fmt.Fprint(str, "NumDownloads: ")
   format.Number.LabelUint64(str, d.NumDownloads)
   fmt.Fprintln(str)
   fmt.Fprint(str, "Size: ")
   format.Size.LabelUint64(str, d.Size)
   fmt.Fprintln(str)
   fmt.Fprintf(str, "Offer: %.2f ", float64(d.Micros)/1_000_000)
   fmt.Fprint(str, d.CurrencyCode)
   return str.String()
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
      "POST", origin + "/checkin", bytes.NewReader(checkinRequest.Marshal()),
   )
   req.Header.Set("Content-Type", "application/x-protobuffer")
   Log.Dump(req)
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
   dev.AndroidID = checkinResponse.GetUint(7, "androidId")
   return &dev, nil
}

func (d Device) String() string {
   return strconv.FormatUint(d.AndroidID, 16)
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}

type response struct {
   *http.Response
}

func (r response) Error() string {
   return strconv.Itoa(r.StatusCode) + " " + r.Status
}
