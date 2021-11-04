# November 3 2021

This issue is helpful:

https://github.com/EFForg/apkeep/issues/17

Fix:

https://github.com/NoMore201/googleplay-api/pull/153

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
