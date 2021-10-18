package main

import (
   "os"
   "os/exec"
)

func main() {
   arg := []string{"-jar", `D:\Desktop\apktool_2.6.0.jar`}
   arg = append(arg, os.Args[1:]...)
   cmd := exec.Command(`D:\Desktop\jdk-17+35\bin\java`, arg...)
   cmd.Stdout = os.Stdout
   cmd.Stderr = os.Stderr
   err := cmd.Run()
   if err != nil {
      panic(err)
   }
}
