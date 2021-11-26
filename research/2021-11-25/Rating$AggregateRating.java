public final class Rating$AggregateRating extends com.google.protobuf.nano.MessageNano {
    public double bayesianMeanRating;
    public long commentCount;
    public long fiveStarRatings;
    public long fourStarRatings;
    public boolean hasBayesianMeanRating;
    public boolean hasCommentCount;
    public boolean hasFiveStarRatings;
    public boolean hasFourStarRatings;
    public boolean hasOneStarRatings;
    public boolean hasRatingsCount;
    public boolean hasStarRating;
    public boolean hasThreeStarRatings;
    public boolean hasThumbsDownCount;
    public boolean hasThumbsUpCount;
    public boolean hasTwoStarRatings;
    public boolean hasType;
    public long oneStarRatings;
    public long ratingsCount;
    public float starRating;
    public long threeStarRatings;
    public long thumbsDownCount;
    public long thumbsUpCount;
    public com.google.android.finsky.protos.Rating$Tip[] tip;
    public long twoStarRatings;
    public int type;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p9)
    {
        if ((this.type != 1) || (this.hasType)) {
            p9.writeInt32(1, this.type);
        }
        if ((this.hasStarRating) || (Float.floatToIntBits(this.starRating) != Float.floatToIntBits(0))) {
            p9.writeFloat(2, this.starRating);
        }
        if ((this.hasRatingsCount) || (this.ratingsCount != 0)) {
            p9.writeUInt64(3, this.ratingsCount);
        }
        if ((this.hasOneStarRatings) || (this.oneStarRatings != 0)) {
            p9.writeUInt64(4, this.oneStarRatings);
        }
        if ((this.hasTwoStarRatings) || (this.twoStarRatings != 0)) {
            p9.writeUInt64(5, this.twoStarRatings);
        }
        if ((this.hasThreeStarRatings) || (this.threeStarRatings != 0)) {
            p9.writeUInt64(6, this.threeStarRatings);
        }
        if ((this.hasFourStarRatings) || (this.fourStarRatings != 0)) {
            p9.writeUInt64(7, this.fourStarRatings);
        }
        if ((this.hasFiveStarRatings) || (this.fiveStarRatings != 0)) {
            p9.writeUInt64(8, this.fiveStarRatings);
        }
        if ((this.hasThumbsUpCount) || (this.thumbsUpCount != 0)) {
            p9.writeUInt64(9, this.thumbsUpCount);
        }
        if ((this.hasThumbsDownCount) || (this.thumbsDownCount != 0)) {
            p9.writeUInt64(10, this.thumbsDownCount);
        }
        if ((this.hasCommentCount) || (this.commentCount != 0)) {
            p9.writeUInt64(11, this.commentCount);
        }
        if ((this.hasBayesianMeanRating) || (Double.doubleToLongBits(this.bayesianMeanRating) != Double.doubleToLongBits(0))) {
            p9.writeDouble(12, this.bayesianMeanRating);
        }
        if ((this.tip != null) && (this.tip.length > 0)) {
            int v1 = 0;
            while (v1 < this.tip.length) {
                com.google.android.finsky.protos.Rating$Tip v0 = this.tip[v1];
                if (v0 != null) {
                    p9.writeMessage(13, v0);
                }
                v1++;
            }
        }
        super.writeTo(p9);
        return;
    }
}
