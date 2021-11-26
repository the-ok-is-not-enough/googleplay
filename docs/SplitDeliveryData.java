public final class SplitDeliveryData extends com.google.protobuf.nano.MessageNano {
    private static volatile com.google.android.finsky.protos.SplitDeliveryData[] _emptyArray;
    public long downloadSize;
    public String downloadUrl;
    public long gzippedDownloadSize;
    public String gzippedDownloadUrl;
    public boolean hasDownloadSize;
    public boolean hasDownloadUrl;
    public boolean hasGzippedDownloadSize;
    public boolean hasGzippedDownloadUrl;
    public boolean hasId;
    public boolean hasSignature;
    public String id;
    public com.google.android.finsky.protos.AndroidAppPatchData patchData;
    public String signature;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p7)
    {
        if ((this.hasId) || (!this.id.equals(""))) {
            p7.writeString(1, this.id);
        }
        if ((this.hasDownloadSize) || (this.downloadSize != 0)) {
            p7.writeInt64(2, this.downloadSize);
        }
        if ((this.hasGzippedDownloadSize) || (this.gzippedDownloadSize != 0)) {
            p7.writeInt64(3, this.gzippedDownloadSize);
        }
        if ((this.hasSignature) || (!this.signature.equals(""))) {
            p7.writeString(4, this.signature);
        }
        if ((this.hasDownloadUrl) || (!this.downloadUrl.equals(""))) {
            p7.writeString(5, this.downloadUrl);
        }
        if ((this.hasGzippedDownloadUrl) || (!this.gzippedDownloadUrl.equals(""))) {
            p7.writeString(6, this.gzippedDownloadUrl);
        }
        if (this.patchData != null) {
            p7.writeMessage(7, this.patchData);
        }
        super.writeTo(p7);
        return;
    }
}
