from . import googleplay_pb2
from base64 import b64decode, urlsafe_b64encode
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.asymmetric.utils import encode_dss_signature
from cryptography.hazmat.primitives.serialization import load_der_public_key
from datetime import datetime
from os import path
from re import match
from sys import version_info
from time import time
import requests

BASE = "https://android.clients.google.com/"
FDFE = BASE + "fdfe/"

AUTH_URL = BASE + "auth"
CHECKIN_URL = BASE + "checkin"
CONTENT_TYPE_PROTO = "application/x-protobuf"
UPLOAD_URL = FDFE + "uploadDeviceConfig"
ssl_verify = True

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
        self.deviceBuilder = DeviceBuilder(device_codename)


    def getHeaders(self, upload_fields=False):
      agent = "Android-Finsky/8.4.19.V-all [0] [FP] 175058788 ("
      # NEED THIS:
      agent += 'versionCode=' + selfDevice.get('vending.version') + ','
      # NEED THIS:
      agent += 'sdk=' + selfDevice.get('build.version.sdk_int') + ','
      agent += ')'
      headers = {'User-Agent': agent}
      if self.gsfId is not None:
          headers["X-DFE-Device-Id"] = "{0:x}".format(self.gsfId)
      if self.authSubToken is not None:
          headers["Authorization"] = "GoogleLogin auth=%s" % self.authSubToken
      if self.device_config_token is not None:
          headers["X-DFE-Device-Config-Token"] = self.device_config_token
      if self.deviceCheckinConsistencyToken is not None:
          headers["X-DFE-Device-Checkin-Consistency-Token"] = self.deviceCheckinConsistencyToken
      if self.dfeCookie is not None:
          headers["X-DFE-Cookie"] = self.dfeCookie
      return headers

    def checkin(self, email, ac2dmToken):
      headers = self.getHeaders()
      headers["Content-Type"] = CONTENT_TYPE_PROTO
      request = self.deviceBuilder.getAndroidCheckinRequest()
      stringRequest = request.SerializeToString()
      res = requests.post(CHECKIN_URL, data=stringRequest,
           headers=headers, verify=ssl_verify,
           proxies=self.proxies_config)
      print(res.status_code, res.url)
      response = googleplay_pb2.AndroidCheckinResponse()
      response.ParseFromString(res.content)
      self.deviceCheckinConsistencyToken = response.deviceCheckinConsistencyToken
      return response.androidId

    def uploadDeviceConfig(self):
      upload = googleplay_pb2.UploadDeviceConfigRequest()
      # NEED THIS:
      upload.deviceConfiguration.CopyFrom(self.deviceBuilder.getDeviceConfig())
      headers = self.getHeaders(upload_fields=True)
      stringRequest = upload.SerializeToString()
      response = requests.post(UPLOAD_URL, data=stringRequest,
            headers=headers,
            verify=ssl_verify,
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
       headers = self.deviceBuilder.getAuthHeaders(self.gsfId)
       headers['app'] = 'com.android.vending'
       response = requests.post(AUTH_URL,
            data=params,
            headers=headers,
            verify=ssl_verify,
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

selfDevice = {
   'build.device': 'sailfish',
   'build.fingerprint': 'google/sailfish/sailfish:8.1.0/OPM4.171019.016.B1/4720843:user/release-keys',
   'build.hardware': 'sailfish',
   'build.id': 'OPM2.171019.029.B1',
   'build.manufacturer': 'Google',
   'build.model': 'Pixel',
   'build.product': 'sailfish',
   'build.radio': 'unknown',
   'build.version.release': '8.1.0',
   'build.version.sdk_int': '27',
   'celloperator': '311480',
   'client': 'android-google',
   'gl.version': '196610',
   'gsf.version': '12685052',
   'hasfivewaynavigation': 'false',
   'hashardkeyboard': 'false',
   'keyboard': '1',
   'navigation': '1',
   'platforms': 'arm64-v8a,armeabi-v7a,armeabi',
   'roaming': 'mobile-notroaming',
   'screen.density': '420',
   'screen.height': '1080',
   'screen.width': '1794',
   'screenlayout': '2',
   'simoperator': '311480',
   'timezone': 'America/Los_Angeles',
   'touchscreen': '3',
   'userreadablename': 'Google Pixel (api27)',
   'vending.version': '81031200',
   'vending.versionstring': '10.3.12-all [0] [PR] 198814133',
   # NEED THIS
   'features': ','.join([
      'android.hardware.telephony.cdma',
      'android.hardware.telephony.gsm',
      'android.hardware.touchscreen',
      'android.hardware.touchscreen.multitouch',
      'android.hardware.touchscreen.multitouch.distinct',
      'android.hardware.vulkan.compute',
      'android.hardware.vulkan.level',
      'android.hardware.vulkan.version',
      'android.hardware.wifi',
      'android.hardware.wifi.direct',
   ])
}

class DeviceBuilder(object):

    def __init__(self, device):
      self.device = {}
    

    def getAuthHeaders(self, gsfid):
      headers = {
         "User-Agent": ("GoogleAuth/1.4 ({device} {id})").format(
            device=selfDevice.get('build.device'),
            id=selfDevice.get('build.id')
         )
      }
      if gsfid is not None:
            headers['device'] = "{0:x}".format(gsfid)
      return headers
    
    # NEED THIS
    def getAndroidCheckinRequest(self):
      request = googleplay_pb2.AndroidCheckinRequest()
      request.id = 0
      request.checkin.CopyFrom(self.getAndroidCheckin())
      request.version = 3
      request.fragment = 0
      return request

    def getDeviceConfig(self):
      featureList = selfDevice['features'].split(",")
      platforms = selfDevice['platforms'].split(",")
      hasFiveWayNavigation = (selfDevice['hasfivewaynavigation'] == 'true')
      hasHardKeyboard = (selfDevice['hashardkeyboard'] == 'true')
      deviceConfig = googleplay_pb2.DeviceConfigurationProto()
      deviceConfig.touchScreen = int(selfDevice['touchscreen'])
      deviceConfig.keyboard = int(selfDevice['keyboard'])
      deviceConfig.navigation = int(selfDevice['navigation'])
      deviceConfig.screenLayout = int(selfDevice['screenlayout'])
      deviceConfig.hasHardKeyboard = hasHardKeyboard
      deviceConfig.hasFiveWayNavigation = hasFiveWayNavigation
      deviceConfig.screenDensity = int(selfDevice['screen.density'])
      deviceConfig.screenWidth = int(selfDevice['screen.width'])
      deviceConfig.screenHeight = int(selfDevice['screen.height'])
      deviceConfig.glEsVersion = int(selfDevice['gl.version'])
      for x in platforms:
          deviceConfig.nativePlatform.append(x)
      for x in featureList:
          deviceConfig.systemAvailableFeature.append(x)
      return deviceConfig

    # NEED THIS
    def getAndroidCheckin(self):
        androidCheckin = googleplay_pb2.AndroidCheckinProto()
        androidCheckin.lastCheckinMsec = 0
        androidCheckin.cellOperator = selfDevice['celloperator']
        androidCheckin.simOperator = selfDevice['simoperator']
        androidCheckin.roaming = selfDevice['roaming']
        androidCheckin.userNumber = 0
        return androidCheckin
