package main

import (
   "flag"
   "fmt"
   "github.com/89z/googleplay"
   "strings"
)

func main() {
   // a
   var app string
   flag.StringVar(&app, "a", "", "app")
   // d
   var device bool
   flag.BoolVar(&device, "d", false, "create device")
   // e
   var email string
   flag.StringVar(&email, "e", "", "email")
   // p
   var password string
   flag.StringVar(&password, "p", "", "password")
   // purchase
   var (
      pur strings.Builder
      purchase bool
   )
   pur.WriteString("Purchase app.")
   pur.WriteString(" Only needs to be done once per Google account.")
   flag.BoolVar(&purchase, "purchase", false, pur.String())
   // s
   var single bool
   flag.BoolVar(&single, "s", false, "single APK")
   // v
   var version int64
   flag.Int64Var(&version, "v", 0, "version")
   // verbose
   var verbose bool
   flag.BoolVar(&verbose, "verbose", false, "dump requests")
   flag.Parse()
   if verbose {
      googleplay.LogLevel = 1
   }
   switch {
   case email != "":
      err := doToken(email, password)
      if err != nil {
         panic(err)
      }
   case device:
      cache, err := doDevice()
      if err != nil {
         panic(err)
      }
      fmt.Println("Create", cache)
   case app != "" && !purchase && version == 0:
      res, err := doDetails(app)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%+v\n", res)
   case app != "" && purchase:
      err := doPurchase(app)
      if err != nil {
         panic(err)
      }
   case app != "" && version != 0:
      err := doDelivery(app, version, single)
      if err != nil {
         panic(err)
      }
   default:
      fmt.Println("googleplay [flags]")
      flag.PrintDefaults()
   }
}
