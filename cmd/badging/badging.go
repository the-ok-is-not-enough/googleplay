package main

import (
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
      buf, err := exec.Command("aapt", "dump", "badging", name).Output()
      if err != nil {
         panic(err)
      }
      lines := strings.FieldsFunc(string(buf), func(r rune) bool {
         return r == '\n'
      })
      for _, line := range lines {
         if verbose ||
         strings.HasPrefix(line, "  uses-feature:") ||
         strings.HasPrefix(line, "native-code:") ||
         strings.HasPrefix(line, "uses-library-not-required:") ||
         strings.HasPrefix(line, "uses-library:") {
            fmt.Println(line)
         }
      }
   } else {
      flag.Usage()
   }
}
