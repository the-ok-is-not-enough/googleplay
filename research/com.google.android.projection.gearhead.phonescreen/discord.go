package discord

import (
   "github.com/89z/format/protobuf"
   "net/http"
   "strconv"
   "strings"
)

func checkin() (uint64, error) {
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/checkin",
      strings.NewReader(checkBody),
   )
   if err != nil {
      return 0, err
   }
   req.Header = http.Header{
      "App": {"com.google.android.gms"},
      "Content-Type": {"application/x-protobuffer"},
      "User-Agent": {"GoogleAuth/1.4 sargo PQ3B.190705.003"},
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return 0, err
   }
   defer res.Body.Close()
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(7), nil
}

func details() (uint64, error) {
   id, err := checkin()
   if err != nil {
      return 0, err
   }
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/details", nil,
   )
   if err != nil {
      return 0, err
   }
   req.Header = http.Header{
      "X-Dfe-Device-Id":[]string{strconv.FormatUint(id, 16)},
      "Accept-Language":[]string{"en-GB"},
      "Authorization":[]string{"Bearer " + auth},
      "User-Agent":[]string{"Android-Finsky/15.8.23-all [0] [PR] 259261889 (api=3,versionCode=81582300,sdk=28,device=sargo,hardware=sargo,product=sargo,platformVersionRelease=9,model=Pixel 3a,buildId=PQ3B.190705.003,isWideScreen=0,supportedAbis=arm64-v8a;armeabi-v7a;armeabi)"},
      "X-Ad-Id":[]string{"LawadaMera"},
      "X-Dfe-Client-Id":[]string{"am-android-google"},
      "X-Dfe-Content-Filters":[]string{""},
      "X-Dfe-Cookie":[]string{"EAEYACICVVMyUENqZ2FOZ29UTkRJeE5UWXhOamcyT1RNd05qSTROVEUzTWhJZkNoQXhOalF4TURFeU1UY3lNakE1TWpReUVnc0l6TGVfamdZUWtJX2pZdz09QhUKBVVTLVRYEgwIzIvTjwYQ0OC9swFKEgoCVVMSDAjMi9OPBhDYgcCzAQ"},
      "X-Dfe-Device-Checkin-Consistency-Token":[]string{"ABFEt1URAJX3Okomj-L3xxsrntpt82bd-fHem4xV8GHm2IXdjt-26CLO4KRXYyCTJlkTz2pz3-I_6b8E0pD-dbHcqD9sAONl-aGeXjjmIGd1BdJ-0uo4rKd7YoXvHN613xzRzN3z-n5sQ0bP32bQ_2hemJxBJf5oooRBsnJkvgBt6Fw2gBxOe3iUL7LRiRMtZ0j_uwqPXR79vu71zSLxI4RLWOQRmzwP7mfekvIkfJfXHdRq3TlwNH2n4vlmwhQdZ6QWsvrD-QleZeiptJTGhASYxDJNUDUPlm-HY17pBRzIPc9XAxOTh-8"},
      "X-Dfe-Device-Config-Token":[]string{"CjgaNgoTNDIxNTYxNjg2OTMwNjI4NTE3MhIfChAxNjQxMDEyMTcyMjA5MjQyEgsIzLe_jgYQkI_jYw=="},
      "X-Dfe-Encoded-Targets":[]string{"CAESN/qigQYC2AMBFfUbyA7SM5Ij/CvfBoIDgxHqGP8R3xzIBvoQtBKFDZ4HAY4FrwSVMasHBO0O2Q8akgYRAQECAQO7AQEpKZ0CnwECAwRrAQYBr9PPAoK7sQMBAQMCBAkIDAgBAwEDBAICBAUZEgMEBAMLAQEBBQEBAcYBARYED+cBfS8CHQEKkAEMMxcBIQoUDwYHIjd3DQ4MFk0JWGYZEREYAQOLAYEBFDMIEYMBAgICAgICOxkCD18LGQKEAcgDBIQBAgGLARkYCy8oBTJlBCUocxQn0QUBDkkGxgNZQq0BZSbeAmIDgAEBOgGtAaMCDAOQAZ4BBIEBKUtQUYYBQscDDxPSARA1oAEHAWmnAsMB2wFyywGLAxol+wImlwOOA80CtwN26A0WjwJVbQEJPAH+BRDeAfkHK/ABASEBCSAaHQemAzkaRiu2Ad8BdXeiAwEBGBUBBN4LEIABK4gB2AFLfwECAdoENq0CkQGMBsIBiQEtiwGgA1zyAUQ4uwS8AwhsvgPyAcEDF27vApsBHaICGhl3GSKxAR8MC6cBAgItmQYG9QIeywLvAeYBDArLAh8HASI4ELICDVmVBgsY/gHWARtcAsMBpALiAdsBA7QBpAJmIArpByn0AyAKBwHTARIHAX8D+AMBcRIBBbEDmwUBMacCHAciNp0BAQF0OgQLJDuSAh54kwFSP0eeAQQ4M5EBQgMEmwFXywFo0gFyWwMcapQBBugBPUW2AVgBKmy3AR6PAbMBGQxrUJECvQR+8gFoWDsYgQNwRSczBRXQAgtRswEW0ALMAREYAUEBIG6yATYCRE8OxgER8gMBvQEDRkwLc8MBTwHZAUOnAXiiBakDIbYBNNcCIUmuArIBSakBrgFHKs0EgwV/G3AD0wE6LgECtQJ4xQFwFbUCjQPkBS6vAQqEAUZF3QIM9wEhCoYCQhXsBCyZArQDugIziALWAdIBlQHwBdUErQE6qQaSA4EEIvYBHir9AQVLmgMCApsCKAwHuwgrENsBAjNYswEVmgIt7QJnN4wDEnta+wGfAcUBxgEtEFXQAQWdAUAeBcwBAQM7rAEJATJ0LENrdh73A6UBhAE+qwEeASxLZUMhDREuH0CGARbd7K0GlQo"},
      "X-Dfe-Mccmnc":[]string{"20815"}, "X-Dfe-Network-Type":[]string{"4"},
      "X-Dfe-Phenotype":[]string{"H4sIAAAAAAAAAB3OO3KjMAAA0KRNuWXukBkBQkAJ2MhgAZb5u2GCwQZbCH_EJ77QHmgvtDtbv-Z9_H63zXXU0NVPB1odlyGy7751Q3CitlPDvFd8lxhz3tpNmz7P92CFw73zdHU2Ie0Ad2kmR8lxhiErTFLt3RPGfJQHSDy7Clw10bg8kqf2owLokN4SecJTLoSwBnzQSd652_MOf2d1vKBNVedzg4ciPoLz2mQ8efGAgYeLou-l-PXn_7Sna1MfhHuySxt-4esulEDp8Sbq54CPPKjpANW-lkU2IZ0F92LBI-ukCKSptqeq1eXU96LD9nZfhKHdtjSWwJqUm_2r6pMHOxk01saVanmNopjX3YxQafC4iC6T55aRbC8nTI98AF_kItIQAJb5EQxnKTO7TZDWnr01HVPxelb9A2OWX6poidMWl16K54kcu_jhXw-JSBQkVcD_fPsLSZu6joIBAAA"},
      "X-Dfe-Request-Params":[]string{"timeoutMs=4000"},
      "X-Dfe-Userlanguages":[]string{"en_GB"},
      "X-Limit-Ad-Tracking-Enabled":[]string{"false"},
   }
   req.URL.RawQuery = "doc=com.google.android.projection.gearhead.phonescreen"
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return 0, err
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(1,2,4,13,1,3), nil
}
