# Google Play

APK apps

You can download APKs from different sites [1][2]. If you encounter an XAPK,
this is a Zip format file, so you will need to extract the APK using 7-Zip or
similar. To install, drag file to emulator home screen. To uninstall, long press
on the app, and drag to **Uninstall**. To force stop, long press on the app, and
drag to **App info**.

1. https://apkpure.com
2. https://apkmirror.com

## Java

- https://adoptium.net
- https://jdk.java.net
- https://sap.github.io/SapMachine

## Mozilla Firefox

Open a new tab. Click **Open menu**, **Web Developer**, **Network**. Then click
**Network Settings**, **Persist Logs**. Also check **Disable Cache**. Then
browse to the page you want to capture. Once you are ready, click **Pause**,
then click **Network Settings**, **Save All As HAR**.

## Stack Exchange

Check questions, and ask question if need be:

https://android.stackexchange.com/search?q=certificate+pinning

## EdXposed

- https://github.com/ElderDrivers/EdXposed
- https://github.com/ViRb3/TrustMeAlready

First need to get Magisk:

- https://blog.nviso.eu/2017/12/22/intercepting-https-traffic-from-apps-on-android-7-using-magisk-burp
- https://github.com/topjohnwu/Magisk

## APK MITM

CLI application that automatically removes certificate pinning from Android APK
files:

https://github.com/shroudedcode/apk-mitm/issues/72

## JustTrustMe

https://github.com/Fuzion24/JustTrustMe/issues/61

## Objection

~~~
objection patchapk -s Vimeo_v3.50.1_apkpure.com.apk
~~~

- https://github.com/sensepost/objection/issues/498
- https://joncooperworks.medium.com/disabling-okhttps-ssl-pinning-on-android-bd116aa74e05
- https://levelup.gitconnected.com/bypassing-ssl-pinning-on-android-3c82f5c51d86
- https://stackoverflow.com/questions/44942851/install-user-certificate-via-adb
- https://vavkamil.cz/2019/09/15/how-to-bypass-android-certificate-pinning-and-intercept-ssl-traffic
- https://www.gabriel.urdhr.fr/2021/03/17/frida-disable-certificate-check-on-android
