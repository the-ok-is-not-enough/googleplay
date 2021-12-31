package gplayapi

import (
   "fmt"
   "gplayapi/gpproto"
   "strings"
   "time"
)

func (i *DeviceInfo) GetUserAgent() string {
   params := []string{
      "versionCode=81582300",
      "sdk=28",
   }
   return "Android-Finsky/15.8.23-all [0] [PR] 259261889 (" + strings.Join(params, ",") + ")"
}
var sargo = ptrStr("sargo")

var Pixel3a = &DeviceInfo{
   Build: &DeviceBuildInfo{
      AndroidBuildProto: &gpproto.AndroidBuildProto{
         // need these:
         Radio:          ptrStr("g670-00011-190411-B-5457439"),
         Bootloader:     ptrStr("b4s4-0.1-5613380"),
         // need these?
         Carrier:        ptrStr("google"),
         Device:         sargo,
         Model:          ptrStr("Pixel 3a"),
         Manufacturer:   ptrStr("google"),
         Product:        sargo,
         Id:             ptrStr("PQ3B.190705.003"),
         SdkVersion:     ptrInt32(29),
         Client:         ptrStr("android-google"),
         GoogleServices: ptrInt32(203615028),
         OtaInstalled:   ptrBool(false),
      },
      VersionRelease: 9,
   },
   SimOperator:  "20815",
   Platforms:    []string{"arm64-v8a", "armeabi-v7a", "armeabi"},
   CellOperator: "334050",
   Roaming:      "mobile-notroaming",
   TimeZone:     "America/Mexico_City",
   TouchScreen:  3,
   Keyboard:     1,
   Navigation:   1,
   ScreenLayout: 2,
   Screen: &DeviceInfoScreen{
      Density: 490,
      Width:   1080,
      Height:  2073,
   },
   GLVersion:       196610,
   Features:        []string{
      "android.hardware.faketouch",
      "android.hardware.screen.portrait",
   },
}


type (
	DeviceInfo struct {
		Build           *DeviceBuildInfo
		SimOperator     string
		Platforms       []string
		OtaInstalled    bool
		CellOperator    string
		Roaming         string
		TimeZone        string
		TouchScreen     int32
		Keyboard        int32
		Navigation      int32
		ScreenLayout    int32
		Screen          *DeviceInfoScreen
		GLVersion       int32
		GLExtensions    []string
		SharedLibraries []string
		Features        []string
		Locales         []string
	}

	DeviceBuildInfo struct {
		*gpproto.AndroidBuildProto
		VersionRelease int
	}

	DeviceInfoScreen struct {
		Density int32
		Width   int32
		Height  int32
	}
)

func (i *DeviceInfo) GetAuthUserAgent() string {
   return fmt.Sprintf("GoogleAuth/1.4 %s %s", i.Build.GetDevice(), i.Build.GetId())
}


func (i *DeviceInfo) GenerateAndroidCheckInRequest() *gpproto.AndroidCheckinRequest {
	var int0 int64 = 0
	timestamp := time.Now().Unix()
	i.Build.Timestamp = &timestamp

	return &gpproto.AndroidCheckinRequest{
		Id: &int0,
		Checkin: &gpproto.AndroidCheckinProto{
			Build:           i.Build.AndroidBuildProto,
			LastCheckinMsec: &int0,
			CellOperator:    &i.CellOperator,
			SimOperator:     &i.SimOperator,
			Roaming:         &i.Roaming,
			UserNumber:      ptrInt32(0),
		},
		Locale:              ptrStr("en_GB"),
		TimeZone:            &i.TimeZone,
		Version:             ptrInt32(3),
		DeviceConfiguration: i.GetDeviceConfigProto(),
		Fragment:            ptrInt32(0),
	}
}

func (i *DeviceInfo) GetDeviceConfigProto() *gpproto.DeviceConfigurationProto {
	var mem int64 = 8589935000
	return &gpproto.DeviceConfigurationProto{
		TouchScreen:            &i.TouchScreen,
		Keyboard:               &i.Keyboard,
		Navigation:             &i.Navigation,
		ScreenLayout:           &i.ScreenLayout,
		HasHardKeyboard:        ptrBool(false),
		HasFiveWayNavigation:   ptrBool(false),
		ScreenDensity:          &i.Screen.Density,
		GlEsVersion:            &i.GLVersion,
		SystemSharedLibrary:    i.SharedLibraries,
		SystemAvailableFeature: i.Features,
		ScreenWidth:            &i.Screen.Width,
		ScreenHeight:           &i.Screen.Height,
		SystemSupportedLocale:  i.Locales,
		GlExtension:            i.GLExtensions,
		DeviceClass:            ptrInt32(0),
		LowRamDevice:           ptrInt32(0),
		TotalMemoryBytes:       &mem,
		MaxNumOf_CPUCores:      ptrInt32(8),
		DeviceFeature:          i.GetDeviceFeatures(),
	}
}

func (i *DeviceInfo) GetDeviceFeatures() (ret []*gpproto.DeviceFeature) {
	var int0 int32 = 0
	for _, f := range i.Features {
		name := f
		ret = append(ret, &gpproto.DeviceFeature{Name: &name, Value: &int0})
	}
	return ret
}
