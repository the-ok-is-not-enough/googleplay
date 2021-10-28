package googleplay

import (
   "bytes"
   "crypto/rand"
   "crypto/rsa"
   "crypto/sha1"
   "encoding/base64"
   "fmt"
   "github.com/89z/parse/tls"
   "io"
   "math/big"
   "net/http"
   "net/url"
   "strings"
)

const (
   Origin = "https://android.clients.google.com"
   androidKeyBase64 = "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="
)

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ParseQuery(query []byte) (url.Values, error) {
   res := make(url.Values)
   for _, pair := range bytes.Split(query, []byte{'\n'}) {
      nv := bytes.SplitN(pair, []byte{'='}, 2)
      if len(nv) != 2 {
         return nil, fmt.Errorf("%q", query)
      }
      res.Add(string(nv[0]), string(nv[1]))
   }
   return res, nil
}

func Signature(email, password string) (string, error) {
   androidKeyBytes, err := base64.StdEncoding.DecodeString(androidKeyBase64)
   if err != nil {
      return "", err
   }
   i := bytesToLong(androidKeyBytes[:4]).Int64()
   j := bytesToLong(androidKeyBytes[i+4 : i+8]).Int64()
   androidKey := &rsa.PublicKey{
      E: int(bytesToLong(androidKeyBytes[i+8 : i+8+j]).Int64()),
      N: bytesToLong(androidKeyBytes[4 : 4+i]),
   }
   hash := sha1.Sum(androidKeyBytes)
   msg := append([]byte(email), 0)
   msg = append(msg, []byte(password)...)
   encryptedLogin, err := rsa.EncryptOAEP(
      sha1.New(), rand.Reader, androidKey, msg, nil,
   )
   if err != nil {
      return "", err
   }
   sig := append([]byte{0}, hash[:4]...)
   sig = append(sig, encryptedLogin...)
   return base64.URLEncoding.EncodeToString(sig), nil
}

func bytesToLong(b []byte) *big.Int {
   return new(big.Int).SetBytes(b)
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
