package main

import (
   "flag"
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
   // email
   var email string
   flag.StringVar(&email, "email", "", "your email")
   // p
   var id int64
   flag.Int64Var(&id, "p", 0, googleplay.Platforms.String())
   // password
   var password string
   flag.StringVar(&password, "password", "", "your password")
   // purchase
   var (
      buf strings.Builder
      purchase bool
   )
   buf.WriteString("Purchase app. ")
   buf.WriteString("Only needs to be done once per Google account.")
   flag.BoolVar(&purchase, "purchase", false, buf.String())
   // s
   var single bool
   flag.BoolVar(&single, "s", false, "single APK")
   // v
   var version uint64
   flag.Uint64Var(&version, "v", 0, "version")
   // verbose
   var verbose bool
   flag.BoolVar(&verbose, "verbose", false, "dump requests")
   flag.Parse()
   if verbose {
      googleplay.LogLevel = 1
   }
   if email != "" {
      err := doToken(email, password)
      if err != nil {
         panic(err)
      }
   } else {
      platform := googleplay.Platforms[id]
      if device {
         err := doDevice(platform)
         if err != nil {
            panic(err)
         }
      } else if app != "" {
         head, err := doHeader(platform, single)
         if err != nil {
            panic(err)
         }
         if purchase {
            err := head.Purchase(app)
            if err != nil {
               panic(err)
            }
         } else if version >= 1 {
            err := doDelivery(head, app, version)
            if err != nil {
               panic(err)
            }
         } else {
            err := doDetails(head, app)
            if err != nil {
               panic(err)
            }
         }
      } else {
         flag.Usage()
      }
   }
}
