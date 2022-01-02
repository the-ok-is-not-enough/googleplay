package discord

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

const auth = "ya29.a0ARrdaM8eYkSiwt61iquDUJHjzARxAjBFIxWYoKa1C_M_Y2aLAZLz2-5ppwo08s4W3r9LME8Qqb6Cx0U55P8uNtgO_XWozJd6__uTznnI9Tt2ULcv9IqP1-MK-E8cNlR2tDYk7PCf_o3Wru-RtbbUQE1jqorMf24HivgTRr-14Dtn4oXjBMD_xL136zqh1fxyQpcBZ-jdEDvJuqUHDTZrEFoh4VT72RRShLD696PAUqQ8R30Sat6mKOPW6dRiXL6oj1yxze1nZK-K188ZTjxWqxK-XM7yYzj5cAgDu_pVRuqCUy4fUKQ"

var body1 = protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:4, String:""}:protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:6, String:""}:"334050",
protobuf.Tag{Number:7, String:""}:"20815",
protobuf.Tag{Number:8, String:""}:"mobile-notroaming",
protobuf.Tag{Number:9, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:protobuf.Message{protobuf.Tag{Number:10, String:""}:uint64(29),
protobuf.Tag{Number:6, String:""}:"android-google",
protobuf.Tag{Number:14, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"PQ3B.190705.003",
protobuf.Tag{Number:4, String:""}:"g670-00011-190411-B-5457439",
protobuf.Tag{Number:13, String:""}:"sargo",
protobuf.Tag{Number:9, String:""}:"sargo",
protobuf.Tag{Number:11, String:""}:"Pixel 3a",
protobuf.Tag{Number:12, String:""}:"google",
protobuf.Tag{Number:2, String:""}:"sargo",
protobuf.Tag{Number:3, String:""}:"google",
protobuf.Tag{Number:5, String:""}:"b4s4-0.1-5613380",
protobuf.Tag{Number:7, String:""}:uint64(1641084042),
protobuf.Tag{Number:8, String:""}:uint64(203615028)}},
protobuf.Tag{Number:6, String:""}:"en_GB",
protobuf.Tag{Number:12, String:""}:"America/Mexico_City",
protobuf.Tag{Number:14, String:""}:uint64(3),
protobuf.Tag{Number:18, String:""}:protobuf.Message{protobuf.Tag{Number:5, String:""}:uint64(0),
protobuf.Tag{Number:14, String:""}:[]string{"af", "af_ZA", "am", "am_ET", "ar", "ar_EG", "ar_SA", "ar_XB", "as", "ast", "az", "be", "be_BY", "bg", "bg_BG", "bh_IN", "bn", "bs", "ca", "ca_ES", "cs", "cs_CZ", "cy_GB", "da", "da_DK", "de", "de_DE", "el", "el_GR", "en", "en_AU", "en_CA", "en_GB", "en_IN", "en_US", "en_XA", "es", "es_ES", "es_US", "et", "et_EE", "eu", "fa", "fa_IR", "fi", "fi_FI", "fil", "fil_PH", "fr", "fr_CA", "fr_FR", "gl", "gl_ES", "gu", "hi", "hi_IN", "hr", "hr_HR", "hu", "hu_HU", "hy", "in", "in_ID", "is", "it", "it_IT", "iw", "iw_IL", "ja", "ja_JP", "ka", "kab_DZ", "kk", "km", "kn", "ko", "ko_KR", "ky", "lo", "lt", "lt_LT", "lv", "lv_LV", "mk", "ml", "mn", "mr", "ms", "ms_MY", "my", "nb", "nb_NO", "ne", "nl", "nl_NL", "or", "pa", "pa_IN", "pl", "pl_PL", "pt", "pt_BR", "pt_PT", "ro", "ro_RO", "ru", "ru_RU", "sc_IT", "si", "sk", "sk_SK", "sl", "sl_SI", "sq", "sr", "sr_Latn", "sr_RS", "sv", "sv_SE", "sw", "sw_TZ", "ta", "te", "th", "th_TH", "tr", "tr_TR", "uk", "uk_UA", "ur", "uz", "vi", "vi_VN", "zh_CN", "zh_HK", "zh_TW", "zu", "zu_ZA"},
protobuf.Tag{Number:26, String:""}:[]protobuf.Message{protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.audio.low_latency",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.audio.output",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.audio.pro",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"android.hardware.bluetooth"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.bluetooth_le",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"android.hardware.camera"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.camera.any",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.camera.autofocus",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.camera.capability.manual_post_processing",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"android.hardware.camera.capability.manual_sensor"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.camera.capability.raw",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.camera.flash",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.camera.front",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.camera.level.full",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.faketouch",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.fingerprint",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.location",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.location.gps",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.location.network",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.microphone",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.nfc",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.nfc.any",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.nfc.hce",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.nfc.hcef",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"android.hardware.opengles.aep"}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"android.hardware.ram.normal"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.screen.landscape",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.screen.portrait",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.accelerometer",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.assist",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.barometer",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.compass",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.gyroscope",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.hifi_sensors",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.light",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.proximity",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.stepcounter",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.stepdetector",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.strongbox_keystore",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.telephony",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.telephony.carrierlock",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.telephony.cdma",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.telephony.euicc",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.telephony.gsm",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.touchscreen",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.touchscreen.multitouch",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.touchscreen.multitouch.distinct",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.touchscreen.multitouch.jazzhand",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"android.hardware.usb.accessory"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.usb.host",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.vulkan.compute",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.vulkan.level",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.vulkan.version",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.wifi",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.wifi.aware",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.wifi.direct",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.wifi.passpoint",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.wifi.rtt",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.activities_on_secondary_displays",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.app_widgets",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.autofill",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.backup",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.cant_save_state",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.companion_device_setup",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.connectionservice",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.cts",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.device_admin",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.device_id_attestation",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.file_based_encryption",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.home_screen",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.input_methods",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.live_wallpaper",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.managed_users",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.midi",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"android.software.picture_in_picture"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.print",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"android.software.securely_removes_users"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.sip",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.sip.voip",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.verified_boot",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.voice_recognizers",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.software.webview",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.apps.dialer.SUPPORTED",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"com.google.android.apps.photos.PIXEL_2019_MIDYEAR_PRELOAD"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.feature.EXCHANGE_6_2",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.feature.GOOGLE_BUILD",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"com.google.android.feature.GOOGLE_EXPERIENCE"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.feature.PIXEL_2017_EXPERIENCE",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.feature.PIXEL_2018_EXPERIENCE",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.feature.PIXEL_2019_MIDYEAR_EXPERIENCE",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.feature.PIXEL_EXPERIENCE",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.feature.TURBO_PRELOAD",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.android.feature.WELLBEING",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
protobuf.Tag{Number:1, String:""}:"com.google.android.feature.ZERO_TOUCH"}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.google.hardware.camera.easel_2018",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.verizon.hardware.telephony.ehrpd",
protobuf.Tag{Number:2, String:""}:uint64(0)}, protobuf.Message{protobuf.Tag{Number:1, String:""}:"com.verizon.hardware.telephony.lte",
protobuf.Tag{Number:2, String:""}:uint64(0)}},
protobuf.Tag{Number:1, String:""}:uint64(3),
protobuf.Tag{Number:7, String:""}:uint64(490),
protobuf.Tag{Number:10, String:""}:[]string{"android.hardware.audio.low_latency", "android.hardware.audio.output", "android.hardware.audio.pro", "android.hardware.bluetooth", "android.hardware.bluetooth_le", "android.hardware.camera", "android.hardware.camera.any", "android.hardware.camera.autofocus", "android.hardware.camera.capability.manual_post_processing", "android.hardware.camera.capability.manual_sensor", "android.hardware.camera.capability.raw", "android.hardware.camera.flash", "android.hardware.camera.front", "android.hardware.camera.level.full", "android.hardware.faketouch", "android.hardware.fingerprint", "android.hardware.location", "android.hardware.location.gps", "android.hardware.location.network", "android.hardware.microphone", "android.hardware.nfc", "android.hardware.nfc.any", "android.hardware.nfc.hce", "android.hardware.nfc.hcef", "android.hardware.opengles.aep", "android.hardware.ram.normal", "android.hardware.screen.landscape", "android.hardware.screen.portrait", "android.hardware.sensor.accelerometer", "android.hardware.sensor.assist", "android.hardware.sensor.barometer", "android.hardware.sensor.compass", "android.hardware.sensor.gyroscope", "android.hardware.sensor.hifi_sensors", "android.hardware.sensor.light", "android.hardware.sensor.proximity", "android.hardware.sensor.stepcounter", "android.hardware.sensor.stepdetector", "android.hardware.strongbox_keystore", "android.hardware.telephony", "android.hardware.telephony.carrierlock", "android.hardware.telephony.cdma", "android.hardware.telephony.euicc", "android.hardware.telephony.gsm", "android.hardware.touchscreen", "android.hardware.touchscreen.multitouch", "android.hardware.touchscreen.multitouch.distinct", "android.hardware.touchscreen.multitouch.jazzhand", "android.hardware.usb.accessory", "android.hardware.usb.host", "android.hardware.vulkan.compute", "android.hardware.vulkan.level", "android.hardware.vulkan.version", "android.hardware.wifi", "android.hardware.wifi.aware", "android.hardware.wifi.direct", "android.hardware.wifi.passpoint", "android.hardware.wifi.rtt", "android.software.activities_on_secondary_displays", "android.software.app_widgets", "android.software.autofill", "android.software.backup", "android.software.cant_save_state", "android.software.companion_device_setup", "android.software.connectionservice", "android.software.cts", "android.software.device_admin", "android.software.device_id_attestation", "android.software.file_based_encryption", "android.software.home_screen", "android.software.input_methods", "android.software.live_wallpaper", "android.software.managed_users", "android.software.midi", "android.software.picture_in_picture", "android.software.print", "android.software.securely_removes_users", "android.software.sip", "android.software.sip.voip", "android.software.verified_boot", "android.software.voice_recognizers", "android.software.webview", "com.google.android.apps.dialer.SUPPORTED", "com.google.android.apps.photos.PIXEL_2019_MIDYEAR_PRELOAD", "com.google.android.feature.EXCHANGE_6_2", "com.google.android.feature.GOOGLE_BUILD", "com.google.android.feature.GOOGLE_EXPERIENCE", "com.google.android.feature.PIXEL_2017_EXPERIENCE", "com.google.android.feature.PIXEL_2018_EXPERIENCE", "com.google.android.feature.PIXEL_2019_MIDYEAR_EXPERIENCE", "com.google.android.feature.PIXEL_EXPERIENCE", "com.google.android.feature.TURBO_PRELOAD", "com.google.android.feature.WELLBEING", "com.google.android.feature.ZERO_TOUCH", "com.google.hardware.camera.easel_2018", "com.verizon.hardware.telephony.ehrpd", "com.verizon.hardware.telephony.lte"},
protobuf.Tag{Number:12, String:""}:uint64(1080),
protobuf.Tag{Number:16, String:""}:uint64(0),
protobuf.Tag{Number:21, String:""}:uint64(8),
protobuf.Tag{Number:3, String:""}:uint64(1),
protobuf.Tag{Number:4, String:""}:uint64(2),
protobuf.Tag{Number:9, String:""}:[]string{"android.ext.services", "android.ext.shared", "android.test.base", "android.test.mock", "android.test.runner", "com.android.future.usb.accessory", "com.android.ims.rcsmanager", "com.android.location.provider", "com.android.media.remotedisplay", "com.android.mediadrm.signer", "com.google.android.camera.experimental2018", "com.google.android.dialer.support", "com.google.android.gms", "com.google.android.hardwareinfo", "com.google.android.lowpowermonitordevicefactory", "com.google.android.lowpowermonitordeviceinterface", "com.google.android.maps", "com.google.android.poweranomalydatafactory", "com.google.android.poweranomalydatamodeminterface", "com.qti.snapdragon.sdk.display", "com.qualcomm.embmslibrary", "com.qualcomm.qcrilhook", "com.qualcomm.qti.QtiTelephonyServicelibrary", "com.qualcomm.qti.imscmservice@1.0-java", "com.qualcomm.qti.lpa.uimlpalibrary", "com.qualcomm.qti.ltedirectdiscoverylibrary", "com.qualcomm.qti.remoteSimlock.uimremotesimlocklibrary", "com.qualcomm.qti.uim.uimservicelibrary", "com.quicinc.cne", "com.quicinc.cneapiclient", "com.verizon.embms", "com.verizon.provider", "com.vzw.apnlib", "javax.obex", "org.apache.http.legacy"},
protobuf.Tag{Number:15, String:""}:[]string{"GL_AMD_compressed_ATC_texture", "GL_AMD_performance_monitor", "GL_ANDROID_extension_pack_es31a", "GL_APPLE_texture_2D_limited_npot", "GL_ARB_vertex_buffer_object", "GL_ARM_shader_framebuffer_fetch_depth_stencil", "GL_EXT_EGL_image_array", "GL_EXT_EGL_image_external_wrap_modes", "GL_EXT_EGL_image_storage", "GL_EXT_YUV_target", "GL_EXT_blend_func_extended", "GL_EXT_blit_framebuffer_params", "GL_EXT_buffer_storage", "GL_EXT_clip_control", "GL_EXT_clip_cull_distance", "GL_EXT_color_buffer_float", "GL_EXT_color_buffer_half_float", "GL_EXT_copy_image", "GL_EXT_debug_label", "GL_EXT_debug_marker", "GL_EXT_discard_framebuffer", "GL_EXT_disjoint_timer_query", "GL_EXT_draw_buffers_indexed", "GL_EXT_external_buffer", "GL_EXT_geometry_shader", "GL_EXT_gpu_shader5", "GL_EXT_memory_object", "GL_EXT_memory_object_fd", "GL_EXT_multisampled_render_to_texture", "GL_EXT_multisampled_render_to_texture2", "GL_EXT_primitive_bounding_box", "GL_EXT_protected_textures", "GL_EXT_robustness", "GL_EXT_sRGB", "GL_EXT_sRGB_write_control", "GL_EXT_shader_framebuffer_fetch", "GL_EXT_shader_io_blocks", "GL_EXT_shader_non_constant_global_initializers", "GL_EXT_tessellation_shader", "GL_EXT_texture_border_clamp", "GL_EXT_texture_buffer", "GL_EXT_texture_cube_map_array", "GL_EXT_texture_filter_anisotropic", "GL_EXT_texture_format_BGRA8888", "GL_EXT_texture_format_sRGB_override", "GL_EXT_texture_norm16", "GL_EXT_texture_sRGB_R8", "GL_EXT_texture_sRGB_decode", "GL_EXT_texture_type_2_10_10_10_REV", "GL_KHR_blend_equation_advanced", "GL_KHR_blend_equation_advanced_coherent", "GL_KHR_debug", "GL_KHR_no_error", "GL_KHR_robust_buffer_access_behavior", "GL_KHR_texture_compression_astc_hdr", "GL_KHR_texture_compression_astc_ldr", "GL_NV_shader_noperspective_interpolation", "GL_OES_EGL_image", "GL_OES_EGL_image_external", "GL_OES_EGL_image_external_essl3", "GL_OES_EGL_sync", "GL_OES_blend_equation_separate", "GL_OES_blend_func_separate", "GL_OES_blend_subtract", "GL_OES_compressed_ETC1_RGB8_texture", "GL_OES_compressed_paletted_texture", "GL_OES_depth24", "GL_OES_depth_texture", "GL_OES_depth_texture_cube_map", "GL_OES_draw_texture", "GL_OES_element_index_uint", "GL_OES_framebuffer_object", "GL_OES_get_program_binary", "GL_OES_matrix_palette", "GL_OES_packed_depth_stencil", "GL_OES_point_size_array", "GL_OES_point_sprite", "GL_OES_read_format", "GL_OES_rgb8_rgba8", "GL_OES_sample_shading", "GL_OES_sample_variables", "GL_OES_shader_image_atomic", "GL_OES_shader_multisample_interpolation", "GL_OES_standard_derivatives", "GL_OES_stencil_wrap", "GL_OES_surfaceless_context", "GL_OES_texture_3D", "GL_OES_texture_compression_astc", "GL_OES_texture_cube_map", "GL_OES_texture_env_crossbar", "GL_OES_texture_float", "GL_OES_texture_float_linear", "GL_OES_texture_half_float", "GL_OES_texture_half_float_linear", "GL_OES_texture_mirrored_repeat", "GL_OES_texture_npot", "GL_OES_texture_stencil8", "GL_OES_texture_storage_multisample_2d_array", "GL_OES_vertex_array_object", "GL_OES_vertex_half_float", "GL_OVR_multiview", "GL_OVR_multiview2", "GL_OVR_multiview_multisampled_render_to_texture", "GL_QCOM_alpha_test", "GL_QCOM_extended_get", "GL_QCOM_shader_framebuffer_fetch_noncoherent", "GL_QCOM_texture_foveated", "GL_QCOM_tiled_rendering"},
protobuf.Tag{Number:19, String:""}:uint64(0),
protobuf.Tag{Number:20, String:""}:uint64(8589935000),
protobuf.Tag{Number:2, String:""}:uint64(1),
protobuf.Tag{Number:6, String:""}:uint64(0),
protobuf.Tag{Number:8, String:""}:uint64(196610),
protobuf.Tag{Number:13, String:""}:uint64(2073)},
protobuf.Tag{Number:20, String:""}:uint64(0)}

