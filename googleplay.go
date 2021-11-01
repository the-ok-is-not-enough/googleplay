package googleplay

import (
   "bytes"
   "crypto/rsa"
   "crypto/sha1"
   "encoding/base64"
   "fmt"
   "github.com/89z/parse/crypto"
   "io"
   "math/big"
   "net/http"
   "net/url"
   "strings"
)

const androidJA3 =
   "769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-" +
   "51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

const androidKey =
   "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp" +
   "5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLN" +
   "WgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

const origin = "https://android.clients.google.com"

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ParseQuery(query []byte) url.Values {
   res := make(url.Values)
   for _, pair := range bytes.Split(query, []byte{'\n'}) {
      nv := bytes.SplitN(pair, []byte{'='}, 2)
      if len(nv) != 2 {
         return nil
      }
      res.Add(string(nv[0]), string(nv[1]))
   }
   return res
}

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
      zero devZero
   )
   msg.WriteString(email)
   msg.WriteByte(0)
   msg.WriteString(password)
   login, err := rsa.EncryptOAEP(
      sha1.New(), zero, &key, msg.Bytes(), nil,
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
   url.Values
}

// device is Google Service Framework.
func (a Auth) Details(device, app string) ([]byte, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   val := req.URL.Query()
   val.Set("doc", app)
   req.URL.RawQuery = val.Encode()
   req.Header.Set("Authorization", "Bearer " + a.Get("Auth"))
   req.Header.Set("X-DFE-Device-ID", device)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("status %q", res.Status)
   }
   return io.ReadAll(res.Body)
}

type Token struct {
   url.Values
}

// Request refresh token.
func NewToken(email, password string) (*Token, error) {
   hello, err := crypto.ParseJA3(androidJA3)
   if err != nil {
      return nil, err
   }
   sig, err := signature(email, password)
   if err != nil {
      return nil, err
   }
   val := make(url.Values)
   val.Set("Email", email)
   val.Set("EncryptedPasswd", sig)
   val.Set("sdk_version", "17")
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := crypto.NewTransport(hello.ClientHelloSpec).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   if val := ParseQuery(query); val != nil {
      return &Token{val}, nil
   }
   return nil, fmt.Errorf("parseQuery %q", query)
}

// Exchange refresh token for access token.
func (t Token) Auth() (*Auth, error) {
   val := make(url.Values)
   val.Set("Token", t.Get("Token"))
   val.Set("service", "oauth2:https://www.googleapis.com/auth/googleplay")
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   if val := ParseQuery(query); val != nil {
      return &Auth{val}, nil
   }
   return nil, fmt.Errorf("parseQuery %q", query)
}

type devZero struct{}

func (devZero) Read(b []byte) (int, error) {
   return len(b), nil
}
