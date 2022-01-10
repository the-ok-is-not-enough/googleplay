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
   "math/big"
   "net/http"
   "net/url"
   "os"
   "strings"
)

const googlePublicKey =
   "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp" +
   "5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLN" +
   "WgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

func cryptPass(email, password string) (string, error) {
   buf, err := base64.StdEncoding.DecodeString(googlePublicKey)
   if err != nil {
      return "", err
   }
   var key rsa.PublicKey
   read := crypto.NewReader(buf)
   // modulus_length | modulus | exponent_length | exponent
   _, mod, ok := read.ReadUint32LengthPrefixed()
   if ok {
      key.N = new(big.Int).SetBytes(mod)
   }
   _, exp, ok := read.ReadUint32LengthPrefixed()
   if ok {
      exp := new(big.Int).SetBytes(exp).Int64()
      key.E = int(exp)
   }
   var (
      mes bytes.Buffer
      nop nopSource
   )
   mes.WriteString(email)
   mes.WriteByte(0)
   mes.WriteString(password)
   login, err := rsa.EncryptOAEP(
      sha1.New(), nop, &key, mes.Bytes(), nil,
   )
   if err != nil {
      return "", err
   }
   hash := sha1.Sum(buf)
   mes.Reset()
   mes.WriteByte(0)
   mes.Write(hash[:4])
   mes.Write(login)
   return base64.URLEncoding.EncodeToString(mes.Bytes()), nil
}

type Token struct {
   Token string
}

func OpenToken(name string) (*Token, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   tok := new(Token)
   if err := json.NewDecoder(file).Decode(tok); err != nil {
      return nil, err
   }
   return tok, nil
}

func (t Token) Create(name string) error {
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(t)
}

// Request refresh token.
func NewToken(email, password string) (*Token, error) {
   hello, err := crypto.ParseJA3(crypto.AndroidAPI26)
   if err != nil {
      return nil, err
   }
   cryptedPass, err := cryptPass(email, password)
   if err != nil {
      return nil, err
   }
   val := url.Values{
      "Email": {email},
      "EncryptedPasswd": {cryptedPass},
      "sdk_version": {"20"}, // Newer versions fail.
   }.Encode()
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   format.Log.Dump(req)
   res, err := crypto.Transport(hello).RoundTrip(req)
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
   format.Log.Dump(req)
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
