# Charles Proxy

First, download the program [1]. Start the program, and click the **Sequence**
tab. Then click **Proxy**, and **Proxy Settings**. The default port should be
`8888`. Assuming that, go to Android Emulator Extended Controls, and enter:

~~~
127.0.0.1:8888
~~~

Then in Android Emulator, open Google Chrome and browse to:

~~~
http://example.com
~~~

1. https://charlesproxy.com/download

## HTTPS

In Android Emulator, open Google Chrome and browse to [1]:

~~~
https://charlesproxy.com/getssl
~~~

Under **Certificate name** enter CHARLES, then click **OK**. Then browse to:

~~~
https://example.com
~~~

Next, in the Charles program under **Sequence**, right click the entry for the
address above, and click **Enable SSL Proxying**. To confirm, click **Proxy**,
then click **SSL Proxying Settings**. Then in Android Emulator, browse to:

~~~
https://example.com
~~~

Then in the Charles program, click **Stop Recording**, then **Stop SSL
Proxying**, then **Clear the current Session**. Then click **Start SSL
Proxying**, then **Start Recording**. Then in Android Emulator, browse to:

~~~
https://example.com
~~~

Then in the Charles program, click **File**, **Quit**. Start the program again.
Then in Android Emulator, browse to:

~~~
https://example.com
~~~

Then with PowerShell, enter:

~~~
Get-NetIPAddress
~~~

Will look like this:

~~~
IPAddress         : 192.168.0.4
InterfaceIndex    : 11
InterfaceAlias    : Ethernet
~~~

Then in Android Emulator, change the proxy to your IP address [2], for example:

~~~
192.168.0.4:8888
~~~

Then browse to:

~~~
https://example.com
~~~

In all cases, we are getting CONNECT Methods instead of GET Methods, which is
not correct. Contact support [3] to request a solution.

1. https://charlesproxy.com/documentation/faqs/ssl-connections-from-within-iphone-applications
2. https://charlesproxy.com/documentation/faqs/using-charles-from-an-iphone
3. https://charlesproxy.com/support/contact

## Is Charles open source?

No.

## Reset configuration

~~~
C:\Users\Steven\AppData\Roaming\Charles
~~~
