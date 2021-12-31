package gplayapi

import (
   "bytes"
   "fmt"
   "google.golang.org/protobuf/proto"
   "gplayapi/gpproto"
   "net/http"
)

func (client *GooglePlayClient) GenerateGsfID() (string, error) {
   req := &gpproto.AndroidCheckinRequest{
      Version:             ptrInt32(3),
      Checkin: &gpproto.AndroidCheckinProto{
         Build:           client.DeviceInfo.Build.AndroidBuildProto,
      },
      DeviceConfiguration: client.DeviceInfo.GetDeviceConfigProto(),
   }
   checkInResp, err := client.checkIn(req)
   if err != nil {
      return "", err
   }
   client.AuthData.GsfID = fmt.Sprintf("%x", checkInResp.GetAndroidId())
   return client.AuthData.GsfID, nil
}

func (client *GooglePlayClient) checkIn(req *gpproto.AndroidCheckinRequest) (resp *gpproto.AndroidCheckinResponse, err error) {
   b, err := proto.Marshal(req)
   if err != nil {
      return
   }
   r, _ := http.NewRequest("POST", UrlBase + "/checkin", bytes.NewReader(b))
   r.Header.Set("Content-Type", "application/x-protobuffer")
   b, _, err = doReq(r)
   if err != nil {
      return
   }
   resp = &gpproto.AndroidCheckinResponse{}
   err = proto.Unmarshal(b, resp)
   return
}

var Pixel3a = &DeviceInfo{
   Build: &DeviceBuildInfo{
      AndroidBuildProto: &gpproto.AndroidBuildProto{
         Radio:          ptrStr("g670-00011-190411-B-5457439"),
         Bootloader:     ptrStr("b4s4-0.1-5613380"),
         SdkVersion:     ptrInt32(29),
      },
   },
   Features:        []string{
      "android.hardware.faketouch",
      "android.hardware.screen.portrait",
   },
   GLVersion:       196610,
   Platforms:    []string{"arm64-v8a", "armeabi-v7a", "armeabi"},
   Screen: &DeviceInfoScreen{},
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

func (i *DeviceInfo) GetDeviceConfigProto() *gpproto.DeviceConfigurationProto {
   return &gpproto.DeviceConfigurationProto{
      GlEsVersion:            &i.GLVersion,
      GlExtension:            i.GLExtensions,
      HasFiveWayNavigation:   ptrBool(false),
      HasHardKeyboard:        ptrBool(false),
      Keyboard:               &i.Keyboard,
      Navigation:             &i.Navigation,
      ScreenDensity:          &i.Screen.Density,
      ScreenLayout:           &i.ScreenLayout,
      SystemAvailableFeature: i.Features,
      SystemSharedLibrary:    i.SharedLibraries,
      TouchScreen:            &i.TouchScreen,
      DeviceFeature:          i.GetDeviceFeatures(),
   }
}

func (i *DeviceInfo) GetDeviceFeatures() (ret []*gpproto.DeviceFeature) {
   var int0 int32
   for _, f := range i.Features {
      name := f
      ret = append(ret, &gpproto.DeviceFeature{Name: &name, Value: &int0})
   }
   return ret
}
