package main

import (
   "io"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
)

func main() {
   var req http.Request
   req.Body = io.NopCloser(body)
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

var body = strings.NewReader("\n\x16com.google.android.gms\x10\xd6\xf1\xabi\x1a\x1e\n\x0eGames.optional\x1a\f220920000000 \x01 \x02 \x03 \x05(\x02(\x018\x01")
