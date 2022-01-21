# Docs

## How to determine required features?

Use a command like this:

~~~
aapt dump badging file.apk
~~~

or check the `cmd/badging` folder.

## How to get Protocol Buffer fields?

Check Google Play Store (`com.android.vending`) with JADX, with the last
working version (2016):

~~~
versionCode='80441400' versionName='6.1.14'
~~~

- https://apkmirror.com/apk/google-inc/google-play-store
- https://github.com/whyorean/GPlayApi/blob/master/src/main/proto/GooglePlay.proto

## How to get TV apps?

https://play.google.com/store/apps/details?id=com.iqiyi.i18n.tv

In general, you can probably just get the regular app instead:

https://play.google.com/store/apps/details?id=com.iqiyi.i18n

## How to install split APK?

Bash:

~~~
adb install-multiple *.apk
~~~

PowerShell:

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

## Will you add features?

I created this tool for a very limited use (getting latest version string for
`com.google.android.youtube`). I am happy to fix bugs or add support for apps,
but I am not really interested in adding any features at this time. However if
you want to create an issue to make a suggestion, go ahead, but know that it
will probably not be implemented. However the module itself is open source [1],
so people can easily make their own tools using the module.

If you insist that some feature be implemented, I am willing to implement
features for people who are willing to make a donation to me. If that is your
situation, make sure to mention that in any communication.

1. https://godocs.io/github.com/89z/googleplay
