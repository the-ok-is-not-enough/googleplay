package main

import (
   "fmt"
   "flag"
   "github.com/89z/googleplay"
   "time"
)

func main() {
   var (
      check bool
      app, email, password string
      version int
   )
   flag.BoolVar(&check, "c", false, "checkin")
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
   case check:
      cache, err := checkin()
      if err != nil {
         panic(err)
      }
      fmt.Printf("Sleeping %v for server to process\n", googleplay.Sleep)
      time.Sleep(googleplay.Sleep)
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


