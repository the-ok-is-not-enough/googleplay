# Vimeo

Size 35.5 MB.

Requires Android 26.

https://apkpure.com/vimeo/com.vimeo.android.videoapp

so, where are the exceptions?

~~~
smali_classes2\g0\k1\h\n.smali
1743:    invoke-direct {p2, p1}, Ljavax/net/ssl/SSLPeerUnverifiedException;-><init>(Ljava/lang/String;)V
1778:    invoke-direct {p1, p2}, Ljavax/net/ssl/SSLPeerUnverifiedException;-><init>(Ljava/lang/String;)V

smali_classes2\g0\k1\m\r\b.smali
111:    invoke-direct {p2, v0}, Ljavax/net/ssl/SSLPeerUnverifiedException;-><init>(Ljava/lang/String;)V

smali_classes2\g0\k1\p\a.smali
221:    invoke-direct {p1, p2}, Ljavax/net/ssl/SSLPeerUnverifiedException;-><init>(Ljava/lang/String;)V
243:    invoke-direct {p2, p1}, Ljavax/net/ssl/SSLPeerUnverifiedException;-><init>(Ljava/lang/String;)V

smali_classes2\g0\r.smali
796:    invoke-direct {v1, v0}, Ljavax/net/ssl/SSLPeerUnverifiedException;-><init>(Ljava/lang/String;)V
~~~

what are the return types?

~~~
smali_classes2\g0\k1\h\n.smali
V

smali_classes2\g0\k1\m\r\b.smali
Ljava/util/List;

smali_classes2\g0\k1\p\a.smali
Ljava/util/List;

smali_classes2\g0\r.smali
V
~~~
