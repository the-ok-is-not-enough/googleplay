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

func doDelivery(head *gp.Header, output, app string, ver uint64) error {
   del, err := head.Delivery(app, ver)
   if err != nil {
      return err
   }
   dst := filename(output, app, "", ver)
   if err := download(del.DownloadURL, dst); err != nil {
      return err
   }
   for _, split := range del.SplitDeliveryData {
      dst := filename(output, app, split.ID, ver)
      err := download(split.DownloadURL, dst)
      if err != nil {
         return err
      }
   }
   return nil
}

func doDevice() error {
   dev, err := gp.DefaultConfig.Checkin()
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

func filename(output, app, id string, ver uint64) string {
   var buf []byte
   if output != "" {
      buf = append(buf, output...)
      buf = append(buf, '/')
   }
   buf = append(buf, app...)
   buf = append(buf, '-')
   if id != "" {
      buf = append(buf, id...)
      buf = append(buf, '-')
   }
   buf = strconv.AppendUint(buf, ver, 10)
   buf = append(buf, ".apk"...)
   return string(buf)
}

func newHeader(single bool) (*gp.Header, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, err
   }
   tok, err := gp.OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   dev, err := gp.OpenDevice(cache, "googleplay/device.json")
   if err != nil {
      return nil, err
   }
   if single {
      return tok.SingleAPK(dev)
   }
   return tok.Header(dev)
}
