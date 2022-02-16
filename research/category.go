package main

import (
   "net/http"
   "net/url"
   "encoding/json"
   "github.com/89z/format/protobuf"
   "os"
)

func main() {
   var req http.Request
   req.Header = make(http.Header)
   req.URL = new(url.URL)
   req.URL.Host = "android.clients.google.com"
   val := make(url.Values)
   req.URL.Scheme = "https"
   val["cat"] = []string{"VIDEO_PLAYERS"}
   req.Header["Authorization"] = []string{bearer}
   req.Header["X-Dfe-Device-Id"] = []string{"3958f0cc913d5165"}
   req.Header["User-Agent"] = []string{"Android-Finsky (sdk=9,versionCode=99999999)"}
   // You can change this to "4", but it will fail with newer versionCode:
   val["c"] = []string{"3"}
   // You can also use "/fdfe/browse" or "/fdfe/homeV2", but they do Prefetch,
   // and seem to ignore the X-DFE-No-Prefetch:true header:
   req.URL.Path = "/fdfe/getHomeStream"
   req.URL.RawQuery = val.Encode()
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      panic(err)
   }
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.SetEscapeHTML(false)
   enc.Encode(mes)
}
