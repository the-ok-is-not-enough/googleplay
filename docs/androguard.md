# AndroGuard

~~~py
from androguard.misc import AnalyzeAPK

a,d,dx= AnalyzeAPK('Bandcamp for Artists and Labels_v1.0.16.apk')
f = open('BAL 1.0.16.java', 'w')

for dd in d:
   for clas in dd.get_classes():
      name = clas.get_name()
      if 'bandcamp' in name:
         print(name, file=f)
         src = dd.get_class(name).get_source()
         print(src, file=f)
~~~
