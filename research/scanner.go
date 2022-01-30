package main

import (
   "fmt"
   "strings"
   "text/scanner"
)

// no final newline
var src = strings.NewReader(`ExpiresInDurationSec=3599
Auth=ya29.A0ARrdaM_dhAAGoHSMRwEH4-egUZFEFdzGcW4NDBbrhlJbS7N3qeaUkotalmG5AWiJ-dBo
Expiry=1643558947
grantedScopes=https://www.googleapis.com/auth/googleplay
isTokenSnowballed=0
issueAdvice=auto
services=friendview,multilogin,cloudconsole,nova,sierra,mobile,ahsid,billing,hist,androiddeveloper,dynamite,cl,talk,wise,omaha,sitemaps,youtube,mail,googleplay,gerritcodereview,android
storeConsentRemotely=0`)

func main() {
   var text scanner.Scanner
   text.Init(src)
   text.IsIdentRune = func(r rune, i int) bool {
      return r != '=' && r != '\n' && r != scanner.EOF
   }
   for text.Scan() != scanner.EOF {
      fmt.Printf("%q\n", text.TokenText())
   }
}
