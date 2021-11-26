public final class AndroidAppDeliveryData extends com.google.protobuf.nano.MessageNano {
    public com.google.android.finsky.protos.AppFileMetadata[] additionalFile;
    public com.google.android.finsky.protos.HttpCookie[] downloadAuthCookie;
    public long downloadSize;
    public String downloadUrl;
    public com.google.android.finsky.protos.EncryptionParams encryptionParams;
    public boolean everExternallyHosted;
    public boolean forwardLocked;
    public long gzippedDownloadSize;
    public String gzippedDownloadUrl;
    public boolean hasDownloadSize;
    public boolean hasDownloadUrl;
    public boolean hasEverExternallyHosted;
    public boolean hasForwardLocked;
    public boolean hasGzippedDownloadSize;
    public boolean hasGzippedDownloadUrl;
    public boolean hasImmediateStartNeeded;
    public boolean hasInstallLocation;
    public boolean hasPostInstallRefundWindowMillis;
    public boolean hasRefundTimeout;
    public boolean hasServerInitiated;
    public boolean hasSignature;
    public boolean immediateStartNeeded;
    public int installLocation;
    public com.google.android.finsky.protos.AndroidAppPatchData patchData;
    public long postInstallRefundWindowMillis;
    public long refundTimeout;
    public boolean serverInitiated;
    public String signature;
    public com.google.android.finsky.protos.SplitDeliveryData[] splitDeliveryData;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p10)
    {
        if ((this.hasDownloadSize) || (this.downloadSize != 0)) {
            p10.writeInt64(1, this.downloadSize);
        }
        if ((this.hasSignature) || (!this.signature.equals(""))) {
            p10.writeString(2, this.signature);
        }
        if ((this.hasDownloadUrl) || (!this.downloadUrl.equals(""))) {
            p10.writeString(3, this.downloadUrl);
        }
        if ((this.additionalFile != null) && (this.additionalFile.length > 0)) {
            int v1_2 = 0;
            while (v1_2 < this.additionalFile.length) {
                com.google.android.finsky.protos.SplitDeliveryData v0_2 = this.additionalFile[v1_2];
                if (v0_2 != null) {
                    p10.writeMessage(4, v0_2);
                }
                v1_2++;
            }
        }
        if ((this.downloadAuthCookie != null) && (this.downloadAuthCookie.length > 0)) {
            int v1_1 = 0;
            while (v1_1 < this.downloadAuthCookie.length) {
                com.google.android.finsky.protos.SplitDeliveryData v0_1 = this.downloadAuthCookie[v1_1];
                if (v0_1 != null) {
                    p10.writeMessage(5, v0_1);
                }
                v1_1++;
            }
        }
        if ((this.hasForwardLocked) || (this.forwardLocked)) {
            p10.writeBool(6, this.forwardLocked);
        }
        if ((this.hasRefundTimeout) || (this.refundTimeout != 0)) {
            p10.writeInt64(7, this.refundTimeout);
        }
        if ((this.hasServerInitiated) || (this.serverInitiated != 1)) {
            p10.writeBool(8, this.serverInitiated);
        }
        if ((this.hasPostInstallRefundWindowMillis) || (this.postInstallRefundWindowMillis != 0)) {
            p10.writeInt64(9, this.postInstallRefundWindowMillis);
        }
        if ((this.hasImmediateStartNeeded) || (this.immediateStartNeeded)) {
            p10.writeBool(10, this.immediateStartNeeded);
        }
        if (this.patchData != null) {
            p10.writeMessage(11, this.patchData);
        }
        if (this.encryptionParams != null) {
            p10.writeMessage(12, this.encryptionParams);
        }
        if ((this.hasGzippedDownloadUrl) || (!this.gzippedDownloadUrl.equals(""))) {
            p10.writeString(13, this.gzippedDownloadUrl);
        }
        if ((this.hasGzippedDownloadSize) || (this.gzippedDownloadSize != 0)) {
            p10.writeInt64(14, this.gzippedDownloadSize);
        }
        if ((this.splitDeliveryData != null) && (this.splitDeliveryData.length > 0)) {
            int v1_0 = 0;
            while (v1_0 < this.splitDeliveryData.length) {
                com.google.android.finsky.protos.SplitDeliveryData v0_0 = this.splitDeliveryData[v1_0];
                if (v0_0 != null) {
                    p10.writeMessage(15, v0_0);
                }
                v1_0++;
            }
        }
        if ((this.installLocation != 0) || (this.hasInstallLocation)) {
            p10.writeInt32(16, this.installLocation);
        }
        if ((this.hasEverExternallyHosted) || (this.everExternallyHosted)) {
            p10.writeBool(17, this.everExternallyHosted);
        }
        super.writeTo(p10);
        return;
    }
}
