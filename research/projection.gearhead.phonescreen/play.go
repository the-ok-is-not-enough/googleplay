package discord

import (
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/http/httputil"
   "os"
   "net/url"
   "strconv"
)

const auth = "ya29.a0ARrdaM9OGGr4x_vWHEK7qTD9YuK1K51UCzf2zKQlKaUP0vvRw3hTSnUdy2WHA6jTKHKZs7cy7RJfEnlTwqWAl818fTcyU10PabF0zBdgkxjl3Cqfmb5_VCrc4vKIeQIMvHN5DMpHO5yNkJKcNeHQ-ow10jlUGGsT2jNwQX-a62u0u7Jyal29sdAc-0M7jzvDD04ZoQR0T2uY9t5PHC086fiSf0tbfWPN_jg7LqKq-5fYNb2Tx7MdM59BNZsEP8HlicHlzN9i0paG21nyv2ymPoV-na5W5tIdnbwt_xbLMHns4FJzu70"

func checkin() (uint64, error) {
   var req0 = &http.Request{Method:"POST", URL:&url.URL{Scheme:"https",
Opaque:"", User:(*url.Userinfo)(nil), Host:"android.clients.google.com",
Path:"/checkin", RawPath:"", ForceQuery:false, RawQuery:"", Fragment:"",
RawFragment:""}, Header:http.Header{"App":[]string{"com.google.android.gms"},
"Content-Type":[]string{"application/x-protobuffer"},
"Host":[]string{"android.clients.google.com"},
"User-Agent":[]string{"GoogleAuth/1.4 sargo PQ3B.190705.003"}},
Body:io.NopCloser(body0)}
   res, err := new(http.Transport).RoundTrip(req0)
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

func details(app string) (uint64, error) {
   id, err := checkin()
   if err != nil {
      return 0, err
   }
   sID := strconv.FormatUint(id, 16)
   var req5 = &http.Request{Method:"GET", URL:&url.URL{Scheme:"https",
      Opaque:"", User:(*url.Userinfo)(nil), Host:"android.clients.google.com",
      Path:"/fdfe/details", RawPath:"", ForceQuery:false, RawQuery:"doc=" + app,
      Fragment:"", RawFragment:""},
      Header:http.Header{
         "Authorization":[]string{"Bearer " + auth},
         "X-Dfe-Device-Id":[]string{sID},
         "Accept-Language":[]string{"en-GB"},
         "User-Agent":[]string{"Android-Finsky/15.8.23-all [0] [PR] 259261889 (api=3,versionCode=81582300,sdk=28,device=sargo,hardware=sargo,product=sargo,platformVersionRelease=9,model=Pixel 3a,buildId=PQ3B.190705.003,isWideScreen=0,supportedAbis=arm64-v8a;armeabi-v7a;armeabi)"},
         "X-Ad-Id":[]string{"LawadaMera"},
         "X-Dfe-Client-Id":[]string{"am-android-google"},
         "X-Dfe-Content-Filters":[]string{""},
         "X-Dfe-Cookie":[]string{"EAEYACICVVMyUENqZ2FOZ29UTkRNeU9EUXlPVGcyTURFeU9USTVOalUwTnhJZkNoQXhOalF4TURnME1EUXdNakUzT1RreEVnc0lpT25EamdZUTJJNzVadz09QhUKBVVTLVRYEgwIiL3XjwYQgP6itwFKEgoCVVMSDAiIvdePBhCIlqi3AQ"},
         "X-Dfe-Device-Checkin-Consistency-Token":[]string{"ABFEt1UYSrT0WNYJ9UY896IamMZw30SWXrVV08ilR27i6zQVaDtUNCk6rk9Wcgk2yb6ADiezHJjK4YwErnEomiDRMkoc4WKCBmO-cSaSYjgCKSlTvBHk43JAmHn8vTyuUqh_Sr_9iCBj8boqFhQH7uBn0-EytAof0GO2lo0zb9E6kfsnw85XhPBY5Xw4sHdKS2J6YsImT-sFpKA6eaK7pFf1i7vY1NIMtdXe4jWSsInWJoGiNQIHDLFiJf-nG0tJbMcyyESIsPEu7tLIKzDr3_1v_IbVAy1wjDvfB6tNVhd57CvGcuWvY_g"},
         "X-Dfe-Device-Config-Token":[]string{"CjgaNgoTNDMyODQyOTg2MDEyOTI5NjU0NxIfChAxNjQxMDg0MDQwMjE3OTkxEgsIiOnDjgYQ2I75Zw=="},
         "X-Dfe-Encoded-Targets":[]string{"CAESN/qigQYC2AMBFfUbyA7SM5Ij/CvfBoIDgxHqGP8R3xzIBvoQtBKFDZ4HAY4FrwSVMasHBO0O2Q8akgYRAQECAQO7AQEpKZ0CnwECAwRrAQYBr9PPAoK7sQMBAQMCBAkIDAgBAwEDBAICBAUZEgMEBAMLAQEBBQEBAcYBARYED+cBfS8CHQEKkAEMMxcBIQoUDwYHIjd3DQ4MFk0JWGYZEREYAQOLAYEBFDMIEYMBAgICAgICOxkCD18LGQKEAcgDBIQBAgGLARkYCy8oBTJlBCUocxQn0QUBDkkGxgNZQq0BZSbeAmIDgAEBOgGtAaMCDAOQAZ4BBIEBKUtQUYYBQscDDxPSARA1oAEHAWmnAsMB2wFyywGLAxol+wImlwOOA80CtwN26A0WjwJVbQEJPAH+BRDeAfkHK/ABASEBCSAaHQemAzkaRiu2Ad8BdXeiAwEBGBUBBN4LEIABK4gB2AFLfwECAdoENq0CkQGMBsIBiQEtiwGgA1zyAUQ4uwS8AwhsvgPyAcEDF27vApsBHaICGhl3GSKxAR8MC6cBAgItmQYG9QIeywLvAeYBDArLAh8HASI4ELICDVmVBgsY/gHWARtcAsMBpALiAdsBA7QBpAJmIArpByn0AyAKBwHTARIHAX8D+AMBcRIBBbEDmwUBMacCHAciNp0BAQF0OgQLJDuSAh54kwFSP0eeAQQ4M5EBQgMEmwFXywFo0gFyWwMcapQBBugBPUW2AVgBKmy3AR6PAbMBGQxrUJECvQR+8gFoWDsYgQNwRSczBRXQAgtRswEW0ALMAREYAUEBIG6yATYCRE8OxgER8gMBvQEDRkwLc8MBTwHZAUOnAXiiBakDIbYBNNcCIUmuArIBSakBrgFHKs0EgwV/G3AD0wE6LgECtQJ4xQFwFbUCjQPkBS6vAQqEAUZF3QIM9wEhCoYCQhXsBCyZArQDugIziALWAdIBlQHwBdUErQE6qQaSA4EEIvYBHir9AQVLmgMCApsCKAwHuwgrENsBAjNYswEVmgIt7QJnN4wDEnta+wGfAcUBxgEtEFXQAQWdAUAeBcwBAQM7rAEJATJ0LENrdh73A6UBhAE+qwEeASxLZUMhDREuH0CGARbd7K0GlQo"},
         "X-Dfe-Mccmnc":[]string{"20815"},
         "X-Dfe-Network-Type":[]string{"4"},
         "X-Dfe-Phenotype":[]string{"H4sIAAAAAAAAAB3OO3KjMAAA0KRNuWXukBkBQkAJ2MhgAZb5u2GCwQZbCH_EJ77QHmgvtDtbv-Z9_H63zXXU0NVPB1odlyGy7751Q3CitlPDvFd8lxhz3tpNmz7P92CFw73zdHU2Ie0Ad2kmR8lxhiErTFLt3RPGfJQHSDy7Clw10bg8kqf2owLokN4SecJTLoSwBnzQSd652_MOf2d1vKBNVedzg4ciPoLz2mQ8efGAgYeLou-l-PXn_7Sna1MfhHuySxt-4esulEDp8Sbq54CPPKjpANW-lkU2IZ0F92LBI-ukCKSptqeq1eXU96LD9nZfhKHdtjSWwJqUm_2r6pMHOxk01saVanmNopjX3YxQafC4iC6T55aRbC8nTI98AF_kItIQAJb5EQxnKTO7TZDWnr01HVPxelb9A2OWX6poidMWl16K54kcu_jhXw-JSBQkVcD_fPsLSZu6joIBAAA"},
         "X-Dfe-Request-Params":[]string{"timeoutMs=4000"},
         "X-Dfe-Userlanguages":[]string{"en_GB"},
         "X-Limit-Ad-Tracking-Enabled":[]string{"false"},
      },
   }
   buf, err := httputil.DumpRequest(req5, false)
   if err != nil {
      return 0, err
   }
   os.Stdout.Write(buf)
   res, err := new(http.Transport).RoundTrip(req5)
   if err != nil {
      return 0, err
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(1,2,4,13,1,3), nil
}
