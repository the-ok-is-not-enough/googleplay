// Lcom/google/android/finsky/protos/DeliveryResponse;
public final class DeliveryResponse extends com.google.protobuf.nano.MessageNano {
    public com.google.android.finsky.protos.AndroidAppDeliveryData appDeliveryData;
    public boolean hasStatus;
    public int status;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p3)
    {
        if ((this.status != 1) || (this.hasStatus)) {
            p3.writeInt32(1, this.status);
        }
        if (this.appDeliveryData != null) {
            p3.writeMessage(2, this.appDeliveryData);
        }
        super.writeTo(p3);
        return;
    }
}
