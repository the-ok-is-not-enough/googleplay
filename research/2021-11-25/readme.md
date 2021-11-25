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
match up with the other apps?
