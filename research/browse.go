package main

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "os"
)

func main() {
   buf, err := os.ReadFile("ignore.txt")
   if err != nil {
      panic(err)
   }
   mes, err := protobuf.Unmarshal(buf)
   if err != nil {
      panic(err)
   }
   fmt.Println(mes)
}
