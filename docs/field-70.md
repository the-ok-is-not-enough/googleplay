# Field 70

Making a request like this:

~~~
GET /fdfe/details?doc=com.gameo2.CatColoring HTTP/1.1
Host: android.clients.google.com
Authorization: Bearer ya29.a0ARrdaM_YOSds-EApEzdd8FP_0Q62kHbkDOAtv-mLk2Wkeh-6W...
X-DFE-Device-ID: 3a1c36f387b...
~~~

In the response we get this:

~~~json
{
  "1": {
    "2": {
      "4": {
        "13": {
          "1": {
            "1": "GameO2",
            "3": 2,
            "4": "1.2",
            "9": 34655364,
            "13": "500+",
            "70": 740,
          }
        }
      }
    }
  }
}
~~~

Where `13` is `numDownloads`, and `70` seems to be the actual number. Does it
match up with the other apps? Yes it does:

~~~
com.elitegames.hillcarracing:
{VersionCode:2 NumDownloads:10+ Something:39}

com.ANOOGAMES.WordsSearchPremium:
{VersionCode:5 NumDownloads:100+ Something:158}

com.gameo2.CatColoring:
{VersionCode:2 NumDownloads:500+ Something:740}

com.fbig.crosswords:
{VersionCode:219 NumDownloads:1,000+ Something:2316}

org.friends.dragonli:
{VersionCode:10 NumDownloads:10,000+ Something:46140}

com.techorus.HyperSodaGeyser:
{VersionCode:25 NumDownloads:50,000+ Something:63459}

com.ketchapp.jabbycat:
{VersionCode:125575 NumDownloads:100,000+ Something:216888}

jp.co.ponos.nyanko_odorobo:
{VersionCode:120 NumDownloads:1,000,000+ Something:1191290}

com.peacocktv.peacockandroid:
{VersionCode:121021022 NumDownloads:10,000,000+ Something:12105026}

com.reddit.frontpage:
{VersionCode:387663 NumDownloads:50,000,000+ Something:83264717}

com.discord:
{VersionCode:102017 NumDownloads:100,000,000+ Something:282988890}

com.netflix.mediaclient:
{VersionCode:40080 NumDownloads:1,000,000,000+ Something:1673255914}

com.google.android.youtube:
{VersionCode:1524493760 NumDownloads:10,000,000,000+ Something:10788627915}
~~~

So, what is this field called? I check online, and every single ProtoBuf file is
old:

1. <https://github.com/a-sarja/Google_App_Downloader/issues/2>
1. https://github.com/EFForg/rs-google-play/issues/5
1. https://github.com/GoneToneStudio/node-google-play-api/issues/14
1. https://github.com/Huma123456/gapps/issues/2
1. https://github.com/Juby210/gplayapi-go/issues/5
1. https://github.com/Mayankagg44/App-crawler/issues/6
1. https://github.com/MrChota/gplayapi/issues/2
1. https://github.com/anod/AppWatcher/issues/92
1. https://github.com/fahadakbar24/google-play-api/issues/4
1. https://github.com/kagasu/GooglePlayStoreApi/issues/15
1. https://github.com/opengapps/apkcrawler/issues/75

Then I checked `com.android.vending` with AndroGuard, with the last working
version (2016):

~~~
versionCode='80441400' versionName='6.1.14'
~~~

This version:

~~~
versionCode='80621000' versionName='6.2.10.A-all [0] 2590673'
~~~

still has the fields, but the code looks different:

~~~java
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.c(1, this.a);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.b(2, this.c);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.b(3, this.e);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.c(4, this.g);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.b(5, this.i);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.b(6, this.k);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.c(7, this.m);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.c(8, this.o);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.c(9, this.q);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(10, this.s);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(11, this.t);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.e(12, this.u);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.e(13, this.w);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(14, this.y);
v0 += (com.google.protobuf.nano.CodedOutputByteBufferNano.c(15) + 1);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(17, this.B);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(18, this.C);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.c(19, this.D);
v0 += (com.google.protobuf.nano.CodedOutputByteBufferNano.c(20) + 1);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(21, this.H);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.b(23, this.I);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(24, this.K);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.b(25, this.L);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(26, this.N);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(27, this.O);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(28, this.P);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(29, this.Q);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.b(30, this.R);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(31, this.T);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(32, this.U);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(33, this.V);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(34, this.W);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(35, this.X);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(36, this.Y);
v0 += (com.google.protobuf.nano.CodedOutputByteBufferNano.c(38) + 1);
v0 += com.google.protobuf.nano.CodedOutputByteBufferNano.d(39, this.ab);
~~~