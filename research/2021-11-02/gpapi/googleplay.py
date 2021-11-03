from . import googleplay_pb2
import requests

BASE = "https://android.clients.google.com/"
FDFE = BASE + "fdfe/"

AUTH_URL = BASE + "auth"
CHECKIN_URL = BASE + "checkin"
CONTENT_TYPE_PROTO = "application/x-protobuf"
UPLOAD_URL = FDFE + "uploadDeviceConfig"

selfDevice = {
   # NEED ALL THESE:
   'gl.version': '196610',
   'hasfivewaynavigation': 'false',
   'hashardkeyboard': 'false',
   'keyboard': '1',
   'navigation': '1',
   'platforms': 'arm64-v8a,armeabi-v7a,armeabi',
   'screen.density': '420',
   'screenlayout': '2',
   'touchscreen': '3',
   'build.version.sdk_int': '27',
   'vending.version': '81031200',
   'features': ','.join([
      'android.hardware.touchscreen',
      'android.hardware.wifi',
   ])
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
      # NEED ALL THESE:
      request = googleplay_pb2.AndroidCheckinRequest()
      request.version = 3
      androidCheckin = googleplay_pb2.AndroidCheckinProto()
      request.checkin.CopyFrom(androidCheckin)
      # b'"\x00p\x03'
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
      upload = googleplay_pb2.UploadDeviceConfigRequest()
      deviceConfig = googleplay_pb2.DeviceConfigurationProto()
      # NEED ALL THESE:
      deviceConfig.glEsVersion = int(selfDevice['gl.version'])
      deviceConfig.hasFiveWayNavigation = (selfDevice['hasfivewaynavigation'] == 'true')
      deviceConfig.hasHardKeyboard = (selfDevice['hashardkeyboard'] == 'true')
      deviceConfig.keyboard = int(selfDevice['keyboard'])
      deviceConfig.navigation = int(selfDevice['navigation'])
      deviceConfig.screenDensity = int(selfDevice['screen.density'])
      deviceConfig.screenLayout = int(selfDevice['screenlayout'])
      deviceConfig.touchScreen = int(selfDevice['touchscreen'])
      for x in selfDevice['features'].split(","):
         deviceConfig.systemAvailableFeature.append(x)
      for x in selfDevice['platforms'].split(","):
         deviceConfig.nativePlatform.append(x)
      upload.deviceConfiguration.CopyFrom(deviceConfig)
      # NEED THIS:
      agent = 'versionCode=' + selfDevice.get('vending.version') + ','
      # NEED THIS:
      agent += 'sdk=' + selfDevice.get('build.version.sdk_int')
      headers = {
         'User-Agent': "Android-Finsky (" + agent + ')'
      }
      # NEED THIS:
      headers["X-DFE-Device-Id"] = "{0:x}".format(self.gsfId)
      # NEED THIS:
      headers["Authorization"] = 'Bearer ' + ac2dmToken
      stringRequest = upload.SerializeToString()
      response = requests.post(
         UPLOAD_URL,
         data=stringRequest,
         headers=headers,
      )
      print(response.status_code, response.url)
