package googleplay

import (
   "bytes"
   "crypto/rsa"
   "crypto/sha1"
   "encoding/base64"
   "encoding/json"
   "github.com/89z/parse/crypto"
   "github.com/89z/parse/net"
   "io"
   "math/big"
   "net/http"
   "net/url"
   "strings"
)

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
   Auth string
}

type Token struct {
   Token string
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
   body := url.Values{
      "Email": {email},
      "EncryptedPasswd": {sig},
      "sdk_version": {"20"}, // Newer versions fail.
   }.Encode()
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   dumpRequest(req)
   res, err := crypto.NewTransport(hello.ClientHelloSpec).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   tok := net.ParseQuery(res.Body).Get("Token")
   return &Token{tok}, nil
}

// Exchange refresh token for access token.
func (t Token) Auth() (*Auth, error) {
   body := url.Values{
      "Token": {t.Token},
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Content-Type": {"application/x-www-form-urlencoded"},
   }
   dumpRequest(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   auth := net.ParseQuery(res.Body).Get("Auth")
   return &Auth{auth}, nil
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

type devZero struct{}

func (devZero) Read(b []byte) (int, error) {
   return len(b), nil
}
