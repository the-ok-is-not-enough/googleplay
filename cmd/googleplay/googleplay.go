package main

import (
   "fmt"
   "github.com/89z/format"
   "io"
   "net/http"
   "os"
   "time"
   gp "github.com/89z/googleplay"
)

func doDelivery(head *gp.Header, app string, ver uint64) error {
   download := func(addr gp.String, name string) error {
      fmt.Println("GET", addr)
      res, err := http.Get(string(addr))
      if err != nil {
         return err
      }
      defer res.Body.Close()
      file, err := os.Create(name)
      if err != nil {
         return err
      }
      defer file.Close()
      pro := format.ProgressBytes(file, res.ContentLength)
      if _, err := io.Copy(pro, res.Body); err != nil {
         return err
      }
      return nil
   }
   del, err := head.Delivery(app, ver)
   if err != nil {
      return err
   }
   for _, split := range del.SplitDeliveryData {
      err := download(split.DownloadURL, del.Split(split.ID))
      if err != nil {
         return err
      }
   }
   for _, file := range del.AdditionalFile {
      err := download(file.DownloadURL, del.Additional(file.FileType))
      if err != nil {
         return err
      }
   }
   return download(del.DownloadURL, del.Download())
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

type native struct {
   path string
   platform gp.String
}

func newNative(armeabi, arm64 bool) native {
   if armeabi {
      return native{"googleplay/armeabi.json", gp.Armeabi}
   }
   if arm64 {
      return native{"googleplay/arm64.json", gp.Arm64}
   }
   return native{"googleplay/x86.json", gp.X86}
}

func (n native) device() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   device, err := gp.Phone.Checkin(n.platform)
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   return device.Create(cache, n.path)
}

func (n native) header(single bool) (*gp.Header, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, err
   }
   token, err := gp.OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   device, err := gp.OpenDevice(cache, n.path)
   if err != nil {
      return nil, err
   }
   if single {
      return token.SingleAPK(device)
   }
   return token.Header(device)
}
