package googleplay

import (
   "crypto/rsa"
   "crypto/sha1"
   "encoding/base64"
   "fmt"
   "github.com/89z/parse/bytes"
   "github.com/89z/parse/tls"
   "io"
   "math/big"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
   stdbytes "bytes"
)

const (
   Origin = "https://android.clients.google.com"
   androidKey = "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="
)

// device is Google Service Framework.
func Details(oauth url.Values, device, app string) ([]byte, error) {
   req, err := http.NewRequest("GET", Origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {
         "Bearer " + oauth.Get("Auth"),
      },
      "X-DFE-Device-Id": {device},
   }
   val := url.Values{
      "doc": {app},
   }
   req.URL.RawQuery = val.Encode()
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}

// Exchange refresh token (aas_et) for access token (Auth).
func OAuth2(aas_et url.Values) (url.Values, error) {
   val := url.Values{
      "Token": {
         aas_et.Get("Token"),
      },
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }
   req, err := http.NewRequest(
      "POST", Origin + "/auth", strings.NewReader(val.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Content-Type": {"application/x-www-form-urlencoded"},
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return ParseQuery(query)
}

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ParseQuery(query []byte) (url.Values, error) {
   res := make(url.Values)
   for _, pair := range stdbytes.Split(query, []byte{'\n'}) {
      nv := stdbytes.SplitN(pair, []byte{'='}, 2)
      if len(nv) != 2 {
         return nil, fmt.Errorf("%q", query)
      }
      res.Add(string(nv[0]), string(nv[1]))
   }
   return res, nil
}

func Signature(email, password string) (string, error) {
   data, err := base64.StdEncoding.DecodeString(androidKey)
   if err != nil {
      return "", err
   }
   var key rsa.PublicKey
   buf := bytes.NewBuffer(data)
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
      msg stdbytes.Buffer
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
   // FIXME is this line actually needed?
   msg.WriteByte(0)
   msg.Write(hash[:4])
   msg.Write(login)
   return base64.URLEncoding.EncodeToString(msg.Bytes()), nil
}

// Token=aas_et/AKppINa1sGeVY7ukPvr-v5Djm5fp0-oQY72xu7JmWPSR_GpFmLterv2fjgI8m3...
func Token(email, encryptedPasswd string) (url.Values, error) {
   hello, err := tls.ParseJA3(tls.Android)
   if err != nil {
      return nil, err
   }
   val := url.Values{
      "Email": {email},
      "EncryptedPasswd": {encryptedPasswd},
      "sdk_version": {"17"},
   }
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/auth",
      strings.NewReader(val.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   dum, err := httputil.DumpRequest(req, true)
   if err != nil {
      return nil, err
   }
   os.Stdout.Write(append(dum, '\n'))
   res, err := tls.NewTransport(hello.ClientHelloSpec).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return ParseQuery(query)
}

type devZero struct{}

func (devZero) Read(b []byte) (int, error) {
   return len(b), nil
}
