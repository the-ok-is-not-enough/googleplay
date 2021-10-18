$env:PATH = 'C:\Users\Steven\AppData\Local\Android\Sdk\build-tools\31.0.0'
$env:PATH += ';C:\Users\Steven\AppData\Local\Android\Sdk\platform-tools'
$env:PATH += ';C:\python\Scripts'

objection patchapk -s Vimeo_v3.50.1_apkpure.com.apk
