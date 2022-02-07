package main

import (
   "fmt"
   "github.com/89z/format"
   "net/http"
   "os"
   "path/filepath"
   "strconv"
   "time"
   gp "github.com/89z/googleplay"
)

func doDelivery(app string, ver int64, single bool) error {
   auth, cache, err := getAuth()
   if err != nil {
      return err
   }
   dev, err := gp.OpenDevice(cache + "/googleplay/device.json")
   if err != nil {
      return err
   }
   head := auth.AppBundle(dev)
   if single {
      head = auth.AndroidPackage(dev)
   }
   del, err := head.Delivery(app, ver)
   if err != nil {
      return err
   }
   dst := filename(app, "", ver)
   if err := download(del.DownloadURL, dst); err != nil {
      return err
   }
   for _, split := range del.SplitDeliveryData {
      dst := filename(app, split.ID, ver)
      err := download(split.DownloadURL, dst)
      if err != nil {
         return err
      }
   }
   return nil
}

func doDetails(app string) (*gp.Details, error) {
   auth, cache, err := getAuth()
   if err != nil {
      return nil, err
   }
   dev, err := gp.OpenDevice(cache + "/googleplay/device.json")
   if err != nil {
      return nil, err
   }
   return auth.AppBundle(dev).Details(app)
}

func doDevice() (string, error) {
   dev, err := gp.NewDevice(gp.DefaultConfig)
   if err != nil {
      return "", err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   cache, err := os.UserCacheDir()
   if err != nil {
      return "", err
   }
   cache = filepath.Join(cache, "/googleplay/device.json")
   if err := dev.Create(cache); err != nil {
      return "", err
   }
   return cache, nil
}

func doPurchase(app string) error {
   auth, cache, err := getAuth()
   if err != nil {
      return err
   }
   dev, err := gp.OpenDevice(cache + "/googleplay/device.json")
   if err != nil {
      return err
   }
   return auth.AppBundle(dev).Purchase(app)
}

func doToken(email, password string) error {
   tok, err := gp.NewToken(email, password)
   if err != nil {
      return err
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.Join(cache, "googleplay")
   os.Mkdir(cache, os.ModePerm)
   cache = filepath.Join(cache, "token.json")
   fmt.Println("Create", cache)
   return tok.Create(cache)
}

func download(src, dst string) error {
   fmt.Println("GET", src)
   res, err := http.Get(src)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   file, err := os.Create(dst)
   if err != nil {
      return err
   }
   defer file.Close()
   pro := format.NewProgress(res)
   if _, err := file.ReadFrom(pro); err != nil {
      return err
   }
   return nil
}

func filename(app, id string, ver int64) string {
   var buf []byte
   buf = append(buf, app...)
   buf = append(buf, '-')
   if id != "" {
      buf = append(buf, id...)
      buf = append(buf, '-')
   }
   buf = strconv.AppendInt(buf, ver, 10)
   buf = append(buf, ".apk"...)
   return string(buf)
}

func getAuth() (*gp.Auth, string, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, "", err
   }
   tok, err := gp.OpenToken(cache + "/googleplay/token.json")
   if err != nil {
      return nil, "", err
   }
   auth, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return auth, cache, nil
}
