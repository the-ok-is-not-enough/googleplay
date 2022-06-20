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

func do_device(dir, platform string) error {
   device, err := gp.Phone.Checkin(platform)
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   return device.Create(dir + "/" + platform + ".txt")
}

func do_token(dir, email, password string) error {
   token, err := gp.New_Token(email, password)
   if err != nil {
      return err
   }
   file, err := format.Create(dir + "/token.txt")
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := token.WriteTo(file); err != nil {
      return err
   }
   return nil
}

func do_details(head *gp.Header, app string, parse bool) error {
   detail, err := head.Details(app)
   if err != nil {
      return err
   }
   if parse {
      date, err := detail.Time()
      if err != nil {
         return err
      }
      detail.Upload_Date = date.String()
   }
   fmt.Println(detail)
   return nil
}

func do_delivery(head *gp.Header, app string, ver uint64) error {
   download := func(addr, name string) error {
      fmt.Println("GET", addr)
      res, err := http.Get(addr)
      if err != nil {
         return err
      }
      defer res.Body.Close()
      file, err := os.Create(name)
      if err != nil {
         return err
      }
      defer file.Close()
      pro := format.Progress_Bytes(file, res.ContentLength)
      if _, err := io.Copy(pro, res.Body); err != nil {
         return err
      }
      return nil
   }
   del, err := head.Delivery(app, ver)
   if err != nil {
      return err
   }
   for _, split := range del.Split_Data {
      err := download(split.Download_URL, del.Split(split.ID))
      if err != nil {
         return err
      }
   }
   for _, file := range del.Additional_File {
      err := download(file.Download_URL, del.Additional(file.File_Type))
      if err != nil {
         return err
      }
   }
   return download(del.Download_URL, del.Download())
}
func do_header(dir, platform string, single bool) (*gp.Header, error) {
   token, err := gp.Open_Token(dir + "/token.txt")
   if err != nil {
      return nil, err
   }
   device, err := gp.Open_Device(dir + "/" + platform + ".txt")
   if err != nil {
      return nil, err
   }
   id, err := device.ID()
   if err != nil {
      return nil, err
   }
   return token.Header(id, single)
}
