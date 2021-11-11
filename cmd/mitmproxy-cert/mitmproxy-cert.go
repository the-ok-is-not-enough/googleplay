package main

import (
   "os"
   "os/exec"
)

const (
   data = "/data/local/tmp/cacerts"
   system = "/system/etc/security/cacerts"
)

var commands = [][]string{
   {"adb", "shell", "mkdir", data},
   {"adb", "shell", "cp", system + "/*", data},
   {"adb", "push", "c8750f0d.0", data},
   {"adb", "root"},
   {"adb", "shell", "mount", "-t", "tmpfs", "tmpfs", system},
   {"adb", "shell", "mv", data + "/*", system},
   {"adb", "shell", "chcon", "u:object_r:system_file:s0", system + "/*"},
}

func main() {
   for _, command := range commands {
      cmd := exec.Command(command[0], command[1:]...)
      cmd.Stderr = os.Stderr
      cmd.Stdout = os.Stdout
      err := cmd.Run()
      if err != nil {
         panic(err)
      }
   }
}
