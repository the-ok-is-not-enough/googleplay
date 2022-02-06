# Man In The Middle

## Martian

- https://github.com/google/martian#har-logging
- https://github.com/google/martian#intercepting-https-requests-and-responses

## MITM Proxy

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

To exit, press `q`, then `y`. To capture HTTPS, open Google Chrome on Virtual
Device, and browse to <http://mitm.it>. Click on the Android certificate. Under
`Certificate name` enter `MITM`, then click `OK`. Then browse to:

~~~
https://example.com
~~~

Disable compression:

~~~
set anticomp true
~~~

1. https://mitmproxy.org/downloads

## Python

~~~ps1
$env:HTTPS_PROXY = 'http://127.0.0.1:8080'
$env:REQUESTS_CA_BUNDLE = 'C:\Users\Steven\.mitmproxy\mitmproxy-ca.pem'
$env:SSL_CERT_FILE
~~~
