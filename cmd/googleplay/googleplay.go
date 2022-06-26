package main

import (
   "fmt"
   "github.com/89z/format"
   "io"
   "time"
   gp "github.com/89z/googleplay"
)

func do_header(dir, platform string, single bool) (*gp.Header, error) {
   auth, err := gp.Open_Auth(dir + "/auth.txt")
   if err != nil {
      return nil, err
   }
   if err := auth.Exchange(); err != nil {
      return nil, err
   }
   device, err := gp.Open_Device(dir + "/" + platform + ".bin")
   if err != nil {
      return nil, err
   }
   id, err := device.ID()
   if err != nil {
      return nil, err
   }
   head := auth.Header(id, single)
   return &head, nil
}

func do_details(head *gp.Header, app string) (string, error) {
   detail, err := head.Details(app)
   if err != nil {
      return "", err
   }
   text, err := detail.MarshalText()
   if err != nil {
      return "", err
   }
   return string(text), nil
}

func do_delivery(head *gp.Header, app string, ver uint64) error {
   download := func(addr, name string) error {
      res, err := gp.Client.Redirect(nil).Get(addr)
      if err != nil {
         return err
      }
      defer res.Body.Close()
      file, err := format.Create(name)
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
   for _, split := range del.Split_Data() {
      addr, err := split.Download_URL()
      if err != nil {
         return err
      }
      id, err := split.ID()
      if err != nil {
         return err
      }
      name := gp.File{app, ver}.APK(id)
      if err := download(addr, name); err != nil {
         return err
      }
   }
   for _, file := range del.Additional_File() {
      addr, err := file.Download_URL()
      if err != nil {
         return err
      }
      typ, err := file.File_Type()
      if err != nil {
         return err
      }
      name := gp.File{app, ver}.OBB(typ)
      if err := download(addr, name); err != nil {
         return err
      }
   }
   addr, err := del.Download_URL()
   if err != nil {
      return err
   }
   name := gp.File{app, ver}.APK("")
   return download(addr, name)
}

func do_auth(dir, email, password string) error {
   auth, err := gp.New_Auth(email, password)
   if err != nil {
      return err
   }
   return auth.Create(dir + "/auth.txt")
}

func do_device(dir, platform string) error {
   device, err := gp.Phone.Checkin(platform)
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   return device.Create(dir + "/" + platform + ".bin")
}
