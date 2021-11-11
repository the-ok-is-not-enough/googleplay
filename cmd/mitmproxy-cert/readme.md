# MITM Proxy certificate

~~~
C:\Users\Steven\.mitmproxy
~~~

Get first line:

~~~
outputs the MD5 "hash" of the certificate subject name

openssl x509 -subject_hash_old -in mitmproxy-ca-cert.cer
c8750f0d
~~~

Then copy file with new name:

~~~
c8750f0d.0
~~~

- https://github.com/httptoolkit/httptoolkit-server/blob/master/src/interceptors/android/adb-commands.ts
- https://github.com/mitmproxy/mitmproxy/issues/4838
- <https://packages.msys2.org/package/mingw-w64-x86_64-openssl>
