from . import googleplay_pb2
import requests

BASE = "https://android.clients.google.com/"
FDFE = BASE + "fdfe/"

AUTH_URL = BASE + "auth"
CHECKIN_URL = BASE + "checkin"
CONTENT_TYPE_PROTO = "application/x-protobuf"
UPLOAD_URL = FDFE + "uploadDeviceConfig"

dev = {
   # NEED ALL THESE:
   'build.version.sdk_int': '27',
   'features': 'android.hardware.touchscreen,android.hardware.wifi',
   'gl.version': '196610',
   'hasfivewaynavigation': 'false',
   'hashardkeyboard': 'false',
   'keyboard': '1',
   'navigation': '1',
   'platforms': 'arm64-v8a,armeabi-v7a,armeabi',
   'screen.density': '420',
   'screenlayout': '2',
   'touchscreen': '3',
   'vending.version': '81031200',
}

class LoginError(Exception):

    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)

class GooglePlayAPI(object):

   def __init__(self, locale="en_US", timezone="UTC", device_codename="bacon"):
        self.authSubToken = None
        self.gsfId = None
        self.dfeCookie = None

   def checkin(self, email, ac2dmToken):
      request = googleplay_pb2.AndroidCheckinRequest()
      request.version = 3
      androidCheckin = googleplay_pb2.AndroidCheckinProto()
      request.checkin.CopyFrom(androidCheckin)
      stringRequest = request.SerializeToString()
      res = requests.post(
         CHECKIN_URL,
         data=stringRequest,
         headers = {"Content-Type": CONTENT_TYPE_PROTO},
      )
      print(res.status_code, res.url)
      response = googleplay_pb2.AndroidCheckinResponse()
      response.ParseFromString(res.content)
      return response.androidId

   def uploadDeviceConfig(self, ac2dmToken):
      headers = {
         "Authorization": 'Bearer ' + ac2dmToken,
         "X-DFE-Device-Id": "{0:x}".format(self.gsfId),
         'User-Agent': 'Android-Finsky (versionCode=' + dev.get('vending.version') + ',sdk=' + dev.get('build.version.sdk_int') + ')'
      }
      deviceConfig = googleplay_pb2.DeviceConfigurationProto()
      deviceConfig.glEsVersion = int(dev['gl.version'])
      deviceConfig.hasFiveWayNavigation = (dev['hasfivewaynavigation'] == 'true')
      deviceConfig.hasHardKeyboard = (dev['hashardkeyboard'] == 'true')
      deviceConfig.keyboard = int(dev['keyboard'])
      deviceConfig.navigation = int(dev['navigation'])
      deviceConfig.screenDensity = int(dev['screen.density'])
      deviceConfig.screenLayout = int(dev['screenlayout'])
      deviceConfig.touchScreen = int(dev['touchscreen'])
      for x in dev['features'].split(","):
         deviceConfig.systemAvailableFeature.append(x)
      for x in dev['platforms'].split(","):
         deviceConfig.nativePlatform.append(x)
      upload = googleplay_pb2.UploadDeviceConfigRequest()
      upload.deviceConfiguration.CopyFrom(deviceConfig)
      stringRequest = upload.SerializeToString()
      response = requests.post(
         UPLOAD_URL,
         data=stringRequest,
         headers=headers,
      )
      print(response.status_code, response.url)
