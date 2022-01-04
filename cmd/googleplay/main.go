package main

import (
   "flag"
   "fmt"
   "github.com/89z/googleplay"
)

func main() {
   var (
      app, email, output, password string
      dev, purch, verbose bool
      version int64
   )
   flag.StringVar(&app, "a", "", "app")
   flag.BoolVar(&dev, "d", false, "create device")
   flag.StringVar(&email, "e", "", "email")
   flag.StringVar(&output, "o", "", "output folder, must already exist")
   flag.StringVar(&password, "p", "", "password")
   flag.BoolVar(
      &purch, "purchase", false,
      "Purchase app. Only needs to be done once per Google account.",
   )
   flag.Int64Var(&version, "v", 0, "version")
   flag.BoolVar(&verbose, "verbose", false, "dump requests")
   flag.Parse()
   if verbose {
      googleplay.Log.Level = 1
   }
   switch {
   case email != "":
      cache, err := token(email, password)
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
   case app != "" && !purch && version == 0:
      res, err := details(app)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%+v\n", res)
   case app != "" && purch:
      err := purchase(app)
      if err != nil {
         panic(err)
      }
   case app != "" && version != 0:
      err := delivery(output, app, version)
      if err != nil {
         panic(err)
      }
   default:
      fmt.Println("googleplay [flags]")
      flag.PrintDefaults()
   }
}
