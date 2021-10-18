package main

import (
   "bytes"
   "github.com/89z/mech"
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

func main() {
   val, err := url.ParseQuery(secret)
   if err != nil {
      panic(err)
   }
   ac2 := Ac2dm{val}
   mech.Verbose(true)
   auth, err := ac2.OAuth2()
   if err != nil {
      panic(err)
   }
   data, err := auth.Details("38B5418D8683ADBB", "com.google.android.youtube")
   if err != nil {
      panic(err)
   }
   f, err := os.Create("details.txt")
   if err != nil {
      panic(err)
   }
   defer f.Close()
   f.Write(data)
}

const Origin = "https://android.clients.google.com"

type Ac2dm struct {
   url.Values
}

func (a Ac2dm) OAuth2() (*OAuth2, error) {
   val := url.Values{
      "Token": {
         a.Get("Token"),
      },
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }
   req, err := http.NewRequest(
      "POST", Origin + "/auth", strings.NewReader(val.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Content-Type": {"application/x-www-form-urlencoded"},
   }
   res, err := mech.RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return &OAuth2{
      ParseQuery(query),
   }, nil
}

type OAuth2 struct {
   url.Values
}

func (o OAuth2) Details(device, app string) ([]byte, error) {
   req, err := http.NewRequest("GET", Origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {
         "Bearer " + o.Get("Auth"),
      },
      "X-DFE-Device-Id": {device},
   }
   val := url.Values{
      "doc": {app},
   }
   req.URL.RawQuery = val.Encode()
   res, err := mech.RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}

func ParseQuery(query []byte) url.Values {
   res := make(url.Values)
   for _, pair := range bytes.Split(query, []byte{'\n'}) {
      nv := bytes.SplitN(pair, []byte{'='}, 2)
      res.Add(string(nv[0]), string(nv[1]))
   }
   return res
}
