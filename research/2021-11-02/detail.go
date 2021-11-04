package main

import (
   "bytes"
   "fmt"
   "github.com/89z/googleplay"
   "net/url"
   "os"
)

func main() {
   if len(os.Args) != 2 {
      fmt.Println("play [device]")
      return
   }
   device := os.Args[1]
   ac2dmToken, ok := os.LookupEnv("ac2dmToken")
   if ! ok {
      panic("ac2dmToken")
   }
   auth := googleplay.Auth{
      url.Values{
         "Auth": {ac2dmToken},
      },
   }
   det, err := auth.Details(device, "com.google.android.youtube")
   if err != nil {
      panic(err)
   }
   vers := []string{"16.", "16.4", "16.43.", "16.43.3", "16.43.34"}
   for _, ver := range vers {
      if bytes.Contains(det, []byte(ver)) {
         fmt.Println("pass", ver)
      } else {
         fmt.Println("fail", ver)
      }
   }
}
