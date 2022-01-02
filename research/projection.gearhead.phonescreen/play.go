package play

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

const auth = "ya29.a0ARrdaM_KG6LR-nzfsM8vHVC1Q3WEUrPD9fXkgfngrB5vVGaudcbqSn_wKTpJkCyJ5nrajNJBqDbRbcp8-zcGNeVD-YA89oyF_i93tLkFAnvS49ejme82cjU37lkjxhi4q8HYonp6PYm2hMZmSv4ohZHVczEXeQdNSWdZzorQaRiJt0SAiTJpwJrogWx036ts0mT8AF5OobK2gYUMaaZC8Wbdj5WT3irIyySL30qf9UmfO43S3TCpvyupjoNg7qETOohwZu2UD2vvPRhj3lB5pjWw0PZD3WdwWUhc_4PulcJ5BJqUXvs"

var body1 = protobuf.Message{
   protobuf.Tag{Number:4}:protobuf.Message{
      protobuf.Tag{Number:1}:protobuf.Message{
         protobuf.Tag{Number:10}:uint64(29),
      },
   },
   protobuf.Tag{Number:14}:uint64(3),
   protobuf.Tag{Number:18}: body0,
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

func checkin() (uint64, error) {
   var req0 = &http.Request{
      Method:"POST",
      URL:&url.URL{Scheme:"https",
         Host:"android.clients.google.com",
         Path:"/checkin", 
      },
      Header: http.Header{
         "Content-Type": {"application/x-protobuffer"},
      },
      Body: io.NopCloser(body1.Encode()),
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

func details(app string) (uint64, error) {
   id, err := checkin()
   if err != nil {
      return 0, err
   }
   sID := strconv.FormatUint(id, 16)
   fmt.Println(sID)
   var req5 = &http.Request{Method:"GET", URL:&url.URL{Scheme:"https",
      Host:"android.clients.google.com",
      Path:"/fdfe/details", RawQuery:"doc=" + app,
      },
      Header:http.Header{
         "Authorization":[]string{"Bearer " + auth},
         "X-Dfe-Device-Id":[]string{sID},
      },
   }
   res, err := new(http.Transport).RoundTrip(req5)
   if err != nil {
      return 0, err
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(1,2,4,13,1,3), nil
}
