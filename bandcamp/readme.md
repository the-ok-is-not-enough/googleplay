# Bandcamp

- https://apkpure.com/bandcamp/com.bandcamp.android
- https://github.com/mitmproxy/mitmproxy/issues/4838

## apktool

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
apktool b Bandcamp_v2.4.11_apkpure.com
~~~

Sign:

~~~
jarsigner dist\Bandcamp_v2.4.11_apkpure.com.apk mykey
~~~

- https://netspi.com/blog/technical/four-ways-bypass-android-ssl-verification-certificate-pinning
- https://stackoverflow.com/questions/52862256/charles-proxy-for-mobile-apps

## adb

> adding client certificates to the system-wide trust store, which is by default
> trusted by all apps

~~~
C:\Users\Steven\.mitmproxy
~~~

Get first line:

~~~
openssl x509 -inform PEM -subject_hash_old -in mitmproxy-ca.pem
~~~

Then copy file with new name:

~~~
c8750f0d.0
~~~

Then append to file:

~~~
openssl x509 -inform PEM -text -in mitmproxy-ca.pem >> c8750f0d.0
~~~

Then push new file:

~~~
adb shell mount -o rw,remount,rw /system
adb push c8750f0d.0 /system/etc/security/cacerts/
adb shell mount -o ro,remount,ro /system
adb reboot
~~~

- https://blog.nviso.eu/2017/12/22/intercepting-https-traffic-from-apps-on-android-7-using-magisk-burp
- https://docs.mitmproxy.org/stable/concepts-certificates
- https://stackoverflow.com/questions/44942851/install-user-certificate-via-adb
