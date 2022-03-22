package googleplay

import (
   "bufio"
   "encoding/json"
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "io"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "strings"
   "time"
)

const Sleep = 4 * time.Second

var LogLevel format.LogLevel

func decode(value any, elem ...string) error {
   name := filepath.Join(elem...)
   file, err := os.Open(name)
   if err != nil {
      return err
   }
   defer file.Close()
   return json.NewDecoder(file).Decode(value)
}

func encode(val interface{}, elem ...string) error {
   name := filepath.Join(elem...)
   err := os.MkdirAll(filepath.Dir(name), os.ModeDir)
   if err != nil {
      return err
   }
   os.Stdout.WriteString("Create " + name + "\n")
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   return json.NewEncoder(file).Encode(val)
}

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
   Token string
}

// You can also use host "android.clients.google.com", but it also uses
// TLS fingerprinting.
func NewToken(email, password string) (*Token, error) {
   val := url.Values{
      "Email": {email},
      "Passwd": {password},
      "client_sig": {""},
      "droidguard_results": {""},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth", strings.NewReader(val),
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
   var tok Token
   tok.Token = parseQuery(res.Body).Get("Token")
   return &tok, nil
}

func OpenToken(elem ...string) (*Token, error) {
   tok := new(Token)
   err := decode(tok, elem...)
   if err != nil {
      return nil, err
   }
   return tok, nil
}

func (t Token) Create(elem ...string) error {
   return encode(t, elem...)
}

func (t Token) Header(dev *Device) (*Header, error) {
   return t.headerVersion(dev, 9999_9999)
}

func (t Token) SingleAPK(dev *Device) (*Header, error) {
   return t.headerVersion(dev, 8091_9999)
}

type errorString string

func (e errorString) Error() string {
   return string(e)
}

func (d Details) Format(f fmt.State, verb rune) {
   fmt.Fprintln(f, "Title:", d.Title)
   fmt.Fprintln(f, "UploadDate:", d.UploadDate)
   fmt.Fprintln(f, "VersionString:", d.VersionString)
   fmt.Fprintln(f, "VersionCode:", d.VersionCode)
   fmt.Fprintln(f, "NumDownloads:", d.NumDownloads)
   fmt.Fprintln(f, "Size:", d.Size)
   fmt.Fprintln(f, "Files:", d.Files)
   fmt.Fprint(f, "Offer: ", d.Micros, " ", d.CurrencyCode)
}

////////////////////////////////////////////////////////////////////////////////

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
