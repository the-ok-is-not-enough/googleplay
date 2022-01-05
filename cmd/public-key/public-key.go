package main

import (
   "bytes"
   "fmt"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
)

func main() {
   src := protobuf.Message{
      {3, "digest"}: "", {4, "checkin"}: protobuf.Message{},
   }
   res, err := http.Post(
      "http://android.clients.google.com/checkin",
      "application/x-protobuffer", bytes.NewReader(src.Marshal()),
   )
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      panic(err)
   }
   dst, err := protobuf.Unmarshal(buf)
   if err != nil {
      panic(err)
   }
   fmt.Println(dst)
}
