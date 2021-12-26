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

Create a file containing Token (`aas_et`) for future requests:

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
&{InstallationSize:38.717 MB NumDownloads:10.812 B Offer:0.00 USD Title:YouTube
VersionCode:1524493760 VersionString:16.46.37}
~~~

Purchase app. Only needs to be done once per Google account:

~~~
googleplay -a com.google.android.youtube -purchase
~~~

Download APK. If APK is split, all pieces will be downloaded:

~~~
googleplay -a com.google.android.youtube -v 1524493760
~~~

## Repo

https://github.com/89z/googleplay

[1]://f002.backblazeb2.com/file/ql8mlh/Spooked+%28The+Office%29.mp4
