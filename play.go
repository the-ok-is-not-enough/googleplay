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

func (d Details) Format(f fmt.State, verb rune) {
   // Title
   fmt.Fprintln(f, "Title:", d.Title)
   // UploadDate
   fmt.Fprintln(f, "UploadDate:", d.UploadDate)
   // VersionString
   fmt.Fprintln(f, "VersionString:", d.VersionString)
   // VersionCode
   fmt.Fprintln(f, "VersionCode:", d.VersionCode)
   // NumDownloads
   fmt.Fprintln(f, "NumDownloads:", d.NumDownloads)
   // Size
   fmt.Fprintln(f, "Size:", format.LabelSize(d.Size))
   // Files
   fmt.Fprintln(f, "Files:", d.Files)
   // Offer
   fmt.Fprint(f, "Offer: ", d.Micros, " ", d.CurrencyCode)
   // Images
   if verb == 'a' {
      for _, ima := range d.Images {
         fmt.Fprint(f, "\nType:", ima.Type)
         fmt.Fprint(f, " URL:", ima.URL)
      }
   }
}

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

func encode(value any, elem ...string) error {
   name := filepath.Join(elem...)
   err := os.MkdirAll(filepath.Dir(name), os.ModeDir)
   if err != nil {
      return err
   }
   fmt.Println("Create", name)
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   return json.NewEncoder(file).Encode(value)
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

type Delivery struct {
   DownloadURL string
   SplitDeliveryData []SplitDeliveryData
}

type Details struct {
   CurrencyCode string
   Files int
   Images []Image
   Micros uint64
   NumDownloads uint64
   Size uint64
   Title string
   UploadDate string
   VersionCode uint64
   VersionString string
}

func (d Details) Icon() string {
   for _, image := range d.Images {
      if image.Type == 4 {
         return image.URL
      }
   }
   return ""
}

type Image struct {
   Type uint64
   URL string
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
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
   auth := parseQuery(res.Body).Get("Auth")
   if auth != "" {
      head.Set("Authorization", "Bearer " + auth)
   }
   // User-Agent
   head.Set(
      "User-Agent", fmt.Sprint("Android-Finsky (sdk=9,versionCode=", version),
   )
   // X-DFE-Device-ID
   head.Set("X-DFE-Device-ID", fmt.Sprintf("%x", dev.AndroidID))
   return &head, nil
}

type errorString string

func (e errorString) Error() string {
   return string(e)
}
