# November 25 2021

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
            "16": "Oct 9, 2021",
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

- <https://github.com/a-sarja/Google_App_Downloader/issues/2>
- https://github.com/EFForg/rs-google-play/issues/5
- https://github.com/GoneToneStudio/node-google-play-api/issues/14
- https://github.com/Huma123456/gapps/issues/2
- https://github.com/Juby210/gplayapi-go/issues/5
- https://github.com/Mayankagg44/App-crawler/issues/6
- https://github.com/MrChota/gplayapi/issues/2
- https://github.com/fahadakbar24/google-play-api/issues/4
- https://github.com/kagasu/GooglePlayStoreApi/issues/15
- https://github.com/opengapps/apkcrawler/issues/75

Then I ran this Python script:

~~~py
from androguard.misc import AnalyzeAPK
a,d,dx= AnalyzeAPK('com.android.vending.apk')
f = open('file.java', 'w')

for dd in d:
   for clas in dd.get_classes():
      name = clas.get_name()
      if 'proto' in name:
         print(name, file=f)
         src = dd.get_class(name).get_source()
         print(src, file=f)
~~~

which works with older APKs:
