package main

import (
   "github.com/89z/googleplay"
   "os"
)

func main() {
   auth, err := googleplay.Open_Auth("auth.txt")
   if err != nil {
      panic(err)
   }
   auth.Exchange()
   device, err := googleplay.Open_Device("x86.bin")
   if err != nil {
      panic(err)
   }
   id, err := device.ID()
   if err != nil {
      panic(err)
   }
   detail, err := auth.Header(id, false).Details("com.app.xt")
   if err != nil {
      panic(err)
   }
   text, err := detail.MarshalText()
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(text)
}
