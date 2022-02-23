package googleplay

import (
   "bufio"
   "encoding/json"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "io"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "strconv"
   "strings"
   "time"
)

const Sleep = 4 * time.Second

var LogLevel format.LogLevel

func decode(val interface{}, elem ...string) error {
   name := filepath.Join(elem...)
   file, err := os.Open(name)
   if err != nil {
      return err
   }
   defer file.Close()
   return json.NewDecoder(file).Decode(val)
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
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(val)
}

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

type Delivery struct {
   DownloadURL string
   SplitDeliveryData []SplitDeliveryData
}

type Details struct {
   Title string
   UploadDate string
   VersionString string
   VersionCode uint64
   NumDownloads uint64
   Size uint64
   Micros uint64
   CurrencyCode string
}

func (d Details) String() string {
   buf := []byte("Title: ")
   buf = append(buf, d.Title...)
   buf = append(buf, "\nUploadDate: "...)
   buf = append(buf, d.UploadDate...)
   buf = append(buf, "\nVersionString: "...)
   buf = append(buf, d.VersionString...)
   buf = append(buf, "\nVersionCode: "...)
   buf = strconv.AppendUint(buf, d.VersionCode, 10)
   buf = append(buf, "\nNumDownloads: "...)
   buf = append(buf, format.Number.GetUint64(d.NumDownloads)...)
   buf = append(buf, "\nSize: "...)
   buf = append(buf, format.Size.GetUint64(d.Size)...)
   buf = append(buf, "\nOffer: "...)
   buf = strconv.AppendFloat(buf, float64(d.Micros)/1e6, 'f', 2, 64)
   buf = append(buf, ' ')
   buf = append(buf, d.CurrencyCode...)
   return string(buf)
}

type Document struct {
   ID string
   Title string
   Creator string
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}

type Token struct {
   Token string
}

// Request refresh token.
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

type errorString string

func (e errorString) Error() string {
   return string(e)
}
