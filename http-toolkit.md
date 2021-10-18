## HTTP Toolkit

First download the program [1]. Start the program, and click **Android device
connected via ADB**. Then browse to:

~~~
https://example.com
~~~

To quit, click **File**, **Quit**. However a bug currently exists [2], where a
zombie `node.exe` keep running, so you will need to kill that before the next
start.

1. https://github.com/httptoolkit/httptoolkit-desktop/releases
2. https://github.com/httptoolkit/httptoolkit/issues/171

## Certificate pinning

https://httptoolkit.tech/blog/frida-certificate-pinning

## Is HTTP Toolkit open source?

No.

## Issues

- https://github.com/httptoolkit/httptoolkit/issues/171
- https://github.com/httptoolkit/httptoolkit/issues/192
