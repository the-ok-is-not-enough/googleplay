package main

import (
   "crypto/rand"
   "crypto/rsa"
   "crypto/sha1"
   "encoding/base64"
   "fmt"
   "github.com/89z/parse/tls"
   "math/big"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
)

const (
   androidKeyBase64 = "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="
   email = "srpen6@gmail.com"
)

func bytesToLong(b []byte) *big.Int {
   return new(big.Int).SetBytes(b)
}

func signature(email, password string) (string, error) {
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
   encryptedLogin, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, androidKey, msg, nil)
   if err != nil {
      return "", err
   }
   sig := append([]byte{0}, hash[:4]...)
   sig = append(sig, encryptedLogin...)
   return base64.URLEncoding.EncodeToString(sig), nil
}

func main() {
   hello, err := tls.ParseJA3(tls.Android)
   if err != nil {
      panic(err)
   }
   encryptedPasswd, err := signature(email, password)
   if err != nil {
      panic(err)
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
      panic(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   fmt.Println("RoundTrip")
   res, err := tls.NewTransport(hello.ClientHelloSpec).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   fmt.Println("DumpResponse")
   dum, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(dum)
}
