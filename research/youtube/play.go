package play

import (
   "encoding/json"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func details(app string) (uint64, error) {
   id, err := checkinProto()
   if err != nil {
      return 0, err
   }
   sID := strconv.FormatUint(id, 16)
   if err := upload(sID); err != nil {
      return 0, err
   }
   req5 := &http.Request{
      Header:http.Header{
         "Authorization":[]string{"Bearer " + auth},
         "X-Dfe-Device-Id":[]string{sID}
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
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(1,2,4,13,1,3), nil
}

type response struct {
   *http.Response
}

func (r response) Error() string {
   return r.Status
}

func checkinProto() (uint64, error) {
   req0 := &http.Request{
      Body: io.NopCloser(checkinBody.Encode()),
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
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(7), nil
}

var body1 = protobuf.Message{
   protobuf.Tag{Number:1, String:""}: body0,
}

func upload(id string) error {
   var req0 = &http.Request{
      Method:"POST",
      URL:&url.URL{
         Scheme:"https",
         Host:"android.clients.google.com",
         Path:"/fdfe/uploadDeviceConfig",
      },
      Header:http.Header{
         "Authorization":[]string{"Bearer " + auth},
         "Host":[]string{"android.clients.google.com"},
         "User-Agent":[]string{"Android-Finsky (sdk=99,versionCode=99999999)"},
         "X-Dfe-Device-Id":[]string{id},
      },
      Body:io.NopCloser(body1.Encode()),
   }
   res, err := new(http.Transport).RoundTrip(req0)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   time.Sleep(16 * time.Second)
   return nil
}

var body0 = protobuf.Message{
   protobuf.Tag{Number:1, String:"touchScreen"}:uint64(0),
   protobuf.Tag{Number:2, String:"keyboard"}:uint64(0),
   protobuf.Tag{Number:3, String:"navigation"}:uint64(0),
   protobuf.Tag{Number:4, String:"screenLayout"}:uint64(0),
   protobuf.Tag{Number:5, String:"hasHardKeyboard"}:uint64(0),
   protobuf.Tag{Number:6, String:"hasFiveWayNavigation"}:uint64(0),
   protobuf.Tag{Number:7, String:"screenDensity"}:uint64(0),
   protobuf.Tag{Number:8, String:"glEsVersion"}:uint64(0x3_0000),
   protobuf.Tag{Number:11}:[]string{"x86"},
   protobuf.Tag{Number:12, String:"screenWidth"}:uint64(1),
   protobuf.Tag{Number:26}:[]protobuf.Message{
      protobuf.Message{
         protobuf.Tag{Number:1}:"android.hardware.touchscreen",
      },
      protobuf.Message{
         protobuf.Tag{Number:1}:"android.hardware.screen.portrait",
      },
      protobuf.Message{
         protobuf.Tag{Number:1}:"android.hardware.wifi",
      },
   },
}

////////////////////////////////////////////////////////////////////////////////

var checkinBody = protobuf.Message{
   protobuf.Tag{Number:4, String: "checkin"}:protobuf.Message{},
   protobuf.Tag{Number:14, String: "version"}:uint64(3),
}
