# APK Tool

~~~
apktool d Bandcamp_v2.4.11_apkpure.com.apk
~~~

Change `res/xml/network_security_config.xml` to:

~~~xml
<?xml version="1.0" encoding="utf-8"?>
<network-security-config>
   <base-config>
      <trust-anchors>
         <certificates src="user" />
      </trust-anchors>
   </base-config>
</network-security-config>
~~~

Create `C:\Users\Steven\.keystore`:

~~~
keytool -genkey -keyalg DSA
~~~

Create `dist\Bandcamp_v2.4.11_apkpure.com.apk`:

~~~
apktool b Bandcamp_v2.4.11_apkpure.com --use-aapt2
~~~

Sign:

~~~
jarsigner dist\Bandcamp_v2.4.11_apkpure.com.apk mykey
~~~

- https://bugs.openjdk.java.net/browse/JDK-8212111
- https://github.com/iBotPeaches/Apktool/issues/1978
- https://github.com/iBotPeaches/Apktool/issues/731
- https://ibotpeaches.github.io/Apktool/documentation
- https://netspi.com/blog/technical/four-ways-bypass-android-ssl-verification-certificate-pinning
- https://stackoverflow.com/questions/52862256/charles-proxy-for-mobile-apps
