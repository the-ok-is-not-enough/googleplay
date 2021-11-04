# November 3 2021

When creating `X-DFE-Device-ID`, you must use Android API 25 or higher.

- https://github.com/Juby210/gplayapi-go/issues/4
- https://github.com/4thel00z/google-play

This issue is helpful:

https://github.com/EFForg/apkeep/issues/17

Fix:

https://github.com/NoMore201/googleplay-api/pull/153

The answer is here:

https://github.com/NoMore201/googleplay-api/blob/664c399f/gpapi/googleplay.py#L242-L244

~~~
pip install protobuf
pip install requests
~~~

After getting a new device ID, you have to wait about 9 seconds before you try
to use it. This is interesting:

https://gitlab.com/marzzzello/playstoreapi/-/blob/e3328b7b/playstoreapi/googleplay.py#L25

How to get this:

~~~
X-DFE-Device-ID
~~~

AKA Google Service Framework ID. Does the APK have it? No. This looks like the
answer here:

https://github.com/crow-misia/go-push-receiver/blob/main/instanceid.go

~~~
done
language:go android.clients.google.com sdk_version
android.clients.google.com sdk_version -EncryptedPasswd

BadAuthentication
~~~
