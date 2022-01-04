package googleplay

import (
   "bufio"
   "bytes"
   "crypto/rsa"
   "crypto/sha1"
   "encoding/base64"
   "encoding/json"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "io"
   "math/big"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

const (
   Sleep = 16 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
   origin = "https://android.clients.google.com"
)

var LogLevel format.LogLevel

// Purchase app. Only needs to be done once per Google account.
func (a Auth) Purchase(dev *Device, app string) error {
   query := "doc=" + url.QueryEscape(app)
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/purchase", strings.NewReader(query),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Auth},
      "Content-Type": {"application/x-www-form-urlencoded"},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

// Read Device from file.
func (d *Device) Decode(src io.Reader) error {
   return json.NewDecoder(src).Decode(d)
}

// Write Device to file.
func (d Device) Encode(dst io.Writer) error {
   enc := json.NewEncoder(dst)
   enc.SetIndent("", " ")
   return enc.Encode(d)
}

type NumDownloads struct {
   Value uint64
}

func (n NumDownloads) String() string {
   return format.Number.LabelUint(n.Value)
}

type Size struct {
   Value uint64
}

func (i Size) String() string {
   return format.Size.LabelUint(i.Value)
}

type response struct {
   *http.Response
}

func (r response) Error() string {
   return strconv.Itoa(r.StatusCode) + " " + r.Status
}

var purchaseRequired = response{
   &http.Response{StatusCode: 3, Status: "purchase required"},
}

type Offer struct {
   Micros uint64
   CurrencyCode string
}

func (o Offer) String() string {
   val := float64(o.Micros) / 1_000_000
   return strconv.FormatFloat(val, 'f', 2, 64) + " " + o.CurrencyCode
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}

const androidKey =
   "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp" +
   "5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLN" +
   "WgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

func signature(email, password string) (string, error) {
   data, err := base64.StdEncoding.DecodeString(androidKey)
   if err != nil {
      return "", err
   }
   var key rsa.PublicKey
   buf := crypto.NewBuffer(data)
   // modulus_length | modulus | exponent_length | exponent
   _, mod, ok := buf.ReadUint32LengthPrefixed()
   if ok {
      key.N = new(big.Int).SetBytes(mod)
   }
   _, exp, ok := buf.ReadUint32LengthPrefixed()
   if ok {
      exp := new(big.Int).SetBytes(exp).Int64()
      key.E = int(exp)
   }
   var (
      msg bytes.Buffer
      nop nopSource
   )
   msg.WriteString(email)
   msg.WriteByte(0)
   msg.WriteString(password)
   login, err := rsa.EncryptOAEP(
      sha1.New(), nop, &key, msg.Bytes(), nil,
   )
   if err != nil {
      return "", err
   }
   hash := sha1.Sum(data)
   msg.Reset()
   msg.WriteByte(0)
   msg.Write(hash[:4])
   msg.Write(login)
   return base64.URLEncoding.EncodeToString(msg.Bytes()), nil
}

type Auth struct {
   Auth string
}

type Token struct {
   Token string
}

// Request refresh token.
func NewToken(email, password string) (*Token, error) {
   hello, err := crypto.ParseJA3(crypto.AndroidJA3)
   if err != nil {
      return nil, err
   }
   sig, err := signature(email, password)
   if err != nil {
      return nil, err
   }
   val := url.Values{
      "Email": {email},
      "EncryptedPasswd": {sig},
      "sdk_version": {"20"}, // Newer versions fail.
   }.Encode()
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   LogLevel.Dump(req)
   res, err := hello.Transport().RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf := bufio.NewScanner(res.Body)
   for buf.Scan() {
      kv := strings.SplitN(buf.Text(), "=", 2)
      if len(kv) == 2 && kv[0] == "Token" {
         var tok Token
         tok.Token = kv[1]
         return &tok, nil
      }
   }
   return nil, notFound{"Token"}
}

// Exchange refresh token for access token.
func (t Token) Auth() (*Auth, error) {
   val := url.Values{
      "Token": {t.Token},
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val),
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
      return nil, response{res}
   }
   buf := bufio.NewScanner(res.Body)
   for buf.Scan() {
      kv := strings.SplitN(buf.Text(), "=", 2)
      if len(kv) == 2 && kv[0] == "Auth" {
         var auth Auth
         auth.Auth = kv[1]
         return &auth, nil
      }
   }
   return nil, notFound{"Auth"}
}

// Read Token from file.
func (t *Token) Decode(src io.Reader) error {
   return json.NewDecoder(src).Decode(t)
}

// Write Token to file.
func (t Token) Encode(dst io.Writer) error {
   enc := json.NewEncoder(dst)
   enc.SetIndent("", " ")
   return enc.Encode(t)
}

type nopSource struct{}

func (nopSource) Read(buf []byte) (int, error) {
   return len(buf), nil
}

type notFound struct {
   input string
}

func (n notFound) Error() string {
   return strconv.Quote(n.input) + " not found"
}
