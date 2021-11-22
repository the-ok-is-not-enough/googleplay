package main

import (
   "fmt"
   "github.com/xi2/xz"
   "net/http"
   "os"
   "os/exec"
   "path/filepath"
)

const server =
   "https://github.com/frida/frida/releases/download/" +
   "15.1.10/frida-server-15.1.10-android-x86.xz"

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

func downloadServer(dst string) error {
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
   return command{
      exec.Command(name, arg...), true,
   }
}

func start(name string, arg ...string) command {
   return command{
      exec.Command(name, arg...), false,
   }
}

func main() {
   if len(os.Args) != 2 {
      fmt.Println("frida-script [package]")
      return
   }
   pack := os.Args[1]
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
   cacheServer := filepath.Join(cache, stem(server))
   if err := downloadServer(cacheServer); err != nil {
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
