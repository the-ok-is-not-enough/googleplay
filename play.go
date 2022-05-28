// github.com/89z
package googleplay

import (
   "bufio"
   "errors"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

const Sleep = 4 * time.Second

var LogLevel format.LogLevel

func parseQuery(query io.Reader) url.Values {
   vals := make(url.Values)
   buf := bufio.NewScanner(query)
   for buf.Scan() {
      key, val, ok := strings.Cut(buf.Text(), "=")
      if ok {
         vals.Add(key, val)
      }
   }
   return vals
}

type Header struct {
   Auth string // Authorization
   SDK int64 // User-Agent
   VersionCode int64 // User-Agent
   AndroidID uint64 // X-DFE-Device-ID
}

func (h Header) SetAgent(head http.Header) {
   var buf []byte
   buf = append(buf, "Android-Finsky (sdk="...)
   buf = strconv.AppendInt(buf, h.SDK, 10)
   buf = append(buf, ",versionCode="...)
   buf = strconv.AppendInt(buf, h.VersionCode, 10)
   head.Set("User-Agent", string(buf))
}

func (h Header) SetAuth(head http.Header) {
   head.Set("Authorization", "Bearer " + h.Auth)
}

func (h Header) SetDevice(head http.Header) {
   device := strconv.FormatUint(h.AndroidID, 16)
   head.Set("X-DFE-Device-ID", device)
}

// Purchase app. Only needs to be done once per Google account.
func (h Header) Purchase(app string) error {
   query := "doc=" + url.QueryEscape(app)
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/fdfe/purchase",
      strings.NewReader(query),
   )
   if err != nil {
      return err
   }
   h.SetAuth(req.Header)
   h.SetDevice(req.Header)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errors.New(res.Status)
   }
   return nil
}

type Token struct {
   Services string
   Token string
}

// You can also use host "android.clients.google.com", but it also uses
// TLS fingerprinting.
func NewToken(email, password string) (*Token, error) {
   body := url.Values{
      "Email": {email},
      "Passwd": {password},
      "client_sig": {""},
      "droidguard_results": {""},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth", strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   hello, err := crypto.ParseJA3(crypto.AndroidAPI26)
   if err != nil {
      return nil, err
   }
   LogLevel.Dump(req)
   res, err := crypto.Transport(hello).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   val := parseQuery(res.Body)
   var tok Token
   tok.Services = val.Get("services")
   tok.Token = val.Get("Token")
   return &tok, nil
}

func OpenToken(elem ...string) (*Token, error) {
   return format.Open[Token](elem...)
}

func (t Token) Create(elem ...string) error {
   return format.Create(t, elem...)
}

func (t Token) Header(androidID uint64, single bool) (*Header, error) {
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
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   var head Header
   head.Auth = parseQuery(res.Body).Get("Auth")
   head.SDK = 9
   head.AndroidID = androidID
   if single {
      head.VersionCode = 8091_9999 // single APK
   } else {
      head.VersionCode = 9999_9999
   }
   return &head, nil
}
