# Emulator

## APK

To install, drag file to emulator home screen. To uninstall, long press on the
app, and drag to `Uninstall`. To force stop, long press on the app, and drag
to `App info`.

## Android Studio

First download the package [1]. Start the program, and click `More Actions`,
`AVD Manager`, `Create Virtual Device`. On the `Select Hardware` screen, click
`Next`. On the `System Image` screen, click `x86 Images`. Find this entry:

Release Name | API Level | ABI | Target
-------------|-----------|-----|------------------------
Pie          | 28        | x86 | Android 9 (Google APIs)

If the APK you are using supports `x86`, then you can use lower versions down to
API 24. If the APK you are using supports `arm64-v8a`, then you will need to use
API 30 or higher [2]. Once you have chosen, click `Download`. Then click
`Next`. On the `Android Virtual Device` screen, click `Finish`. On the `Your
Virtual Devices` screen, click `Launch this AVD in the emulator`.

1. https://developer.android.com/studio#downloads
2. https://android.stackexchange.com/questions/222094/install-failed

If you need to configure a proxy, in the emulator click `More`. On the
`Extended Controls` screen, click `Settings`, `Proxy`. Uncheck `Use Android
Studio HTTP proxy settings`. Click `Manual proxy configuration`. Enter `Host
name` and `Port number` as determined by the proxy program you are using. Click
`Apply`, and you should see `Proxy status Success`.

## Genymotion

This emulator sucks. Here is Android Studio:

- 0s, launching emulator
- 49s, click Play Store
- 1m, sign in to Play Store
- 1m37s, install Firefox
- 1m51s, open Firefox
- 2m9s, http://example.com finish loading

and Genymotion:

- 0s, launching emulator
- 44s, install Gapps
- 1m41s, restart device
- 2m32s, click Play Store
- 3m, sign in to Play Store
- 4m53s, install Chrome
- 5m37s, open Chrome
- 6m, Chrome crash
