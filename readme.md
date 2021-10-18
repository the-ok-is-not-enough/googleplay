# Google Play

Using Android API 24 fails, but API 25 or higher works. It applies to all
devices, not just Virtual Devices.

https://accounts.google.com/EmbeddedSetup

## Protocol buffer

I have implemented a parser for unknown bytes.

## EncryptedPasswd

I have an implemention of the algorithm in Python, need to convert to Go.

## TLS fingerprint

I have a hand created hello struct, better would be to parse a JA3 string. I
found a JA3 that works. Need to finish the JA3 encoder.

https://github.com/Danny-Dasilva/CycleTLS/issues/39

## Old APKs

Its possible to get an old APK if you have the version code, but what if you
dont? Do they just start at 1 and go up? Is it possible to get a list of the
versions?

https://github.com/Juby210/gplayapi-go/issues/3
