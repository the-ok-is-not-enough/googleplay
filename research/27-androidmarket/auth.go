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
   req.Header = make(http.Header)
   req.Method = "POST"
   req.URL = new(url.URL)
   req.URL.Host = "android.clients.google.com"
   req.URL.Path = "/auth"
   req.URL.Scheme = "https"
   req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
   body := url.Values{
      "Token":[]string{token},
      "service":[]string{"androidmarket"},
   }.Encode()
   req.Body = io.NopCloser(strings.NewReader(body))
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
