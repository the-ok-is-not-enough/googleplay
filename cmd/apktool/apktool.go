package main

import (
   "os"
   "os/exec"
)

const apktool = `C:\Program Files\Android\apktool_2.6.0.jar`

func main() {
   arg := []string{"-jar", apktool}
   arg = append(arg, os.Args[1:]...)
   cmd := exec.Command("java.exe", arg...)
   cmd.Stderr = os.Stderr
   cmd.Stdout = os.Stdout
   err := cmd.Run()
   if err != nil {
      panic(err)
   }
}
