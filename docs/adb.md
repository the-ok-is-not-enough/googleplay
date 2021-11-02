# adb

~~~
C:\Users\89z\.mitmproxy
~~~

Get first line:

~~~
openssl x509 -subject_hash_old -in mitmproxy-ca-cert.cer
~~~

Then copy file with new name:

~~~
c8750f0d.0
~~~

Start emulator:

~~~
adb start-server
emulator -writable-system -avd Pixel_2_API_24
~~~

Push new file:

~~~ps1
adb root
adb remount
adb push c8750f0d.0 /system/etc/security/cacerts
~~~

- https://android.stackexchange.com/questions/242222/cold-boot-snapshot-failed
- https://docs.mitmproxy.org/stable/howto-install-system-trusted-ca-android
- https://github.com/mitmproxy/mitmproxy/issues/4838
- https://stackoverflow.com/questions/36670592/cannot-change-android-system
- https://stackoverflow.com/questions/44942851/install-user-certificate-via-adb
