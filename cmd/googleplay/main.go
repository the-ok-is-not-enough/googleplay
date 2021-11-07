package main

import (
   "fmt"
   "flag"
)

func main() {
   var (
      app, email, password string
      dev bool
      version int
   )
   flag.BoolVar(&dev, "d", false, "device")
   flag.IntVar(&version, "v", 0, "version")
   flag.StringVar(&app, "a", "", "get app details")
   flag.StringVar(&email, "e", "", "email")
   flag.StringVar(&password, "p", "", "password")
   flag.Parse()
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
   case app != "" && version == 0:
      det, err := details(app)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%+v\n", det)
   case app != "" && version != 0:
      del, err := delivery(app, version)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%+v\n", del)
   default:
      fmt.Println("googleplay [flags]")
      flag.PrintDefaults()
   }
}
