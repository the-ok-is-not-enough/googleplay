package main

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "net/http"
)

func main() {
   src := protobuf.Message{
      {3, ""}: "", {4, ""}: "",
   }
   req, err := http.NewRequest(
      "POST", "http://android.clients.google.com/checkin", src.Encode(),
   )
   if err != nil {
      panic(err)
   }
   req.Header.Set("Content-Type", "application/x-protobuffer")
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   dst, err := protobuf.Decode(res.Body)
   if err != nil {
      panic(err)
   }
   fmt.Println(dst)
}
