package main

import (
   "fmt"
   "os"
   "net/http"
   "path/filepath"
   "time"
   gp "github.com/89z/googleplay"
)

func delivery(app string, ver int) error {
   auth, cache, err := getAuth()
   if err != nil {
      return err
   }
   dev := new(gp.Device)
   rd, err := os.Open(cache + "/googleplay/device.json")
   if err != nil {
      return err
   }
   defer rd.Close()
   if err := dev.Decode(rd); err != nil {
      return err
   }
   del, err := auth.Delivery(dev, app, ver)
   if err != nil {
      return err
   }
   fmt.Println("GET", del.AppDeliveryData.DownloadURL)
   res, err := http.Get(del.AppDeliveryData.DownloadURL)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   name := fmt.Sprintf("%v-%v.apk", app, ver)
   wr, err := os.Create(name)
   if err != nil {
      return err
   }
   defer wr.Close()
   if _, err := wr.ReadFrom(res.Body); err != nil {
      return err
   }
   return nil
}

func purchase(app string) error {
   auth, cache, err := getAuth()
   if err != nil {
      return err
   }
   dev := new(gp.Device)
   rd, err := os.Open(cache + "/googleplay/device.json")
   if err != nil {
      return err
   }
   defer rd.Close()
   if err := dev.Decode(rd); err != nil {
      return err
   }
   return auth.Purchase(dev, app)
}

func details(app string) (*gp.Details, error) {
   auth, cache, err := getAuth()
   if err != nil {
      return nil, err
   }
   dev := new(gp.Device)
   rd, err := os.Open(cache + "/googleplay/device.json")
   if err != nil {
      return nil, err
   }
   defer rd.Close()
   if err := dev.Decode(rd); err != nil {
      return nil, err
   }
   return auth.Details(dev, app)
}

func device() (string, error) {
   auth, cache, err := getAuth()
   if err != nil {
      return "", err
   }
   dev, err := gp.NewDevice(gp.DefaultCheckin)
   if err != nil {
      return "", err
   }
   if err := auth.Upload(dev, gp.DefaultConfig); err != nil {
      return "", err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   cache = filepath.Join(cache, "/googleplay/device.json")
   write, err := os.Create(cache)
   if err != nil {
      return "", err
   }
   defer write.Close()
   if err := dev.Encode(write); err != nil {
      return "", err
   }
   return cache, nil
}

func getAuth() (*gp.Auth, string, error) {
   tok := new(gp.Token)
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, "", err
   }
   rd, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      return nil, "", err
   }
   defer rd.Close()
   if err := tok.Decode(rd); err != nil {
      return nil, "", err
   }
   auth, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return auth, cache, nil
}

func token(email, password string) (string, error) {
   tok, err := gp.NewToken(email, password)
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
