package main

import (
   "encoding/json"
   "fmt"
   "github.com/89z/parse/protobuf"
   "net/http"
   "net/url"
   "time"
)

type app struct {
   numDownloads string
   id string
}

var apps = []app{
   {"10+", "com.elitegames.hillcarracing"},
   {"100+", "com.ANOOGAMES.WordsSearchPremium"},
   {"500+", "com.gameo2.CatColoring"},
   {"1,000+", "com.fbig.crosswords"},
   {"10,000+", "org.friends.dragonli"},
   {"50,000+", "com.techorus.HyperSodaGeyser"},
   {"100,000+", "com.ketchapp.jabbycat"},
   {"1,000,000+", "jp.co.ponos.nyanko_odorobo"},
   {"10,000,000+", "com.peacocktv.peacockandroid"},
   {"50,000,000+", "com.reddit.frontpage"},
   {"100,000,000+", "com.discord"},
   {"1,000,000,000+", "com.netflix.mediaclient"},
   {"10,000,000,000+", "com.google.android.youtube"},
}

func main() {
   for _, app := range apps {
      req, err := http.NewRequest(
         "GET", "https://android.clients.google.com/fdfe/details", nil,
      )
      if err != nil {
         panic(err)
      }
      req.URL.RawQuery = url.Values{
         "doc": {app.id},
      }.Encode()
      req.Header = http.Header{
         "Authorization": {auth}, "X-DFE-Device-ID": {device},
      }
      res, err := new(http.Transport).RoundTrip(req)
      if err != nil {
         panic(err)
      }
      defer res.Body.Close()
      mes, err := protobuf.Decode(res.Body)
      if err != nil {
         panic(err)
      }
      buf, err := mes.MarshalJSON()
      if err != nil {
         panic(err)
      }
      var wrap responseWrapper
      if err := json.Unmarshal(buf, &wrap); err != nil {
         panic(err)
      }
      fmt.Println(app.id + ":")
      fmt.Printf("%+v\n", wrap.Payload.DetailsResponse.DocV2.Details)
      time.Sleep(time.Second)
   }
}

type responseWrapper struct {
   Payload struct {
      DetailsResponse struct {
         DocV2 struct {
            Details struct {
               AppDetails struct {
                  VersionCode int32 `json:"3"`
                  NumDownloads string `json:"13"`
                  Something int `json:"70"`
               } `json:"1"`
            } `json:"13"`
         } `json:"4"`
      } `json:"2"`
   } `json:"1"`
}
