package main

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "net/http"
)

func main() {
   src := protobuf.Message{
      {3, "digest"}: "", {4, "checkin"}: protobuf.Message{},
   }
   res, err := http.Post(
      "http://android.clients.google.com/checkin",
      "application/x-protobuffer", src.Encode(),
   )
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
