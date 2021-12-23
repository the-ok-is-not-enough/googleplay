package main

import (
   "flag"
   "fmt"
   "github.com/xi2/xz"
   "net/http"
   "os"
   "os/exec"
   "path/filepath"
   "strings"
)

const version = "15.1.10"

func newServer(version, cpu string) string {
   var str strings.Builder
   str.WriteString("https://github.com/frida/frida/releases/download/")
   str.WriteString(version)
   str.WriteString("/frida-server-")
   str.WriteString(version)
   str.WriteString("-android-")
   str.WriteString(cpu)
   str.WriteString(".xz")
   return str.String()
}

func main() {
   _, err := exec.LookPath("frida")
   if err != nil {
      panic("pip install frida-tools")
   }
   var x86_64 bool
   flag.BoolVar(&x86_64, "64", false, "x86_64")
   flag.Parse()
   if flag.NArg() != 1 {
      fmt.Println("frida-script [flags] [package]")
      flag.PrintDefaults()
      return
   }
   pack := flag.Arg(0)
   cache, err := os.UserCacheDir()
   if err != nil {
      panic(err)
   }
   cache = filepath.Join(cache, "googleplay")
   os.Mkdir(cache, os.ModeDir)
   // download script
   cacheScript := filepath.Join(cache, filepath.Base(script))
   if err := downloadScript(cacheScript); err != nil {
      panic(err)
   }
   // download server
   cpu := "x86"
   if x86_64 {
      cpu += "_64"
   }
   server := newServer(version, cpu)
   cacheServer := filepath.Join(cache, stem(server))
   if err := downloadServer(server, cacheServer); err != nil {
      panic(err)
   }
   // commands
   commands := []command{
      run("adb", "root"),
      run("adb", "wait-for-device"),
      run("adb", "push", cacheServer, "/data/app/frida-server"),
      run("adb", "shell", "chmod", "755", "/data/app/frida-server"),
      start("adb", "shell", "/data/app/frida-server"),
      run("frida", "--no-pause", "-U", "-l", cacheScript, "-f", pack),
   }
   for _, cmd := range commands {
      cmd.Stderr = os.Stderr
      cmd.Stdout = os.Stdout
      fmt.Println(cmd.Args)
      cmd.Start()
      if cmd.wait {
         err := cmd.Wait()
         if err != nil {
            panic(err)
         }
      } else {
         defer cmd.Wait()
      }
   }
}

const script =
   "https://raw.githubusercontent.com/httptoolkit/frida-android-unpinning/" +
   "main/frida-script.js"

func downloadScript(dst string) error {
   fmt.Println("Stat", dst)
   _, err := os.Stat(dst)
   if err == nil {
      return nil
   }
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

type command struct {
   *exec.Cmd
   wait bool
}

func run(name string, arg ...string) command {
   var cmd command
   cmd.Cmd = exec.Command(name, arg...)
   cmd.wait = true
   return cmd
}

func start(name string, arg ...string) command {
   var cmd command
   cmd.Cmd = exec.Command(name, arg...)
   return cmd
}
