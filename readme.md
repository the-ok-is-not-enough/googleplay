# googleplay

> Fear plays an interesting role in our lives. How dare we let it motivate us?
> How dare we let it into our decision-making, into our livelihoods, into our
> relationships?
>
> It's funny, isn't it, we take a day a year to dress up in costume and
> celebrate fear?
>
> [Spooked (The Office) (2011)][1]

Download APK from Google Play or send API requests

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
{Offer:{FormattedAmount:$0} Details:{AppDetails:{DeveloperName:Google LLC
VersionCode:1524221376 Version:16.44.32 UploadDate:Nov 2, 2021}}}
~~~

Purchase app. Only needs to be done once per Google account:

~~~
googleplay -a com.google.android.youtube -purchase
~~~

Download APK. If APK is split, all pieces will be downloaded:

~~~
googleplay -a com.google.android.youtube -v 1524221376
~~~

## Module

Docs here:

https://godocs.io/github.com/89z/googleplay

Example here:

https://github.com/89z/googleplay/tree/master/cmd

[1]://f002.backblazeb2.com/file/ql8mlh/Spooked+%28The+Office%29.mp4
