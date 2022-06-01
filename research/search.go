package search

import (
   "encoding/json"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "github.com/89z/googleplay"
   "net/http"
   "net/url"
   "os"
)

var LogLevel format.LogLevel

func search(query string) error {
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   token, err := googleplay.OpenToken(home, "googleplay/token.json")
   if err != nil {
      return err
   }
   device, err := googleplay.OpenDevice(home, "googleplay/x86.json")
   if err != nil {
      return err
   }
   head, err := token.Header(device.AndroidID, false)
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/search", nil,
   )
   head.SetAuth(req.Header)
   head.SetDevice(req.Header)
   req.URL.RawQuery = url.Values{
      "c": {"3"},
      "q": {query},
   }.Encode()
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return err
   }
   file, err := os.Create(query + ".json")
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(responseWrapper.Get(1).Get(5))
}
