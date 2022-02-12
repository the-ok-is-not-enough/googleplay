package main

import (
   "bufio"
   "flag"
   "fmt"
   "os/exec"
   "strings"
)

func main() {
   // f
   var name string
   flag.StringVar(&name, "f", "", "file")
   // v
   var verbose bool
   flag.BoolVar(&verbose, "v", false, "verbose")
   flag.Parse()
   if name != "" {
      cmd := exec.Command("aapt", "dump", "badging", name)
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
      flag.PrintDefaults()
   }
}
