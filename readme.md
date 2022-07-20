# GooglePlay

> I’m not really sure what the point of this video is but I guess just be
> generous.
>
> Be kind to people because you never know how much they might need it or how
> far it’ll go.
>
> [NakeyJakey (2018)](//youtube.com/watch?v=Cr0UYNKmrUs)

Download APK from Google Play or send API requests

## Issues

During this time I am trying to support myself, so I am requiring payment for
all issues. That might change at a later date. Any issue without payment of at
least 19 USD will be closed immediately. Payment can be made to the PayPal link
on this page, or if you want to use another method, mention that in the issue
text. For business opportunities, contact me:

- srpen6@gmail.com
- Discord srpen6#6983

## How to install?

This module works with Windows, macOS or Linux. First, [download Go][2] and
extract archive. Then [download&nbsp;GooglePlay][3] and extract archive. Then
navigate to `googleplay-master/cmd/googleplay`, and enter:

~~~
go build
~~~

[2]://go.dev/dl
[3]://github.com/89z/googleplay/archive/refs/heads/master.zip

## Tool examples

Before trying these examples, make sure the Google account you are using has
logged into the Play&nbsp;Store at least once before. Also you need to have
accepted the Google Play terms and conditions. Create a file containing token
(`aas_et`) for future requests:

~~~
googleplay -email EMAIL -password PASSWORD
~~~

Create a file containing `X-DFE-Device-ID` (GSF ID) for future requests:

~~~
googleplay -device
~~~

Get app details:

~~~
> googleplay -a com.google.android.youtube
Title: YouTube
Creator: Google LLC
Upload Date: Jul 7, 2022
Version: 17.27.35
Version Code: 1530387904
Num Downloads: 12.18 billion
Installation Size: 48.51 megabyte
File: APK APK APK APK
Offer: 0 USD
~~~

Purchase app. Only needs to be done once per Google account:

~~~
googleplay -a com.google.android.youtube -purchase
~~~

Download APK. You need to specify any valid version code. The latest code is
provided by the previous details command. If APK is split, all pieces will be
downloaded:

~~~
googleplay -a com.google.android.youtube -v 1529992640
~~~

## API

https://github.com/89z/googleplay
