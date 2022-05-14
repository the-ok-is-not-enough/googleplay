package googleplay

import (
   "bufio"
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
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
      return nil, errorString(res.Status)
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
   if res.StatusCode != http.StatusOK {
      return nil, errorString(res.Status)
   }
   var head Header
   head.Header = make(http.Header)
   // Authorization
   head.Set(
      "Authorization",
      "Bearer " + parseQuery(res.Body).Get("Auth"),
   )
   // User-Agent
   head.Set(
      "User-Agent",
      fmt.Sprintf("Android-Finsky (sdk=9,versionCode=%v", version),
   )
   // X-DFE-Device-ID
   head.Set(
      "X-DFE-Device-ID",
      fmt.Sprintf("%x", dev.AndroidID),
   )
   return &head, nil
}

type errorString string

func (e errorString) Error() string {
   return string(e)
}

const (
   // com.kakaogames.twodin
   Arm64 String = "arm64-v8a"
   // com.miui.weather2
   Armeabi String = "armeabi-v7a"
   // com.google.android.youtube
   X86 String = "x86"
)

type Device struct {
   AndroidID Fixed64
   TimeMsec Varint
}

func OpenDevice(elem ...string) (*Device, error) {
   return format.Open[Device](elem...)
}

func (d Device) Create(elem ...string) error {
   return format.Create(d, elem...)
}

type Fixed64 = protobuf.Fixed64

type Message = protobuf.Message

type String = protobuf.String

type Varint = protobuf.Varint

type Header struct {
   http.Header
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
   h.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header = h.Header
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errorString(res.Status)
   }
   return nil
}
