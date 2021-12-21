package main

import (
   "flag"
   "fmt"
   "github.com/89z/googleplay"
   "net/http"
   "os"
)

func download(addr, id, app string, ver int) error {
   fmt.Println("GET", addr)
   res, err := http.Get(addr)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   var name string
   if id != "" {
      name = fmt.Sprintf("%v-%v-%v.apk", app, id, ver)
   } else {
      name = fmt.Sprintf("%v-%v.apk", app, ver)
   }
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := file.ReadFrom(res.Body); err != nil {
      return err
   }
   return nil
}

func main() {
   var (
      app, email, pass string
      dev, purch, verbose bool
      version int
   )
   flag.StringVar(&app, "a", "", "app")
   flag.BoolVar(&dev, "d", false, "create device")
   flag.StringVar(&email, "e", "", "email")
   flag.StringVar(&pass, "p", "", "password")
   flag.BoolVar(
      &purch, "purchase", false,
      "Purchase app. Only needs to be done once per Google account.",
   )
   flag.IntVar(&version, "v", 0, "version")
   flag.BoolVar(&verbose, "verbose", false, "dump requests")
   flag.Parse()
   googleplay.Verbose = verbose
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
      err := delivery(app, version)
      if err != nil {
         panic(err)
      }
   default:
      fmt.Println("googleplay [flags]")
      flag.PrintDefaults()
   }
}
