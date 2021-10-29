package encrypt

import (
   "crypto/rsa"
   "crypto/sha1"
   "encoding/base64"
   "encoding/binary"
   "fmt"
   "math/big"
)

var (
   _ = binary.BigEndian
   _ = fmt.Print
)

type devZero struct{}

func (devZero) Read(b []byte) (int, error) {
   return len(b), nil
}

const androidKey = "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

func bytesToLong(b []byte) *big.Int {
   return new(big.Int).SetBytes(b)
}

func signature(email, password string) (string, error) {
   data, err := base64.StdEncoding.DecodeString(androidKey)
   if err != nil {
      return "", err
   }
   nLen := binary.BigEndian.Uint32(data[:4])
   var key rsa.PublicKey
   key.N = bytesToLong(data[4:4+nLen])
   key.E = int(bytesToLong(data[4+nLen+4:]).Int64())
   hash := sha1.Sum(data)
   msg := append([]byte(email), 0)
   msg = append(msg, []byte(password)...)
   var zero devZero
   encryptedLogin, err := rsa.EncryptOAEP(
      sha1.New(), zero, &key, msg, nil,
   )
   if err != nil {
      return "", err
   }
   sig := append([]byte{0}, hash[:4]...)
   sig = append(sig, encryptedLogin...)
   return base64.URLEncoding.EncodeToString(sig), nil
}
