package googleplay

import (
   "bufio"
   "github.com/89z/format"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

var LogLevel format.LogLevel

func parseQuery(query io.Reader) url.Values {
   val := make(url.Values)
   buf := bufio.NewScanner(query)
   for buf.Scan() {
      var key string
      for i, elem := range strings.SplitN(buf.Text(), "=", 2) {
         switch i {
         case 0:
            key = elem
         case 1:
            val.Add(key, elem)
         }
      }
   }
   return val
}

type Device struct {
   AndroidID uint64
}

type Header struct {
   http.Header
}

type Token struct {
   Token string
}

func (t Token) Header(dev *Device) (*Header, error) {
   return t.headerVersion(dev, 9999_9999)
}

func (t Token) SingleAPK(dev *Device) (*Header, error) {
   return t.headerVersion(dev, 8091_9999)
}

func (t Token) headerVersion(dev *Device, version int64) (*Header, error) {
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
   var head Header
   head.Header = make(http.Header)
   auth := parseQuery(res.Body).Get("Auth")
   if auth != "" {
      head.Set("Authorization", "Bearer " + auth)
   }
   buf := []byte("Android-Finsky (sdk=9,versionCode=")
   buf = strconv.AppendInt(buf, version, 10)
   head.Set("User-Agent", string(buf))
   head.Set("X-DFE-Device-ID", strconv.FormatUint(dev.AndroidID, 16))
   return &head, nil
}
