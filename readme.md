# GooglePlay

> Fear plays an interesting role in our lives. How dare we let it motivate us?
> How dare we let it into our decision-making, into our livelihoods, into our
> relationships?
>
> It's funny, isn't it, we take a day a year to dress up in costume and
> celebrate fear?
>
> [Spooked (The Office) (2011)][1]

Download APK from Google Play or send API requests

[1]://f002.backblazeb2.com/file/ql8mlh/Spooked+%28The+Office%29.mp4

## How to install?

Download Go from here:

https://golang.org/dl

and extract archive. Then download GooglePlay:

https://github.com/89z/googleplay/archive/refs/heads/master.zip

and extract archive. Then navigate to `googleplay-master/cmd/googleplay`, and
enter:

~~~
go build
~~~

## Tool examples

If you are outside the United States, you might need to create an
[App Password][2]. Create a file containing Token (`aas_et`) for future
requests:

~~~
googleplay -e EMAIL -p PASSWORD
~~~

Create a file containing `Android_ID` (GSF ID) for future requests:

~~~
googleplay -d
~~~

Get app details:

~~~
PS C:\> googleplay -a com.google.android.youtube
Title: YouTube
VersionString: 16.49.37
VersionCode: 1524886976
NumDownloads: 11.025 B
Size: 38.905 MB
Offer: 0.00 USD
~~~

Purchase app. Only needs to be done once per Google account:

~~~
googleplay -a com.google.android.youtube -purchase
~~~

Download APK. You need to specify any valid version code. The latest code is
provided by the previous details command. If APK is split, all pieces will be
downloaded:

~~~
googleplay -a com.google.android.youtube -v 1524886976
~~~

[2]://support.google.com/accounts/answer/185833

## Repo

https://github.com/89z/googleplay
