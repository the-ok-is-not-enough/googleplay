package main

import (
   "fmt"
   "github.com/xi2/xz"
   "net/http"
   "os"
   "os/exec"
   "path/filepath"
   "strings"
)

type command struct {
   *exec.Cmd
   wait bool
}

func newCommand(wait bool, name string, arg ...string) command {
   var cmd command
   cmd.Cmd = exec.Command(name, arg...)
   cmd.Stdout = os.Stdout
   cmd.Stderr = os.Stderr
   cmd.wait = wait
   return cmd
}

func run(name string, arg ...string) command {
   return newCommand(true, name, arg...)
}

func start(name string, arg ...string) command {
   return newCommand(false, name, arg...)
}

func runCommands(server, script, pack string) error {
   commands := []command{
      run("adb", "root"),
      run("adb", "wait-for-device"),
      run("adb", "push", server, "/data/app/frida-server"),
      run("adb", "shell", "chmod", "755", "/data/app/frida-server"),
      start("adb", "shell", "/data/app/frida-server"),
      run("frida", "--no-pause", "-U", "-l", script, "-f", pack),
   }
   for _, command := range commands {
      fmt.Println(command.Args)
      err := command.Start()
      if err != nil {
         return err
      }
      if command.wait {
         err := command.Wait()
         if err != nil {
            return err
         }
      }
   }
   return nil
}

const script =
   "https://raw.githubusercontent.com/httptoolkit/frida-android-unpinning/" +
   "main/frida-script.js"

const version = "15.1.16"

func newServer(version string) string {
   var str strings.Builder
   str.WriteString("https://github.com/frida/frida/releases/download/")
   str.WriteString(version)
   str.WriteString("/frida-server-")
   str.WriteString(version)
   str.WriteString("-android-x86.xz")
   return str.String()
}

func downloadScript(dst string) error {
   fmt.Println("Stat", dst)
   _, err := os.Stat(dst)
   if err != nil {
      fmt.Println("GET", script)
      res, err := http.Get(script)
      if err != nil {
         return err
      }
      defer res.Body.Close()
      file, err := os.Create(dst)
      if err != nil {
         return err
      }
      defer file.Close()
      if _, err := file.ReadFrom(res.Body); err != nil {
         return err
      }
   }
   return nil
}

func downloadServer(server, dst string) error {
   fmt.Println("Stat", dst)
   _, err := os.Stat(dst)
   if err == nil {
      return nil
   }
   fmt.Println("GET", server)
   rHTTP, err := http.Get(server)
   if err != nil {
      return err
   }
   defer rHTTP.Body.Close()
   rXZ, err := xz.NewReader(rHTTP.Body, 0)
   if err != nil {
      return err
   }
   file, err := os.Create(dst)
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := file.ReadFrom(rXZ); err != nil {
      return err
   }
   return nil
}

func stem(s string) string {
   base := filepath.Base(s)
   ext := filepath.Ext(base)
   return base[:len(base)-len(ext)]
}
