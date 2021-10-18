# APK Tool

~~~
$env:PATH = 'D:\Desktop\jdk-17+35\bin'

java -jar apktool_2.6.0.jar d org.metabrainz.android.apk

java -jar apktool_2.6.0.jar b org.metabrainz.android -o app_patched.apk `
--use-aapt2

keytool -genkey -alias keys -keystore keys -keyalg DSA

jarsigner -verbose -keystore keys app_patched.apk keys
~~~

- https://bugs.openjdk.java.net/browse/JDK-8212111
- https://github.com/iBotPeaches/Apktool/issues/1978
- https://github.com/iBotPeaches/Apktool/issues/731
- https://ibotpeaches.github.io/Apktool/documentation
- https://stackoverflow.com/questions/52862256/charles-proxy-for-mobile-apps
