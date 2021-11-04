# November 4 2021

How do we get this:

~~~
Android-Finsky (sdk=27,versionCode=81031200)
~~~

Some commands:

~~~
fail
"Android-Finsky" site:developer.android.com
"Android-Finsky" site:android.com

"User-Agent" "Android-Finsky"

adb shell cat system/build.prop
ro.build.version.sdk=24
ro.build.version.incremental=6696031

sdk=%d
getprop ro.build.version.sdk

versionCode=%d
getprop ro.build.version.incremental
~~~