func checkin() (uint64, error) {
   var req0 = &http.Request{
      Method:"POST",
      URL:&url.URL{Scheme:"https",
         Host:"android.clients.google.com",
         Path:"/checkin", 
      },
      Header:http.Header{
         "Content-Type":[]string{"application/x-protobuffer"},
      },
      Body: io.NopCloser(body1.Encode()),
   }
   res, err := new(http.Transport).RoundTrip(req0)
   if err != nil {
      return 0, err
   }
   defer res.Body.Close()
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(7), nil
}

func details(app string) (uint64, error) {
   id, err := checkin()
   if err != nil {
      return 0, err
   }
   sID := strconv.FormatUint(id, 16)
   fmt.Println(sID)
   var req5 = &http.Request{Method:"GET", URL:&url.URL{Scheme:"https",
      Host:"android.clients.google.com",
      Path:"/fdfe/details", RawQuery:"doc=" + app,
      },
      Header:http.Header{
         "Authorization":[]string{"Bearer " + auth},
         "X-Dfe-Device-Id":[]string{sID},
      },
   }
   res, err := new(http.Transport).RoundTrip(req5)
   if err != nil {
      return 0, err
   }
   mes, err := protobuf.Decode(res.Body)
   if err != nil {
      return 0, err
   }
   return mes.GetUint64(1,2,4,13,1,3), nil
}
