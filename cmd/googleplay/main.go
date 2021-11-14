package main

import (
   "fmt"
   "flag"
)

func main() {
   var (
      app, email, pass string
      dev, purch bool
      ver int
   )
   flag.StringVar(&app, "a", "", "app")
   flag.BoolVar(&dev, "d", false, "create device")
   flag.StringVar(&email, "e", "", "email")
   flag.StringVar(&pass, "p", "", "password")
   flag.BoolVar(
      &purch, "purchase", false,
      "Purchase app. Only needs to be done once per Google account.",
   )
   flag.IntVar(&ver, "v", 0, "version")
   flag.Parse()
   switch {
   case email != "":
      cache, err := token(email, pass)
      if err != nil {
         panic(err)
      }
      fmt.Println("Create", cache)
   case dev:
      cache, err := device()
      if err != nil {
         panic(err)
      }
      fmt.Println("Create", cache)
   case app != "" && !purch && ver == 0:
      det, err := details(app)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%+v\n", det.DocV2)
   case app != "" && purch:
      err := purchase(app)
      if err != nil {
         panic(err)
      }
   case app != "" && ver != 0:
      err := delivery(app, ver)
      if err != nil {
         panic(err)
      }
   default:
      fmt.Println("googleplay [flags]")
      flag.PrintDefaults()
   }
}
