# Magisk

Download this and extract:

https://github.com/shakalaca/MagiskOnEmulator/archive/refs/heads/master.zip

Then download APK to extracted folder:

https://github.com/topjohnwu/Magisk/releases

Rename to `magisk.zip`. Then copy `ramdisk.img` to extracted folder from here:

~~~
C:\Users\Steven\AppData\Local\Android\Sdk\system-images\android-26\google_apis\x86
~~~

Then start AVD. Then, run `patch.bat`.

https://github.com/shakalaca/MagiskOnEmulator/issues/49

Then replace original `ramdisk.img` with modified version. Then Power off and
Cold Boot the device.
