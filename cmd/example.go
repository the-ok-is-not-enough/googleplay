package main

import (
   "fmt"
   "time"
   gp "github.com/89z/googleplay"
)

const (
   app = "com.google.android.youtube"
   ver = 1524094400
)

func main() {
   tok, err := gp.NewToken("EMAIL", "PASSWORD")
   if err != nil {
      panic(err)
   }
   auth, err := tok.Auth()
   if err != nil {
      panic(err)
   }
   dev, err := gp.NewDevice(gp.DefaultCheckin)
   if err != nil {
      panic(err)
   }
   auth.Upload(dev, gp.DefaultConfig)
   time.Sleep(gp.Sleep)
   auth.Purchase(dev, app)
   del, err := auth.Delivery(dev, app, ver)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%+v\n", del)
}
