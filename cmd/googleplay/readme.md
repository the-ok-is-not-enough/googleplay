# GooglePlay

~~~
PS C:\> googleplay -a org.wikipedia
POST https://android.clients.google.com/auth
GET https://android.clients.google.com/fdfe/details?doc=org.wikipedia
&{NumDownloads:55.774 M Offer:0.00 USD Size:12.782 MB Title:Wikipedia
VersionCode:50388 VersionString:2.7.50388-r-2021-12-02}

PS C:\> googleplay -a org.wikipedia -v 50388 -o hello
POST https://android.clients.google.com/auth
GET https://android.clients.google.com/fdfe/delivery?doc=org.wikipedia&vc=50388
~~~

https://github.com/89z/googleplay/issues/7
