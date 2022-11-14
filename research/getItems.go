package main

import (
   "bytes"
   "fmt"
   "github.com/89z/rosso/protobuf"
   "io"
   "net/http"
   "net/url"
   "encoding/base64"
)

func main() {
   var req http.Request
   req_body := protobuf.Message{
      2:protobuf.Message{
         1:protobuf.Message{
            1: protobuf.Message{
               1: protobuf.String("com.balysv.loop"),
            },
         },
      },
   }.Marshal()
   req.Body = io.NopCloser(bytes.NewReader(req_body))
   req.Header = make(http.Header)
   req.Method = "POST"
   req.URL = new(url.URL)
   req.URL.Host = "play-fe.googleapis.com"
   req.URL.Path = "/fdfe/getItems"
   req.URL.Scheme = "https"
   req.Header["X-Dfe-Device-Id"] = []string{"374a9c9111827216"}
   req.Header["X-Dfe-Item-Field-Mask"] = []string{field_mask()}
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      panic(res.Status)
   }
   res_body, err := io.ReadAll(res.Body)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%q\n", res_body)
}

func field_mask() string {
   mask := protobuf.Message{
      3: protobuf.Bytes{0x20, 0xa1, 0xf, 0x9e, 0x5},
      4: protobuf.Bytes{0x0, 0x92, 0xb3, 0x7, 0xb2, 0xff, 0x0, 0x80},
   }.Marshal()
   return base64.StdEncoding.EncodeToString(mask)
}
