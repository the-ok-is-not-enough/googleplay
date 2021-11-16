package main

import (
   "encoding/base64"
   "encoding/json"
   "github.com/89z/parse/protobuf"
   "io"
   "net/http"
   "net/url"
   "os"
)

type stream struct {
   F1 struct {
      F164 struct {
         F1 int `json:"1"`
         F5 struct {
            F1 int `json:"1"`
         } `json:"5"`
         F6 struct {
            F1 string `json:"1"`
            F2 bool `json:"2"`
            F3 int `json:"3"`
         } `json:"6"`
      } `json:"164"`
   } `json:"1"`
}

func similarApps() stream {
   var s stream
   s.F1.F164.F1 = 17
   s.F1.F164.F5.F1 = 8
   s.F1.F164.F6.F1 = "com.soundcloud.android"
   s.F1.F164.F6.F2 = true
   s.F1.F164.F6.F3 = 3
   return s
}

func youMightAlsoLike() stream {
   var s stream
   s.F1.F164.F1 = 18
   s.F1.F164.F5.F1 = 8
   s.F1.F164.F6.F1 = "com.soundcloud.android"
   s.F1.F164.F6.F2 = true
   s.F1.F164.F6.F3 = 3
   return s
}

func main() {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/getStream", nil,
   )
   if err != nil {
      panic(err)
   }
   req.Header = http.Header{
      "Authorization": {"GoogleLogin auth=" + auth},
      "X-DFE-Device-Id": {"3694a3922a861e1d"},
   }
   str := youMightAlsoLike()
   m, err := protobuf.NewEncoder(str)
   if err != nil {
      panic(err)
   }
   buf, err := m.Encode()
   if err != nil {
      panic(err)
   }
   req.URL.RawQuery = url.Values{
      "ecp": {base64.StdEncoding.EncodeToString(buf)},
   }.Encode()
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err = io.ReadAll(res.Body)
   if err != nil {
      panic(err)
   }
   dec := protobuf.NewDecoder(buf)
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.Encode(dec)
}
