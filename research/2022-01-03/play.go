package play

import (
   "bytes"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

var defaultConfig = config{
   deviceFeature: []string{
      "android.hardware.touchscreen",
      "android.hardware.screen.portrait",
      "android.hardware.wifi",
   },
   glEsVersion: 0x3_0000,
   nativePlatform: []string{
      "x86",
   },
   screenWidth: 1,
}

type config struct {
   deviceFeature []string
   glEsVersion uint64
   hasFiveWayNavigation uint64
   hasHardKeyboard uint64
   keyboard uint64
   nativePlatform []string
   navigation uint64
   screenDensity uint64
   screenLayout uint64
   screenWidth uint64
   touchScreen uint64
}

type details struct {
   versionCode uint64
   versionString string
}

func newDetails(dev *device, app string) (*details, error) {
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
   var det details
   docV2 := responseWrapper.
      Get(1, "payload").
      Get(2, "detailsResponse").
      Get(4, "docV2")
   det.versionCode = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetUint64(3, "versionCode")
   det.versionString = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(4, "versionString")
   return &det, nil
}

type device struct {
   androidID uint64
}

// A Sleep is needed after this.
func checkin(con config) (*device, error) {
   checkinRequest := protobuf.Message{
      {4, "checkin"}: protobuf.Message{
         {1, "build"}: protobuf.Message{
            {10, "sdkVersion"}: uint64(29),
         },
      },
      {14, "version"}: uint64(3),
      {18, "deviceConfiguration"}: protobuf.Message{
         {1, "touchScreen"}: con.touchScreen,
         {2, "keyboard"}: con.keyboard,
         {3, "navigation"}: con.navigation,
         {4, "screenLayout"}: con.screenLayout,
         {5, "hasHardKeyboard"}: con.hasHardKeyboard,
         {6, "hasFiveWayNavigation"}: con.hasFiveWayNavigation,
         {7, "screenDensity"}: con.screenDensity,
         {8, "glEsVersion"}: con.glEsVersion,
         {11, "nativePlatform"}: con.nativePlatform,
         {12, "screenWidth"}: con.screenWidth,
      },
   }
   for _, feature := range con.deviceFeature {
      checkinRequest.
      Get(18, "deviceConfiguration").
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
   mes, err := protobuf.Unmarshal(buf)
   if err != nil {
      return nil, err
   }
   var dev device
   dev.androidID = mes.GetUint64(7, "androidId")
   return &dev, nil
}

func (d device) String() string {
   return strconv.FormatUint(d.androidID, 16)
}

type response struct {
   *http.Response
}

func (r response) Error() string {
   return r.Status
}
