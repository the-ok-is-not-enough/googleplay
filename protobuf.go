package googleplay

import (
   "bytes"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
)

// A Sleep is needed after this.
func NewDevice(con Config) (*Device, error) {
   checkin := protobuf.Message{
      tag(4, "checkin"): protobuf.Message{
         tag(1, "build"): protobuf.Message{
            tag(10, "sdkVersion"): uint64(29),
         },
      },
      tag(14, "version"): uint64(3),
      tag(18, "deviceConfiguration"): protobuf.Message{
         tag(1, "touchScreen"): con.TouchScreen,
         tag(2, "keyboard"): con.Keyboard,
         tag(3, "navigation"): con.Navigation,
         tag(4, "screenLayout"): con.ScreenLayout,
         tag(5, "hasHardKeyboard"): con.HasHardKeyboard,
         tag(6, "hasFiveWayNavigation"): con.HasFiveWayNavigation,
         tag(7, "screenDensity"): con.ScreenDensity,
         tag(8, "glEsVersion"): con.GLESversion,
         tag(9, "systemSharedLibrary"): con.SystemSharedLibrary,
         tag(11, "nativePlatform"): con.NativePlatform,
         tag(15, "glExtension"): con.GLextension,
      },
   }
   config := checkin.Get(tag(18, "deviceConfiguration"))
   for _, name := range con.DeviceFeature {
      feature := protobuf.Message{
         tag(1, "name"): name,
      }
      config.Add(tag(26, "deviceFeature"), feature)
   }
   req, err := http.NewRequest(
      "POST", origin + "/checkin", bytes.NewReader(checkin.Marshal()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-protobuffer")
   LogLevel.Dump(req)
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
   dev.AndroidID = checkinResponse.GetFixed64(tag(7, "androidId"))
   return &dev, nil
}

func (h Header) Delivery(app string, ver int64) (*Delivery, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/delivery", nil)
   if err != nil {
      return nil, err
   }
   req.Header = h.Header
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
   status := responseWrapper.Get(tag(1, "payload")).
      Get(tag(21, "deliveryResponse")).
      GetVarint(tag(1, "status"))
   switch status {
   case 2:
      return nil, errorString("Regional lockout")
   case 3:
      return nil, errorString("Purchase required")
   case 5:
      return nil, errorString("Invalid version")
   }
   appData := responseWrapper.Get(tag(1, "payload")).
      Get(tag(21, "deliveryResponse")).
      Get(tag(2, "appDeliveryData"))
   var del Delivery
   del.DownloadURL = appData.GetString(tag(3, "downloadUrl"))
   for _, data := range appData.GetMessages(tag(15, "splitDeliveryData")) {
      var split SplitDeliveryData
      split.ID = data.GetString(tag(1, "id"))
      split.DownloadURL = data.GetString(tag(5, "downloadUrl"))
      del.SplitDeliveryData = append(del.SplitDeliveryData, split)
   }
   return &del, nil
}

func (h Header) Details(app string) (*Details, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = h.Header
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != http.StatusOK {
      return nil, errorString(res.Status)
   }
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   docV2 := responseWrapper.Get(tag(1, "payload")).
      Get(tag(2, "detailsResponse")).
      Get(tag(4, "docV2"))
   var det Details
   det.CurrencyCode = docV2.Get(tag(8, "offer")).
      GetString(tag(2, "currencyCode"))
   det.Micros = docV2.Get(tag(8, "offer")).
      GetVarint(tag(1, "micros"))
   det.NumDownloads = docV2.Get(tag(13, "details")).
      Get(tag(1, "appDetails")).
      GetVarint(tag(70, "numDownloads"))
   // The shorter path 13,1,9 returns wrong size for some packages:
   // com.riotgames.league.wildriftvn
   det.Size = docV2.Get(tag(13, "details")).
      Get(tag(1, "appDetails")).
      Get(tag(34, "installDetails")).
      GetVarint(tag(2, "size"))
   det.Title = docV2.GetString(tag(5, "title"))
   det.UploadDate = docV2.Get(tag(13, "details")).
      Get(tag(1, "appDetails")).
      GetString(tag(16, "uploadDate"))
   det.VersionCode = docV2.Get(tag(13, "details")).
      Get(tag(1, "appDetails")).
      GetVarint(tag(3, "versionCode"))
   det.VersionString = docV2.Get(tag(13, "details")).
      Get(tag(1, "appDetails")).
      GetString(tag(4, "versionString"))
   return &det, nil
}
