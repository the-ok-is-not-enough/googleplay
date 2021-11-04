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
   tok := googleplay.Token{
      url.Values{
         "Token": {token},
      },
   }
   auth, err := tok.Auth()
   if err != nil {
      panic(err)
   }
   det, err := auth.Details(device, "com.google.android.youtube")
   if err != nil {
      panic(err)
   }
   fmt.Printf("%q\n", det)
   vers := []string{"16.", "16.4", "16.43.", "16.43.3", "16.43.34"}
   for _, ver := range vers {
      if bytes.Contains(det, []byte(ver)) {
         fmt.Println("pass", ver)
      } else {
         fmt.Println("fail", ver)
      }
   }
}
