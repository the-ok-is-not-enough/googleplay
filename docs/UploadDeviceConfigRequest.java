public final class UploadDeviceConfig$UploadDeviceConfigRequest extends com.google.protobuf.nano.MessageNano {
    public com.google.android.finsky.protos.DataServiceSubscriber dataServiceSubscriber;
    public com.google.android.finsky.protos.DeviceConfiguration$DeviceConfigurationProto deviceConfiguration;
    public String gcmRegistrationId;
    public boolean hasGcmRegistrationId;
    public boolean hasManufacturer;
    public String manufacturer;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p3)
    {
        if (this.deviceConfiguration != null) {
            p3.writeMessage(1, this.deviceConfiguration);
        }
        if ((this.hasManufacturer) || (!this.manufacturer.equals(""))) {
            p3.writeString(2, this.manufacturer);
        }
        if ((this.hasGcmRegistrationId) || (!this.gcmRegistrationId.equals(""))) {
            p3.writeString(3, this.gcmRegistrationId);
        }
        if (this.dataServiceSubscriber != null) {
            p3.writeMessage(4, this.dataServiceSubscriber);
        }
        super.writeTo(p3);
        return;
    }
}
