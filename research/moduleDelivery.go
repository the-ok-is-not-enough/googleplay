package main

import (
   "bytes"
   "github.com/89z/rosso/protobuf"
   "io"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
)

func main() {
   body := protobuf.Message{
      3:protobuf.Message{
         1: protobuf.String("Games.optional"),
      },
      1:protobuf.String("com.google.android.gms"),
      2:protobuf.Varint(220920022),
   }.Marshal()
   var req http.Request
   req.Body = io.NopCloser(bytes.NewReader(body))
   req.Header = make(http.Header)
   req.Method = "POST"
   req.URL = new(url.URL)
   req.URL.Host = "play-fe.googleapis.com"
   req.URL.Path = "/fdfe/moduleDelivery"
   req.URL.Scheme = "https"
   req.Header["X-Dfe-Device-Id"] = []string{"374a9c9111827216"}
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
}
