# November 2 2021

This issue is helpful:

https://github.com/EFForg/apkeep/issues/17

Fix:

https://github.com/NoMore201/googleplay-api/pull/153

The answer is here:

https://github.com/NoMore201/googleplay-api/blob/664c399f/gpapi/googleplay.py#L242-L244

~~~
pip install cryptography
pip install protobuf
pip install requests
pip install urllib3==1.25.11
~~~

After getting a new device ID, you have to wait about 9 seconds before you try
to use it.
