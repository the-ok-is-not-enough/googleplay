# AndroGuard

~~~py
from androguard.misc import AnalyzeAPK

a,d,dx= AnalyzeAPK('file.apk')
f = open('file.java', 'w')

for dd in d:
   for clas in dd.get_classes():
      name = clas.get_name()
      if 'something' in name:
         print(name, file=f)
         src = dd.get_class(name).get_source()
         print(src, file=f)
~~~
