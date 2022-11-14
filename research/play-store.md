# Google Play Store

If we choose a production build:

API Level | ABI | Target
----------|-----|--------------------------
24        | x86 | Android 7.0 (Google Play)

and start the Google Play device, writable:

~~~
emulator -list-avds
emulator -avd Pixel_3a_XL_API_24 -writable-system
~~~

or not:

~~~
emulator -avd Pixel_3a_XL_API_24
~~~

we cannot run as root:

~~~
> adb root
adbd cannot run as root in production builds
~~~

so create a new device like this:

API Level | ABI | Target
----------|-----|--------------------------
24        | x86 | Android 7.0 (Google APIs)

then go to http://opengapps.org and choose:

Platform | Android | Variant
---------|---------|--------
x86      | 7.0     | pico

and extract:

~~~
Core\vending-x86.tar.lz
~~~

and extract:

~~~
vending-x86\nodpi\priv-app\Phonesky\Phonesky.apk
~~~

and start the Google APIs device:

~~~
emulator -avd Pixel_3a_XL_API_24 -writable-system
~~~

and push:

~~~
adb root

adb remount
adb push Phonesky.apk /system/priv-app
adb reboot
~~~

then start Play Store. If you instead start the Google Play device and pull:

~~~
adb pull /system/priv-app/Phonesky/Phonesky.apk Phonesky.apk
~~~

then start the Google APIs device:

~~~
emulator -avd Pixel_3a_XL_API_24 -writable-system
~~~

and push:

~~~
adb root

adb remount
adb push Phonesky.apk /system/priv-app
adb reboot
~~~

we get this result:

> Google Play services has stopped

So the pulled `Phonesky.apk` is poisoned somehow. Get the `versionCode` from
Google Play device:

~~~
> adb shell dumpsys package com.android.vending | rg versionCode
versionCode=80671500 minSdk=14 targetSdk=23
~~~
