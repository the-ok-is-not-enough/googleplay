package googleplay

type responseWrapper struct { // DONE
   Payload struct { // DONE
      DetailsResponse struct { // DONE
         DocV2 struct { // DONE
            Offer struct { // DONE
               FormattedAmount string `json:"3"` // DONE
            } `json:"8"` // DONE
            Details struct { // DONE
               AppDetails struct {
                  DeveloperName string `json:"1"`
                  VersionCode int32 `json:"3"`
                  Version string `json:"4"`
                  InstallationSize int64 `json:"9"`
                  UploadDate string `json:"16"`
               } `json:"1"`
            } `json:"13"` // DONE
            AggregateRating struct {
               OneStar uint64 `json:"4"`
               TwoStar uint64 `json:"5"`
               ThreeStar uint64 `json:"6"`
               FourStar uint64 `json:"7"`
               FiveStar uint64 `json:"8"`
            } `json:"14"`
         } `json:"4"` // DONE
      } `json:"2"` // DONE
      DeliveryResponse struct {
         Status int32 `json:"1"`
         AppDeliveryData struct {
            DownloadURL string `json:"3"`
            SplitDeliveryData []struct {
               Name string `json:"1"`
               DownloadURL string `json:"5"`
            } `json:"15"`
         } `json:"2"`
      } `json:"21"`
   } `json:"1"` // DONE
} // DONE
