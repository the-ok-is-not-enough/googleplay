# TikTok

https://play.google.com/store/apps/details?id=com.zhiliaoapp.musically

## Download

TikTok only supports ARM, so to get the APK, you need to request using an ARM
device:

~~~
armeabi-v7a
arm64-v8a
~~~

## Genymotion

Create a new device, then start it, then close it, then delete, and it will give
you this error:

~~~
Failed to delete the virtual device
~~~

Close program and start again, and device should be gone. Future deletes should
also work as expected. However even after that, I couldnt really get the
emulator to work. One start would get stuck on black screen, one start would get
stuck on Android screen for multiple minutes.

## Android Studio, ARM

If you create an ARM Virtual Device using Android 11, supposedly the speed
should be good [1]. However I get this message:

~~~
[Window Title]
qemu-system-aarch64.exe

[Main Instruction]
qemu-system-aarch64.exe is not responding
~~~

I tried waiting for two minutes, but the emulator never started. I was able to
start the device by changing Graphics setting to `Hardware - GLES 2.0`. However
even with that, its still take about five minutes to start up [2], and over two
minute to install an app.

1. https://android-developers.googleblog.com/2020/03/run-arm-apps-on-android-emulator.html
2. https://android.stackexchange.com/questions/134545/speed-up-android-sdk

## Android Studio, x86

If your Virtual Device is `x86` or `x86_64`, you might get this error:

~~~
INSTALL_FAILED_NO_MATCHING_ABIS
~~~

If the APK is `armeabi-v7a`, then you can install with `x86` Virtual Device, if
using Android 9 (API 28) or higher. If the APK is `arm64-v8a`, then you can
install with `x86_64` Virtual Device, if using Android 11 (API 30) or higher.

## Run, x86

If you try to start the app, you will get this message:

~~~
TikTok keeps stopping
~~~

Running with Frida shows this:

~~~
Process crashed: Trace/BPT trap
Build fingerprint: 'google/sdk_gphone_x86_arm/generic_x86_arm:9/PSR1.180720.122/6736742:userdebug/dev-keys'
Revision: '0'
ABI: 'x86'
pid: 5266, tid: 5321, name: tp-io-13  >>> com.zhiliaoapp.musically <<<
signal 6 (SIGABRT), code -6 (SI_TKILL), fault addr --------
Abort message: 'vendor/unbundled_google/libs/ndk_translation/runtime/host_call_frame.cc:65: CHECK failed: 3331640760 == 3331641312'
    eax 00000000  ebx 00001492  ecx 000014c9  edx 00000006
    edi 00001492  esi c6c59f4c
    ebp c6c59f18  esp c6c59ea8  eip ef330b39
backtrace:
    #00 pc 00000b39  [vdso:ef330000] (__kernel_vsyscall+9)
    #01 pc 0001fdf8  /system/lib/libc.so (offset 0x1f000) (syscall+40)
    #02 pc 00022ed3  /system/lib/libc.so (offset 0x22000) (abort+115)
    #03 pc 0000040e  <anonymous:eb4c7000>
[Android Emulator 5554::com.zhiliaoapp.musically]->
~~~

## Run, x86\_64

App crashes on start. If you try it with Frida, you get this:

~~~
Spawned `com.zhiliaoapp.musically`. Resuming main thread!
[Android Emulator 5554::com.zhiliaoapp.musically]-> Process crashed: SIGSEGV SI_KERNEL
Build fingerprint: 'google/sdk_gphone_x86_64_arm64/generic_x86_64_arm64:11/RSR1.210722.013/7800151:userdebug/dev-keys'
Revision: '0'
ABI: 'x86_64'
Timestamp: 2021-11-22 19:54:06-0600
pid: 5981, tid: 5981, name: re-initialized>  >>> <pre-initialized> <<<
uid: 10168
signal 11 (SIGSEGV), code 128 (SI_KERNEL), fault addr 0x0
    rax f4fd14c690eb21b7  rbx 00007c01a6c3ec70  rcx 00007ffdca8fa7f4  rdx 0000000000000000
    r8  00007ffdca8faa48  r9  00007c00e5a0dc43  r10 00007c00e5a0dc46  r11 00007ffdca8faa70
    r12 0000000012e05df0  r13 00000000702eb608  r14 00007ffdca8fa7f0  r15 000000007085a498
    rdi 00007c01a6c3ec70  rsi 00007ffdca8fa7f0
    rbp 000000006f5f62f8  rsp 00007ffdca8fa6f8  rip 00007c00e68bc373
backtrace:
      #00 pc 000000000047f373  /apex/com.android.art/lib64/libart.so!libart.so (offset 0x436000) (art::JNI<false>::GetStringUTFChars(_JNIEnv*, _jstring*, unsigned char*)+51) (BuildId: 7fbaf2a1a3317bd634b00eb90e32291e)
      #01 pc 00000000000ac630  /system/lib64/libandroid_runtime.so (android::(anonymous namespace)::SystemProperties_getSS(_JNIEnv*, _jclass*, _jstring*, _jstring*)+80) (BuildId: 84eb9c8bad06a5ac4720d16d40f66380)
      #02 pc 0000000000237cdb  /system/framework/x86_64/boot-framework.oat (art_jni_trampoline+251) (BuildId: 9ae0dca73129fa9275af95cba6a7cdde25868e76)
      #03 pc 0000000000693d26  /system/framework/x86_64/boot-framework.oat!boot-framework.oat (offset 0x42d000) (android.os.SystemProperties.get+54) (BuildId: 9ae0dca73129fa9275af95cba6a7cdde25868e76)
      #04 pc 000000000041d97b  /system/framework/x86_64/boot-framework.oat!boot-framework.oat (offset 0x41d000) (android.app.ActivityThread.handleBindApplication+91) (BuildId: 9ae0dca73129fa9275af95cba6a7cdde25868e76)
~~~

https://issuetracker.google.com/issues/207399356
