public final class Common$Offer extends com.google.protobuf.nano.MessageNano {
    private static volatile com.google.android.finsky.protos.Common$Offer[] _emptyArray;
    public String buyButtonLabel;
    public boolean checkoutFlowRequired;
    public com.google.android.finsky.protos.Common$Offer[] convertedPrice;
    public String currencyCode;
    public String formattedAmount;
    public String formattedDescription;
    public String formattedFullAmount;
    public String formattedName;
    public long fullPriceMicros;
    public boolean hasBuyButtonLabel;
    public boolean hasCheckoutFlowRequired;
    public boolean hasCurrencyCode;
    public boolean hasFormattedAmount;
    public boolean hasFormattedDescription;
    public boolean hasFormattedFullAmount;
    public boolean hasFormattedName;
    public boolean hasFullPriceMicros;
    public boolean hasInstantPurchaseEnabled;
    public boolean hasLicensedOfferType;
    public boolean hasMicros;
    public boolean hasOfferId;
    public boolean hasOfferType;
    public boolean hasOnSaleDate;
    public boolean hasOnSaleDateDisplayTimeZoneOffsetMsec;
    public boolean hasPreorder;
    public boolean hasPreorderFulfillmentDisplayDate;
    public boolean hasRepeatLastPayment;
    public boolean hasTemporarilyFree;
    public boolean instantPurchaseEnabled;
    public com.google.android.finsky.protos.Common$LicenseTerms licenseTerms;
    public int licensedOfferType;
    public long micros;
    public String offerId;
    public com.google.android.finsky.protos.Common$OfferPayment[] offerPayment;
    public int offerType;
    public long onSaleDate;
    public int onSaleDateDisplayTimeZoneOffsetMsec;
    public boolean preorder;
    public long preorderFulfillmentDisplayDate;
    public String[] promotionLabel;
    public com.google.android.finsky.protos.Common$RentalTerms rentalTerms;
    public boolean repeatLastPayment;
    public com.google.android.finsky.protos.Common$SubscriptionContentTerms subscriptionContentTerms;
    public com.google.android.finsky.protos.Common$SubscriptionTerms subscriptionTerms;
    public boolean temporarilyFree;
    public com.google.android.finsky.protos.Common$VoucherOfferTerms voucherTerms;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p10)
    {
        if ((this.hasMicros) || (this.micros != 0)) {
            p10.writeInt64(1, this.micros);
        }
        if ((this.hasCurrencyCode) || (!this.currencyCode.equals(""))) {
            p10.writeString(2, this.currencyCode);
        }
        if ((this.hasFormattedAmount) || (!this.formattedAmount.equals(""))) {
            p10.writeString(3, this.formattedAmount);
        }
        if ((this.convertedPrice != null) && (this.convertedPrice.length > 0)) {
            int v1_2 = 0;
            while (v1_2 < this.convertedPrice.length) {
                com.google.android.finsky.protos.Common$OfferPayment v0_2 = this.convertedPrice[v1_2];
                if (v0_2 != null) {
                    p10.writeMessage(4, v0_2);
                }
                v1_2++;
            }
        }
        if ((this.hasCheckoutFlowRequired) || (this.checkoutFlowRequired)) {
            p10.writeBool(5, this.checkoutFlowRequired);
        }
        if ((this.hasFullPriceMicros) || (this.fullPriceMicros != 0)) {
            p10.writeInt64(6, this.fullPriceMicros);
        }
        if ((this.hasFormattedFullAmount) || (!this.formattedFullAmount.equals(""))) {
            p10.writeString(7, this.formattedFullAmount);
        }
        if ((this.offerType != 1) || (this.hasOfferType)) {
            p10.writeInt32(8, this.offerType);
        }
        if (this.rentalTerms != null) {
            p10.writeMessage(9, this.rentalTerms);
        }
        if ((this.hasOnSaleDate) || (this.onSaleDate != 0)) {
            p10.writeInt64(10, this.onSaleDate);
        }
        if ((this.promotionLabel != null) && (this.promotionLabel.length > 0)) {
            int v1_1 = 0;
            while (v1_1 < this.promotionLabel.length) {
                com.google.android.finsky.protos.Common$OfferPayment v0_1 = this.promotionLabel[v1_1];
                if (v0_1 != null) {
                    p10.writeString(11, v0_1);
                }
                v1_1++;
            }
        }
        if (this.subscriptionTerms != null) {
            p10.writeMessage(12, this.subscriptionTerms);
        }
        if ((this.hasFormattedName) || (!this.formattedName.equals(""))) {
            p10.writeString(13, this.formattedName);
        }
        if ((this.hasFormattedDescription) || (!this.formattedDescription.equals(""))) {
            p10.writeString(14, this.formattedDescription);
        }
        if ((this.hasPreorder) || (this.preorder)) {
            p10.writeBool(15, this.preorder);
        }
        if ((this.hasOnSaleDateDisplayTimeZoneOffsetMsec) || (this.onSaleDateDisplayTimeZoneOffsetMsec != 0)) {
            p10.writeInt32(16, this.onSaleDateDisplayTimeZoneOffsetMsec);
        }
        if ((this.licensedOfferType != 1) || (this.hasLicensedOfferType)) {
            p10.writeInt32(17, this.licensedOfferType);
        }
        if (this.subscriptionContentTerms != null) {
            p10.writeMessage(18, this.subscriptionContentTerms);
        }
        if ((this.hasOfferId) || (!this.offerId.equals(""))) {
            p10.writeString(19, this.offerId);
        }
        if ((this.hasPreorderFulfillmentDisplayDate) || (this.preorderFulfillmentDisplayDate != 0)) {
            p10.writeInt64(20, this.preorderFulfillmentDisplayDate);
        }
        if (this.licenseTerms != null) {
            p10.writeMessage(21, this.licenseTerms);
        }
        if ((this.hasTemporarilyFree) || (this.temporarilyFree)) {
            p10.writeBool(22, this.temporarilyFree);
        }
        if (this.voucherTerms != null) {
            p10.writeMessage(23, this.voucherTerms);
        }
        if ((this.offerPayment != null) && (this.offerPayment.length > 0)) {
            int v1_0 = 0;
            while (v1_0 < this.offerPayment.length) {
                com.google.android.finsky.protos.Common$OfferPayment v0_0 = this.offerPayment[v1_0];
                if (v0_0 != null) {
                    p10.writeMessage(24, v0_0);
                }
                v1_0++;
            }
        }
        if ((this.hasRepeatLastPayment) || (this.repeatLastPayment)) {
            p10.writeBool(25, this.repeatLastPayment);
        }
        if ((this.hasBuyButtonLabel) || (!this.buyButtonLabel.equals(""))) {
            p10.writeString(26, this.buyButtonLabel);
        }
        if ((this.hasInstantPurchaseEnabled) || (this.instantPurchaseEnabled)) {
            p10.writeBool(27, this.instantPurchaseEnabled);
        }
        super.writeTo(p10);
        return;
    }
}
