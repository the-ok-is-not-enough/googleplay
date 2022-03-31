# Deep linking

<https://wikipedia.org/wiki/Mobile_deep_linking>

If you can get the APK, then you can open it in JADX [1]:

~~~
jadx-gui com.pinterest-10098030.apk
~~~

and check the `Androidmanifest.xml` file:

~~~xml
<intent-filter android:autoVerify="true">
   <action android:name="android.nfc.action.NDEF_DISCOVERED"/>
   <action android:name="android.intent.action.VIEW"/>
   <category android:name="android.intent.category.DEFAULT"/>
   <category android:name="android.intent.category.BROWSABLE"/>
   <data android:scheme="https"/>
   <data android:scheme="http"/>
   <data android:host="www.pinterest.com"/>
   <data android:host="post.pinterest.com"/>
   <data android:host="pin.it"/>
   <!-- ... -->
</intent-filter>
~~~

So only link with those host will get noticed by the app. In addition, some apps also filter on the path:

~~~xml
<intent-filter android:label="@string/app_name" android:autoVerify="true">
   <data android:scheme="https" android:host="www.pbs.org"/>
   <data android:pathPattern="/video/.*/"/>
</intent-filter>
~~~

Finally, if you have `adb`, you can use it like this:

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://abc.com/shows/greys-anatomy/episode-guide/season-18/12-the-makings-of-you
~~~

1. https://github.com/skylot/jadx
