package main

import (
   "fmt"
   "os"
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
   del, err := auth.DeliveryResponse(dev, app, ver)
   if err != nil {
      return err
   }
   data := del.AppDeliveryData()
   if err := download(data.DownloadURL(), "", app, ver); err != nil {
      return err
   }
   for _, split := range data.SplitDeliveryData() {
      err := download(split.DownloadURL(), split.ID(), app, ver)
      if err != nil {
         return err
      }
   }
   return nil
}

func detailsResponse(app string) (*gp.DetailsResponse, error) {
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
   return auth.DetailsResponse(dev, app)
}

func device() (string, error) {
   auth, cache, err := getAuth()
   if err != nil {
      return "", err
   }
   dev, err := gp.NewDevice()
   if err != nil {
      return "", err
   }
   con := gp.NewConfig()
   if err := auth.Upload(dev, con); err != nil {
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
   os.Mkdir(cache, os.ModePerm)
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
