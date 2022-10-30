# Research

## GenyMotion

Install Open GApps. Then start proxy:

~~~
mitmproxy
~~~

then set proxy:

~~~
adb shell settings put global http_proxy 192.168.56.1:8080
~~~

WebView Browser cannot download, so install certificate like this:

~~~
adb push mitmproxy-ca-cert.cer /sdcard/Download
~~~

- https://support.genymotion.com/hc/articles/360002778137-How-to-connect
- https://vizz.github.io/blog/2014/01/08/inspect-your-android-device-networking-with-mitmproxy
