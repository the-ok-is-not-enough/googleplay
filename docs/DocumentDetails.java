public final class DocDetails$DocumentDetails extends com.google.protobuf.nano.MessageNano {
    public com.google.android.finsky.protos.AlbumDetails albumDetails;
    public com.google.android.finsky.protos.AppDetails appDetails;
    public com.google.android.finsky.protos.ArtistDetails artistDetails;
    public com.google.android.finsky.protos.BookDetails bookDetails;
    public com.google.android.finsky.protos.BookSeriesDetails bookSeriesDetails;
    public com.google.android.finsky.protos.DeveloperDetails developerDetails;
    public com.google.android.finsky.protos.MagazineDetails magazineDetails;
    public com.google.android.finsky.protos.DocDetails$PersonDetails personDetails;
    public com.google.android.finsky.protos.SongDetails songDetails;
    public com.google.android.finsky.protos.SubscriptionDetails subscriptionDetails;
    public com.google.android.finsky.protos.DocDetails$TalentDetails talentDetails;
    public com.google.android.finsky.protos.TvEpisodeDetails tvEpisodeDetails;
    public com.google.android.finsky.protos.TvSeasonDetails tvSeasonDetails;
    public com.google.android.finsky.protos.TvShowDetails tvShowDetails;
    public com.google.android.finsky.protos.VideoDetails videoDetails;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p3)
    {
        if (this.appDetails != null) {
            p3.writeMessage(1, this.appDetails);
        }
        if (this.albumDetails != null) {
            p3.writeMessage(2, this.albumDetails);
        }
        if (this.artistDetails != null) {
            p3.writeMessage(3, this.artistDetails);
        }
        if (this.songDetails != null) {
            p3.writeMessage(4, this.songDetails);
        }
        if (this.bookDetails != null) {
            p3.writeMessage(5, this.bookDetails);
        }
        if (this.videoDetails != null) {
            p3.writeMessage(6, this.videoDetails);
        }
        if (this.subscriptionDetails != null) {
            p3.writeMessage(7, this.subscriptionDetails);
        }
        if (this.magazineDetails != null) {
            p3.writeMessage(8, this.magazineDetails);
        }
        if (this.tvShowDetails != null) {
            p3.writeMessage(9, this.tvShowDetails);
        }
        if (this.tvSeasonDetails != null) {
            p3.writeMessage(10, this.tvSeasonDetails);
        }
        if (this.tvEpisodeDetails != null) {
            p3.writeMessage(11, this.tvEpisodeDetails);
        }
        if (this.personDetails != null) {
            p3.writeMessage(12, this.personDetails);
        }
        if (this.talentDetails != null) {
            p3.writeMessage(13, this.talentDetails);
        }
        if (this.developerDetails != null) {
            p3.writeMessage(14, this.developerDetails);
        }
        if (this.bookSeriesDetails != null) {
            p3.writeMessage(15, this.bookSeriesDetails);
        }
        super.writeTo(p3);
        return;
    }
}
