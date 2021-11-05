package main

import (
   "fmt"
   "flag"
   "github.com/89z/googleplay"
   "os"
   "path/filepath"
   "time"
)

func checkin() (string, error) {
   auth, cache, err := getAuth()
   if err != nil {
      return "", err
   }
   check, err := googleplay.NewCheckinRequest().Post()
   if err != nil {
      return "", err
   }
   if err := auth.Upload(check.String(), googleplay.NewDevice()); err != nil {
      return "", err
   }
   cache = filepath.Join(cache, "/googleplay/checkin.json")
   write, err := os.Create(cache)
   if err != nil {
      return "", err
   }
   defer write.Close()
   if err := check.Encode(write); err != nil {
      return "", err
   }
   return cache, nil
}

func details(app string) (*googleplay.AppDetails, error) {
   auth, cache, err := getAuth()
   if err != nil {
      return nil, err
   }
   check := new(googleplay.Checkin)
   read, err := os.Open(cache + "/googleplay/checkin.json")
   if err != nil {
      return nil, err
   }
   defer read.Close()
   if err := check.Decode(read); err != nil {
      return nil, err
   }
   return auth.Details(check.String(), app)
}

func getAuth() (*googleplay.Auth, string, error) {
   tok := new(googleplay.Token)
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, "", err
   }
   read, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      return nil, "", err
   }
   defer read.Close()
   if err := tok.Decode(read); err != nil {
      return nil, "", err
   }
   auth, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return auth, cache, nil
}

func token(email, password string) (string, error) {
   tok, err := googleplay.NewToken(email, password)
   if err != nil {
      return "", err
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      return "", err
   }
   cache = filepath.Join(cache, "googleplay")
   os.Mkdir(cache, os.ModeDir)
   cache = filepath.Join(cache, "token.json")
   file, err := os.Create(cache)
   if err != nil {
      return "", err
   }
   defer file.Close()
   if err := tok.Encode(file); err != nil {
      return "", err
   }
   return cache, nil
}

func main() {
   var (
      check bool
      app, email, password string
   )
   flag.BoolVar(&check, "c", false, "checkin")
   flag.StringVar(&app, "a", "", "get app details")
   flag.StringVar(&email, "e", "", "email")
   flag.StringVar(&password, "p", "", "password")
   flag.Parse()
   switch {
   case app != "":
      detail, err := details(app)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%+v\n", detail)
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
   default:
      fmt.Println("googleplay [flags]")
      flag.PrintDefaults()
   }
}
