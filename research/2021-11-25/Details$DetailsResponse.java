public final class Details$DetailsResponse extends com.google.protobuf.nano.MessageNano {
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p6)
    {
        if (this.docV1 != null) {
            p6.writeMessage(1, this.docV1);
        }
        if (this.userReview != null) {
            p6.writeMessage(3, this.userReview);
        }
        if (this.docV2 != null) {
            p6.writeMessage(4, this.docV2);
        }
        if ((this.hasFooterHtml) || (!this.footerHtml.equals(""))) {
            p6.writeString(5, this.footerHtml);
        }
        if ((this.hasServerLogsCookie) || (!java.util.Arrays.equals(this.serverLogsCookie, com.google.protobuf.nano.WireFormatNano.EMPTY_BYTES))) {
            p6.writeBytes(6, this.serverLogsCookie);
        }
        if ((this.discoveryBadge != null) && (this.discoveryBadge.length > 0)) {
            int v1 = 0;
            while (v1 < this.discoveryBadge.length) {
                com.google.android.finsky.protos.Details$DiscoveryBadge v0 = this.discoveryBadge[v1];
                if (v0 != null) {
                    p6.writeMessage(7, v0);
                }
                v1++;
            }
        }
        if ((this.hasEnableReviews) || (this.enableReviews != 1)) {
            p6.writeBool(8, this.enableReviews);
        }
        super.writeTo(p6);
        return;
    }
}
