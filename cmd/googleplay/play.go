package main

import (
   "github.com/89z/googleplay"
   "os"
   "path/filepath"
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

func delivery(app string, ver int) (*googleplay.Delivery, error) {
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
   return auth.Delivery(check.String(), app, ver)
}

func details(app string) (*googleplay.Details, error) {
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
