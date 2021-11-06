# googleplay

Download APK from Google Play or send API requests

https://godocs.io/github.com/89z/googleplay

## Tool examples

Create a file containing Token (`aas_et`) for future requests:

~~~
googleplay -e EMAIL -p PASSWORD
~~~

Create a file containing `Android_ID` for future requests:

~~~
googleplay -d
~~~

Get app details:

~~~
googleplay -a com.google.android.youtube
~~~

Get APK URL:

~~~
googleplay -a com.google.android.youtube -v 1524094400
~~~

## Module example

~~~go
package main

import (
   "fmt"
   "time"
   gp "github.com/89z/googleplay"
)

func main() {
   tok, err := gp.NewToken("EMAIL", "PASSWORD")
   if err != nil {
      panic(err)
   }
   auth, err := tok.Auth()
   if err != nil {
      panic(err)
   }
   dev, err := gp.NewDevice(gp.DefaultCheckin)
   if err != nil {
      panic(err)
   }
   auth.Upload(dev, gp.DefaultConfig)
   time.Sleep(gp.Sleep)
   del, err := auth.Delivery(dev, "com.google.android.youtube", 1524094400)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%+v\n", del)
}
~~~

## How to get Android JA3?

First install:

https://github.com/emanuele-f/PCAPdroid

Start app, then change from HTTP Server to PCAP File. Then click start, if
prompted to save, choose Downloads. Start Android Chrome and wait for a page to
load. Then stop monitoring, and copy file to computer:

~~~
adb ls /sdcard/Download
adb pull /sdcard/Download/PCAPdroid_22_Oct_15_19_28.pcap
~~~

Then my other package can get you the rest of the way:

https://godocs.io/github.com/89z/parse/crypto

## How to get Android public key?

Get this:

https://apkpure.com/google-play-services/com.google.android.gms

Then extract:

~~~
apktool d com.google.android.gms.apk
~~~

The Android public key modulus length should always be 128, which Base64 encoded
looks like:

~~~
AAAAgA
~~~

So you should be able to search the extracted files for one of these:

~~~
AAAAg
public key available
~~~

Result:

~~~
smali\gnt.smali
320: const-string v1, "no public key available, using default"
321-
322- invoke-interface {v0, v1}, Lalyp;->u(Ljava/lang/String;)V
323-
324- const-string v0, "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwg...
~~~

## Thanks

https://github.com/4thel00z/google-play
