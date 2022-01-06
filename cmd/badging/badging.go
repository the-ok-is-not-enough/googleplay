package main

import (
   "bufio"
   "flag"
   "fmt"
   "os/exec"
   "strings"
)

func main() {
   var verbose bool
   flag.BoolVar(&verbose, "v", false, "verbose")
   flag.Parse()
   if flag.NArg() == 1 {
      apk := flag.Arg(0)
      cmd := exec.Command("aapt", "dump", "badging", apk)
      pipe, err := cmd.StdoutPipe()
      if err != nil {
         panic(err)
      }
      if err := cmd.Start(); err != nil {
         panic(err)
      }
      defer cmd.Wait()
      buf := bufio.NewScanner(pipe)
      for buf.Scan() {
         text := buf.Text()
         if verbose ||
         strings.HasPrefix(text, "  uses-feature:") ||
         strings.HasPrefix(text, "native-code:") {
            fmt.Println(text)
         }
      }
   } else {
      fmt.Println("badging [APK]")
      flag.PrintDefaults()
   }
}
