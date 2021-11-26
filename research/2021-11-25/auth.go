package googleplay

type responseWrapper struct { // DONE
   Payload struct { // DONE
      DetailsResponse struct { // DONE
         DocV2 struct { // DONE
            Offer struct { // DONE
               FormattedAmount string `json:"3"` // DONE
            } `json:"8"` // DONE
            Details struct { // DONE
               AppDetails struct { // DONE
                  DeveloperName string `json:"1"` // DONE
                  VersionCode int32 `json:"3"` // DONE
                  Version string `json:"4"` // DONE
                  InstallationSize int64 `json:"9"` // DONE
                  UploadDate string `json:"16"` // DONE
               } `json:"1"` // DONE
            } `json:"13"` // DONE
            AggregateRating struct { // DONE
               OneStar uint64 `json:"4"` // DONE
               TwoStar uint64 `json:"5"` // DONE
               ThreeStar uint64 `json:"6"` // DONE
               FourStar uint64 `json:"7"` // DONE
               FiveStar uint64 `json:"8"` // DONE
            } `json:"14"` // DONE
         } `json:"4"` // DONE
      } `json:"2"` // DONE
      DeliveryResponse struct { // DONE
         Status int32 `json:"1"` // DONE
         AppDeliveryData struct { // DONE
            DownloadURL string `json:"3"` // DONE
            SplitDeliveryData []struct {
               Name string `json:"1"`
               DownloadURL string `json:"5"`
            } `json:"15"`
         } `json:"2"` // DONE
      } `json:"21"` // DONE
   } `json:"1"` // DONE
} // DONE
