// Lcom/google/android/finsky/protos/DocV2;
public final class DocV2 extends com.google.protobuf.nano.MessageNano {
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p6)
    {
        if ((this.hasDocid) || (!this.docid.equals(""))) {
            p6.writeString(1, this.docid);
        }
        if ((this.hasBackendDocid) || (!this.backendDocid.equals(""))) {
            p6.writeString(2, this.backendDocid);
        }
        if ((this.docType != 1) || (this.hasDocType)) {
            p6.writeInt32(3, this.docType);
        }
        if ((this.backendId != 0) || (this.hasBackendId)) {
            p6.writeInt32(4, this.backendId);
        }
        if ((this.hasTitle) || (!this.title.equals(""))) {
            p6.writeString(5, this.title);
        }
        if ((this.hasCreator) || (!this.creator.equals(""))) {
            p6.writeString(6, this.creator);
        }
        if ((this.hasDescriptionHtml) || (!this.descriptionHtml.equals(""))) {
            p6.writeString(7, this.descriptionHtml);
        }
        if ((this.offer != null) && (this.offer.length > 0)) {
            int v1_3 = 0;
            while (v1_3 < this.offer.length) {
                com.google.android.finsky.protos.ReviewTip v0_3 = this.offer[v1_3];
                if (v0_3 != null) {
                    p6.writeMessage(8, v0_3);
                }
                v1_3++;
            }
        }
        if (this.availability != null) {
            p6.writeMessage(9, this.availability);
        }
        if ((this.image != null) && (this.image.length > 0)) {
            int v1_2 = 0;
            while (v1_2 < this.image.length) {
                com.google.android.finsky.protos.ReviewTip v0_2 = this.image[v1_2];
                if (v0_2 != null) {
                    p6.writeMessage(10, v0_2);
                }
                v1_2++;
            }
        }
        if ((this.child != null) && (this.child.length > 0)) {
            int v1_1 = 0;
            while (v1_1 < this.child.length) {
                com.google.android.finsky.protos.ReviewTip v0_1 = this.child[v1_1];
                if (v0_1 != null) {
                    p6.writeMessage(11, v0_1);
                }
                v1_1++;
            }
        }
        if (this.containerMetadata != null) {
            p6.writeMessage(12, this.containerMetadata);
        }
        if (this.details != null) {
            p6.writeMessage(13, this.details);
        }
        if (this.aggregateRating != null) {
            p6.writeMessage(14, this.aggregateRating);
        }
        if (this.annotations != null) {
            p6.writeMessage(15, this.annotations);
        }
        if ((this.hasDetailsUrl) || (!this.detailsUrl.equals(""))) {
            p6.writeString(16, this.detailsUrl);
        }
        if ((this.hasShareUrl) || (!this.shareUrl.equals(""))) {
            p6.writeString(17, this.shareUrl);
        }
        if ((this.hasReviewsUrl) || (!this.reviewsUrl.equals(""))) {
            p6.writeString(18, this.reviewsUrl);
        }
        if ((this.hasBackendUrl) || (!this.backendUrl.equals(""))) {
            p6.writeString(19, this.backendUrl);
        }
        if ((this.hasPurchaseDetailsUrl) || (!this.purchaseDetailsUrl.equals(""))) {
            p6.writeString(20, this.purchaseDetailsUrl);
        }
        if ((this.hasDetailsReusable) || (this.detailsReusable)) {
            p6.writeBool(21, this.detailsReusable);
        }
        if ((this.hasSubtitle) || (!this.subtitle.equals(""))) {
            p6.writeString(22, this.subtitle);
        }
        if ((this.hasTranslatedDescriptionHtml) || (!this.translatedDescriptionHtml.equals(""))) {
            p6.writeString(23, this.translatedDescriptionHtml);
        }
        if ((this.hasServerLogsCookie) || (!java.util.Arrays.equals(this.serverLogsCookie, com.google.protobuf.nano.WireFormatNano.EMPTY_BYTES))) {
            p6.writeBytes(24, this.serverLogsCookie);
        }
        if (this.productDetails != null) {
            p6.writeMessage(25, this.productDetails);
        }
        if ((this.hasMature) || (this.mature)) {
            p6.writeBool(26, this.mature);
        }
        if ((this.hasPromotionalDescription) || (!this.promotionalDescription.equals(""))) {
            p6.writeString(27, this.promotionalDescription);
        }
        if ((this.hasAvailableForPreregistration) || (this.availableForPreregistration)) {
            p6.writeBool(29, this.availableForPreregistration);
        }
        if ((this.tip != null) && (this.tip.length > 0)) {
            int v1_0 = 0;
            while (v1_0 < this.tip.length) {
                com.google.android.finsky.protos.ReviewTip v0_0 = this.tip[v1_0];
                if (v0_0 != null) {
                    p6.writeMessage(30, v0_0);
                }
                v1_0++;
            }
        }
        if ((this.hasSnippetsUrl) || (!this.snippetsUrl.equals(""))) {
            p6.writeString(31, this.snippetsUrl);
        }
        if ((this.hasForceShareability) || (this.forceShareability)) {
            p6.writeBool(32, this.forceShareability);
        }
        if ((this.hasUseWishlistAsPrimaryAction) || (this.useWishlistAsPrimaryAction)) {
            p6.writeBool(33, this.useWishlistAsPrimaryAction);
        }
        super.writeTo(p6);
        return;
    }
}
