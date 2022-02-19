package main

import (
   "fmt"
   "github.com/89z/format"
   "net/http"
   "os"
   "strconv"
   "time"
   gp "github.com/89z/googleplay"
)

func doDelivery(app string, ver int64, single bool) error {
   auth, cache, err := getAuth()
   if err != nil {
      return err
   }
   dev, err := gp.OpenDevice(cache, "googleplay/device.json")
   if err != nil {
      return err
   }
   head := auth.Header(dev)
   if single {
      head = auth.SingleAPK(dev)
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
   dev, err := gp.OpenDevice(cache, "googleplay/device.json")
   if err != nil {
      return nil, err
   }
   return auth.Header(dev).Details(app)
}

func doDevice() error {
   dev, err := gp.NewDevice(gp.DefaultConfig)
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   return dev.Create(cache, "googleplay/device.json")
}

func doPurchase(app string) error {
   auth, cache, err := getAuth()
   if err != nil {
      return err
   }
   dev, err := gp.OpenDevice(cache, "googleplay/device.json")
   if err != nil {
      return err
   }
   return auth.Header(dev).Purchase(app)
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
   return tok.Create(cache, "googleplay/token.json")
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
   tok, err := gp.OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, "", err
   }
   auth, err := tok.Auth()
   if err != nil {
      return nil, "", err
   }
   return auth, cache, nil
}
