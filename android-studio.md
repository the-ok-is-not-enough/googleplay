# Android Studio

First download the package [1]. Start the program, and click **More Actions**,
**AVD Manager**, **Create Virtual Device**. On the **Select Hardware** screen,
click **Next**. On the **System Image** screen, click **x86 Images**. Find this
entry:

Release Name | API Level | ABI | Target
-------------|-----------|-----|------------
Nougat       | 24        | x86 | Google APIs

Note that in some cases you might need a higher version. For example, I believe
Google Play requires API 25. Once you have chosen, click **Download**. Then
click **Next**. On the **Android Virtual Device** screen, click **Finish**. On
the **Your Virtual Devices** screen, click **Launch this AVD in the emulator**.

1. https://developer.android.com/studio#downloads

## Proxy

If you need to configure a proxy, in the emulator click **More**. On the
**Extended Controls** screen, click **Settings**, **Proxy**. Uncheck **Use
Android Studio HTTP proxy settings**. Click **Manual proxy configuration**.
Enter **Host name** and **Port number** as determined by the proxy program you
are using. Click **Apply**, and you should see **Proxy status Success**.

## Uninstall

The two big folders are here:

~~~
C:\Users\Steven\.android
C:\Users\Steven\AppData\Local\Android
~~~
