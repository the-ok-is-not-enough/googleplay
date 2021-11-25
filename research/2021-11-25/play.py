from androguard.misc import AnalyzeAPK
a,d,dx= AnalyzeAPK('com.android.vending_6.2.10.apk')
f = open('file.java', 'w')

for dd in d:
   for clas in dd.get_classes():
      name = clas.get_name()
      #if 'proto' in name:
      print(name, file=f)
      src = dd.get_class(name).get_source()
      print(src, file=f)
