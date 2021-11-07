package googleplay

import (
   "net/url"
   "time"
)

const (
   Sleep = 16 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
   deviceConfiguration = 1
   glEsVersion = 8
   hasFiveWayNavigation = 6
   hasHardKeyboard = 5
   keyboard = 2
   nativePlatform = 11
   navigation = 3
   screenDensity = 7
   screenLayout = 4
   systemAvailableFeature = 10
   touchScreen = 1
)

type Auth struct {
   url.Values
}
