package main

import (
   "crypto/md5"
   "crypto/x509"
   "encoding/hex"
   "encoding/pem"
   "fmt"
   "os"
   "os/exec"
)

const (
   data = "/data/local/tmp/cacerts"
   system = "/system/etc/security/cacerts"
)

// outputs the MD5 "hash" of the certificate subject name
func subjectHashOld(buf []byte) ([]byte, error) {
   block, _ := pem.Decode(buf)
   cert, err := x509.ParseCertificate(block.Bytes)
   if err != nil {
      return nil, err
   }
   md := md5.Sum(cert.RawSubject)
   return []byte{md[3], md[2], md[1], md[0]}, nil
}

func main() {
   home, err := os.UserHomeDir()
   if err != nil {
      panic(err)
   }
   buf, err := os.ReadFile(home + "/.mitmproxy/mitmproxy-ca-cert.cer")
   if err != nil {
      panic(err)
   }
   hash, err := subjectHashOld(buf)
   if err != nil {
      panic(err)
   }
   push := hex.EncodeToString(hash) + ".0"
   commands := [][]string{
      {"adb", "shell", "mkdir", data},
      {"adb", "shell", "cp", system + "/*", data},
      {"adb", "push", "mitmproxy-ca-cert.cer", data + "/" + push},
      {"adb", "root"},
      {"adb", "shell", "mount", "-t", "tmpfs", "tmpfs", system},
      {"adb", "shell", "mv", data + "/*", system},
      {"adb", "shell", "chcon", "u:object_r:system_file:s0", system + "/*"},
   }
   for _, command := range commands {
      cmd := exec.Command(command[0], command[1:]...)
      cmd.Stderr = os.Stderr
      cmd.Stdout = os.Stdout
      fmt.Println(cmd.Args)
      err := cmd.Run()
      if err != nil {
         panic(err)
      }
   }
}
