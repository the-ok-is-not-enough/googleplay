package main

import (
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
)

/*
1 message payload
3 message reviewResponse
1 message getResponse
1 message	
0	
8 string "I have a lot of experience of canceling travel in last 30min.
*/
func main() {
   var req http.Request
   req.Header = make(http.Header)
   req.URL = new(url.URL)
   req.URL.Host = "android.clients.google.com"
   req.URL.Path = "/fdfe/rev"
   val := make(url.Values)
   val["doc"] = []string{"com.comuto"}
   req.URL.RawQuery = val.Encode()
   req.URL.Scheme = "https"
   req.Header["X-Dfe-Device-Id"] = []string{"3958f0cc913d5165"}
   req.Header["Authorization"] = []string{bearer}
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err := httputil.DumpResponse(res, false)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
}
