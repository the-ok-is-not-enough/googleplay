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
   responseWrapper, err := protobuf.Unmarshal(buf)
   if err != nil {
      panic(err)
   }
   child := responseWrapper.Get(3, "preFetch").
      Get(2, "response").
      Get(1, "payload").
      Get(1, "listResponse").
      Get(2, "doc").
      GetMessages(11, "child")
   for _, doc := range child {
      for _, doc := range doc.GetMessages(11, "child") {
         id := doc.GetString(1, "docid")
         fmt.Println(id)
      }
   }
}
