package main

import (
   "fmt"
   "github.com/89z/parse/protobuf"
   "net/http"
   "net/http/httputil"
   "strings"
)

const origin = "https://android.clients.google.com"

/*
[
   {Number:1 Type:0 Value:1}
   {Number:3 Type:0 Value:1635970217779}
   {Number:7 Type:1 Value:4316223298629298718}
   {Number:8 Type:1 Value:1924261969879187925}
   {Number:11 Type:2 Value:MxhuIxnK6ZUeQHjbYkgYHaqqcMd2zhc}
   {Number:12 Type:2 Value:ABFEt1UYsZhI54CRZ1aNncOMh1sdgKbIVmu8tdAUkMnYZTFEOdl040v1SOTZ-g7l_2Vrq1lqa1tMGd9cwGOnpehClE6xlSOTR50xKhXDuU91Q4MPYTJCA5CqwvYWbhS3MAwLflmTZw1neYPcccLjWjySdJIMV2awmwLWUBdCJphgva24dwuMJ7h0T8r_hSWMeK5Yh6blOUIoU0DEQ0Qe80luRx3cO2yBBeX4uYJalGdYqd4b5NrHpV1J79pPuyeFiB2apn6mdcVQBtljAWa6mYrWxWwekWoGcN-RsTigpiLAjYmrhWZtjdE}
]
*/
var buf = []byte("\b\x01\x18\xb3\uec7b\xce/9\x1ej\x10\xc68N\xe6;A\xd5Aq\x01\x11Z\xb4\x1aZ\x1fMxhuIxnK6ZUeQHjbYkgYHaqqcMd2zhcb\xb7\x02ABFEt1UYsZhI54CRZ1aNncOMh1sdgKbIVmu8tdAUkMnYZTFEOdl040v1SOTZ-g7l_2Vrq1lqa1tMGd9cwGOnpehClE6xlSOTR50xKhXDuU91Q4MPYTJCA5CqwvYWbhS3MAwLflmTZw1neYPcccLjWjySdJIMV2awmwLWUBdCJphgva24dwuMJ7h0T8r_hSWMeK5Yh6blOUIoU0DEQ0Qe80luRx3cO2yBBeX4uYJalGdYqd4b5NrHpV1J79pPuyeFiB2apn6mdcVQBtljAWa6mYrWxWwekWoGcN-RsTigpiLAjYmrhWZtjdE")

func checkin() ([]byte, error) {
   req, err := http.NewRequest(
      "POST", origin + "/checkin", strings.NewReader("\"\x00p\x03"),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-protobuf")
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return httputil.DumpResponse(res, true)
}

func main() {
   fields := protobuf.Parse(buf)
   fmt.Printf("%+v\n", fields)
}
