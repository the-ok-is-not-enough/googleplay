# MITM Proxy

First download [1], then start `mitmproxy.exe`. Address and port should be in
the bottom right corner. Default should be:

~~~
*:8080
~~~

Assuming the above, go to Android Emulator and set proxy:

~~~
127.0.0.1:8080
~~~

Then open Google Chrome on Virtual Device, and browse to:

~~~
http://example.com
~~~

To exit, press `q`, then `y`.

1. https://mitmproxy.org/downloads

## HTTPS

Open Google Chrome on Virtual Device, and browse to <http://mitm.it>. Click on
the Android certificate. Under **Certificate name** enter **MITM**, then click
**OK**. Then browse to:

~~~
https://example.com
~~~

## Certificate pinning

https://github.com/mitmproxy/mitmproxy/issues/4838

but might be able to get it working:

> adding client certificates to the system-wide trust store, which is by default
> trusted by all apps

~~~
C:\Users\Steven\.mitmproxy
~~~

- https://blog.nviso.eu/2017/12/22/intercepting-https-traffic-from-apps-on-android-7-using-magisk-burp
- https://docs.mitmproxy.org/stable/concepts-certificates
- https://stackoverflow.com/questions/44942851/install-user-certificate-via-adb
