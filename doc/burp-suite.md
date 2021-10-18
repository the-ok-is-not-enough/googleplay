# Burp Suite

First download [1]. The program installs to:

~~~
C:\Users\Steven\AppData\Local\Programs\BurpSuiteCommunity
~~~

Start the program, and choose **Temporary project**, then click **Next**. Choose
**Use Burp defaults**, and click **Start Burp**. Then click **Proxy**,
**Intercept**. Then click **Intercept is on** to turn it off. Intercept pauses
every request until you click **Forward**, so you probably dont want that. Then
click **Options**. Under **Proxy Listeners**, the default interface should be:

~~~
127.0.0.1:8080
~~~

Enter that address and port into Android Studio. Then click **HTTP history**.
In the emulator, open Google Chrome and browse to:

~~~
http://example.com
~~~

1. https://portswigger.net/burp/communitydownload

## HTTPS

Click **Proxy**, then **Options**. Then click **export CA certificate**. Then
click **Export Certificate in DER format**. Then click **Next**. Save the file
as `cacert.cer`. Drag the file to emulator home screen to copy. Then in
emulator, open Settings App. Then click **Security**. Then click **Install from
SD card**. Then click **Android SDK**. Then click **Download**. Then click
`cacert.cer`. Under **Certificate name**, enter BURP, and click **OK**. Then in
Burp Suite, click **Proxy** and **HTTP history**. In the emulator, open Google
Chrome and browse to:

~~~
https://example.com
~~~

- https://portswigger.net/support/configuring-an-android-device-to-work-with-burp
- https://portswigger.net/support/installing-burp-suites-ca-certificate-in-an-android-device

## Is Burp Suite open source?

No.
