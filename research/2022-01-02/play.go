package play

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
)

type response struct {
   *http.Response
}

func (r response) Error() string {
   return r.Status
}

var defaultConfig = config{
   deviceFeature: []string{
      "android.hardware.screen.portrait",
      "android.hardware.touchscreen",
      "android.hardware.wifi",
   },
   glEsVersion: 0x0003_0001,
   nativePlatform: []string{
      "x86",
   },
   screenWidth: 1,
}

type device struct {
   androidID uint64
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

////////////////////////////////////////////////////////////////////////////////

// A Sleep might be needed after this or the /auth call.
func checkin(con config) (*device, error) {
   androidCheckinRequest := protobuf.Message{
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
   for _, feat := range con.deviceFeature {
      androidCheckinRequest.Add(26, protobuf.Message{
         {1, "name"}: feat,
      })
   }
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/checkin",
      androidCheckinRequest.Encode(),
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Content-Type": {"application/x-protobuffer"},
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return 0, response{res}
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var dev device
   dev.androidID = mes.GetUint64(7)
   return &dev, nil
}

func details(app string) (uint64, error) {
   dev, err := checkin(defaultConfig)
   if err != nil {
      return 0, err
   }
   sID := strconv.FormatUint(dev.androidID, 16)
   fmt.Println(sID)
   req5 := &http.Request{
      Header: http.Header{
         "Authorization": {"Bearer " + auth},
         "X-Dfe-Device-Id": {sID},
      },
      Method: "GET",
      URL: &url.URL{
         Host: "android.clients.google.com",
         Path: "/fdfe/details",
         RawQuery: "doc=" + app,
         Scheme: "https",
      },
   }
   res, err := new(http.Transport).RoundTrip(req5)
   if err != nil {
      return 0, err
   }
   if res.StatusCode != http.StatusOK {
      return 0, response{res}
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(1,2,4,13,1,3), nil
}

