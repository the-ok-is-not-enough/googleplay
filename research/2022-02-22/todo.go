package googleplay

import (
   "bufio"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

const (
   multiple = 9999_9999
   single = 8091_9999
)

func (t Token) Header(dev *Device) (*Header, error) {
   val := url.Values{
      "Token": {t.Token},
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth", strings.NewReader(val),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf := bufio.NewScanner(res.Body)
   var head Header
   head.Header = make(http.Header)
   head.Set("X-DFE-Device-ID", strconv.FormatUint(dev.AndroidID, 16))
   head.Set("User-Agent", "Android-Finsky (sdk=9,versionCode=99999999")
   for buf.Scan() {
      kv := strings.SplitN(buf.Text(), "=", 2)
      if len(kv) == 2 && kv[0] == "Auth" {
         head.Set("Authorization", "Bearer " + kv[1])
      }
   }
   return &head, nil
}
