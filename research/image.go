package main
import gp "github.com/89z/googleplay"

func main() {
   token, err := gp.OpenToken("token.json")
   if err != nil {
      panic(err)
   }
   dev, err := gp.OpenDevice("device.json")
   if err != nil {
      panic(err)
   }
   head, err := token.Header(dev)
   if err != nil {
      panic(err)
   }
   detail, err := head.Details("com.google.android.youtube")
   if err != nil {
      panic(err)
   }
   println(detail.Icon())
}
