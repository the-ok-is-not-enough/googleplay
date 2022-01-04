package play

import (
   "bytes"
   "fmt"
   "github.com/89z/format/protobuf"
   "io"
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

var defaultConfig = config{
   deviceFeature: []string{
      "android.hardware.touchscreen",
      "android.hardware.screen.portrait",
      "android.hardware.wifi",
   }
   glEsVersion: 0x3_0000,
   nativePlatform: []string{
      "x86",
   },
   screenWidth: 1,
}

////////////////////////////////////////////////////////////////////////////////

// A Sleep might be needed after this.
func checkinProto() (uint64, error) {
   androidCheckinRequest := protobuf.Message{
      {4, "checkin"}:protobuf.Message{
         {1, "build"}:protobuf.Message{
            {10, "sdkVersion"}:uint64(29),
         },
      },
      {14, "version"}:uint64(3),
      {18, "deviceConfiguration"}: protobuf.Message{
         {1, "touchScreen"}:uint64(0),
         {2, "keyboard"}:uint64(0),
         {3, "navigation"}:uint64(0),
         {4, "screenLayout"}:uint64(0),
         {5, "hasHardKeyboard"}:uint64(0),
         {6, "hasFiveWayNavigation"}:uint64(0),
         {7, "screenDensity"}:uint64(0),
         {8, "glEsVersion"}:uint64(0x3_0000),
         {11, "nativePlatform"}:[]string{"x86"},
         {12, "screenWidth"}:uint64(1),
         {26, "deviceFeature"}:[]protobuf.Message{
            protobuf.Message{
               {1, "name"}:"android.hardware.touchscreen",
            },
            protobuf.Message{
               {1, "name"}:"android.hardware.screen.portrait",
            },
            protobuf.Message{
               {1, "name"}:"android.hardware.wifi",
            },
         },
      },
   }
   req0 := &http.Request{
      Body: io.NopCloser(bytes.NewReader(androidCheckinRequest.Marshal())),
      Header: http.Header{
         "Content-Type": {"application/x-protobuffer"},
      },
      Method:"POST",
      URL:&url.URL{
         Host:"android.clients.google.com",
         Path:"/checkin",
         Scheme:"https",
      },
   }
   res, err := new(http.Transport).RoundTrip(req0)
   if err != nil {
      return 0, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return 0, response{res}
   }
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return 0, err
   }
   mes, err := protobuf.Unmarshal(buf)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(7, "androidId"), nil
}

func details(id uint64, app string) (uint64, error) {
   sID := strconv.FormatUint(id, 16)
   fmt.Println(sID)
   req5 := &http.Request{
      Header:http.Header{
         "Authorization":[]string{"Bearer " + auth},
         "X-Dfe-Device-Id":[]string{sID},
      },
      Method:"GET",
      URL:&url.URL{
         Host:"android.clients.google.com",
         Path:"/fdfe/details",
         RawQuery:"doc=" + app,
         Scheme:"https",
      },
   }
   res, err := new(http.Transport).RoundTrip(req5)
   if err != nil {
      return 0, err
   }
   if res.StatusCode != http.StatusOK {
      return 0, response{res}
   }
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return 0, err
   }
   mes, err := protobuf.Unmarshal(buf)
   if err != nil {
      return 0, err
   }
   code := mes.Get(1, "payload").
      Get(2, "detailsResponse").
      Get(4, "docV2").
      Get(13, "details").
      Get(1, "appDetails").
      GetUint64(3, "versionCode")
   return code, nil
}
