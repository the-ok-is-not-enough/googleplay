# November 4 2021

This issue is helpful:

https://github.com/EFForg/apkeep/issues/17

After getting a new device ID, you have to wait about 9 seconds before you try
to use it. This is interesting:

https://gitlab.com/marzzzello/playstoreapi/-/blob/e3328b7b/playstoreapi/googleplay.py#L25

~~~
done
language:go android.clients.google.com sdk_version
android.clients.google.com sdk_version -EncryptedPasswd

BadAuthentication
~~~

## To do

- Decode Details response
- Encode method for Device ID
- Encode method for Token
