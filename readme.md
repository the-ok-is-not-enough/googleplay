# googleplay

Download APK from Google Play or send API requests

https://godocs.io/github.com/89z/googleplay

## Tool examples

Create a file containing Token (`aas_et`) for future requests:

~~~
googleplay -e EMAIL -p PASSWORD
~~~

Create a file containing `Android_ID` for future requests:

~~~
googleplay -d
~~~

Get app details:

~~~
googleplay -a com.google.android.youtube
~~~

Purchase app:

~~~
googleplay -a com.google.android.youtube -purchase
~~~

Get APK URL:

~~~
googleplay -a com.google.android.youtube -v 1524094400
~~~

## Module example

See `cmd` folder:

https://github.com/89z/googleplay/tree/master/cmd
