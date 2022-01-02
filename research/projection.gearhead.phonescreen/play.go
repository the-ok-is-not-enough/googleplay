package play

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

const auth = "ya29.a0ARrdaM8WY37NU2AwiJuhmKBCVee_VF0wBYHL1FtWCpZ-_2y-dJgyaEij0GW7Cnvh-NM0GNWpgLFxrFD97wjfSyrkw_f7uiyHnveUKjEJtVQrNobuqXayGZFJWDrXgFIjZdPRf9KXSTWqXtLEDgyMoFsztF1Jw_wma4onRymLnK6fyaBbCqyAFHh8Jo73nVu5ovVOzDQSOR1-dN8A4qdYpcfE7J9drzSUj5QMGmKeh9MhR9qM9EhDgf5ZTn-adEursZYymjRpRxVC8a2EQbbbV6O-IhNav37ptVgK3Ko-yBWxW44UDFY"

var body1 = protobuf.Message{
   protobuf.Tag{Number:2, String:""}:uint64(0),
   protobuf.Tag{Number:4, String:""}:protobuf.Message{
      protobuf.Tag{Number:2, String:""}:uint64(0),
      protobuf.Tag{Number:6, String:""}:"334050",
      protobuf.Tag{Number:7, String:""}:"20815",
      protobuf.Tag{Number:8, String:""}:"mobile-notroaming",
      protobuf.Tag{Number:9, String:""}:uint64(0),
      protobuf.Tag{Number:1, String:""}:protobuf.Message{
         protobuf.Tag{Number:10, String:""}:uint64(29),
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
         protobuf.Tag{Number:8, String:""}:uint64(203615028),
      },
   },
   protobuf.Tag{Number:6, String:""}:"en_GB",
   protobuf.Tag{Number:12, String:""}:"America/Mexico_City",
   protobuf.Tag{Number:14, String:""}:uint64(3),
   protobuf.Tag{Number:18, String:""}:protobuf.Message{
      protobuf.Tag{Number:5, String:""}:uint64(0),
      protobuf.Tag{Number:14, String:""}:[]string{"af", "af_ZA", "am", "am_ET", "ar", "ar_EG", "ar_SA", "ar_XB", "as", "ast", "az", "be", "be_BY", "bg", "bg_BG", "bh_IN", "bn", "bs", "ca", "ca_ES", "cs", "cs_CZ", "cy_GB", "da", "da_DK", "de", "de_DE", "el", "el_GR", "en", "en_AU", "en_CA", "en_GB", "en_IN", "en_US", "en_XA", "es", "es_ES", "es_US", "et", "et_EE", "eu", "fa", "fa_IR", "fi", "fi_FI", "fil", "fil_PH", "fr", "fr_CA", "fr_FR", "gl", "gl_ES", "gu", "hi", "hi_IN", "hr", "hr_HR", "hu", "hu_HU", "hy", "in", "in_ID", "is", "it", "it_IT", "iw", "iw_IL", "ja", "ja_JP", "ka", "kab_DZ", "kk", "km", "kn", "ko", "ko_KR", "ky", "lo", "lt", "lt_LT", "lv", "lv_LV", "mk", "ml", "mn", "mr", "ms", "ms_MY", "my", "nb", "nb_NO", "ne", "nl", "nl_NL", "or", "pa", "pa_IN", "pl", "pl_PL", "pt", "pt_BR", "pt_PT", "ro", "ro_RO", "ru", "ru_RU", "sc_IT", "si", "sk", "sk_SK", "sl", "sl_SI", "sq", "sr", "sr_Latn", "sr_RS", "sv", "sv_SE", "sw", "sw_TZ", "ta", "te", "th", "th_TH", "tr", "tr_TR", "uk", "uk_UA", "ur", "uz", "vi", "vi_VN", "zh_CN", "zh_HK", "zh_TW", "zu", "zu_ZA"},
      protobuf.Tag{Number:26, String:""}:[]protobuf.Message{
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.faketouch", protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.screen.portrait", protobuf.Tag{Number:2, String:""}:uint64(0)},
      },
      protobuf.Tag{Number:1, String:""}:uint64(3),
      protobuf.Tag{Number:7, String:""}:uint64(490),
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
      protobuf.Tag{Number:13, String:""}:uint64(2073),
   },
   protobuf.Tag{Number:20, String:""}:uint64(0),
}

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
