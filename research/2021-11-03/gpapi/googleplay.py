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
   def __init__(self, locale="en_US", timezone="UTC", device_codename="bacon",
                 proxies_config=None):
        self.authSubToken = None
        self.gsfId = None
        self.device_config_token = None
        self.deviceCheckinConsistencyToken = None
        self.dfeCookie = None
        self.proxies_config = proxies_config


   def checkin(self, email, ac2dmToken):
      headers = {}
      headers["Content-Type"] = CONTENT_TYPE_PROTO
      request = googleplay_pb2.AndroidCheckinRequest()
      request.id = 0
      androidCheckin = googleplay_pb2.AndroidCheckinProto()
      request.checkin.CopyFrom(androidCheckin)
      request.version = 3
      request.fragment = 0
      stringRequest = request.SerializeToString()
      res = requests.post(
         CHECKIN_URL,
         data=stringRequest,
         headers=headers,
         proxies=self.proxies_config
      )
      print(res.status_code, res.url)
      response = googleplay_pb2.AndroidCheckinResponse()
      response.ParseFromString(res.content)
      self.deviceCheckinConsistencyToken = response.deviceCheckinConsistencyToken
      return response.androidId

   def uploadDeviceConfig(self):
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
      agent = "Android-Finsky/8.4.19.V-all [0] [FP] 175058788 ("
      # NEED THIS:
      agent += 'versionCode=' + selfDevice.get('vending.version') + ','
      # NEED THIS:
      agent += 'sdk=' + selfDevice.get('build.version.sdk_int')
      agent += ')'
      headers = {'User-Agent': agent}
      if self.gsfId is not None:
         headers["X-DFE-Device-Id"] = "{0:x}".format(self.gsfId)
      if self.authSubToken is not None:
         headers["Authorization"] = "GoogleLogin auth=%s" % self.authSubToken
      stringRequest = upload.SerializeToString()
      response = requests.post(UPLOAD_URL, data=stringRequest,
            headers=headers,
            timeout=60,
            proxies=self.proxies_config)
      print(response.status_code, response.url)
      response = googleplay_pb2.ResponseWrapper.FromString(response.content)
      try:
          if response.payload.HasField('uploadDeviceConfigResponse'):
              self.device_config_token = response.payload.uploadDeviceConfigResponse
              self.device_config_token = self.device_config_token.uploadDeviceConfigToken
      except ValueError:
          pass

   def getAuthSubToken(self, email, passwd, master_token):
      params = {
         "Email": email,
         "EncryptedPasswd": passwd,
         "droidguard_results": "dummy123",
         'service': 'androidmarket',
      }
      if self.gsfId is not None:
         params['androidId'] = "{0:x}".format(self.gsfId)
      params['Token'] = master_token
      params.pop('Email')
      params.pop('EncryptedPasswd')
      response = requests.post(AUTH_URL,
            data=params,
            proxies=self.proxies_config)
      print(response.status_code, response.url)
      data = response.text.split()
      params = {}
      for d in data:
         if "=" not in d:
             continue
         k, v = d.split("=", 1)
         params[k.strip().lower()] = v.strip()
      if "auth" in params:
         second_round_token = params["auth"]
         print(second_round_token)
         self.authSubToken = second_round_token
      elif "error" in params:
          raise LoginError("server says: " + params["error"])
      else:
          raise LoginError("Auth token not found.")
