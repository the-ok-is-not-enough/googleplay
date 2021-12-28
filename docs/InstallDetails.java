public final class InstallDetails extends com.google.protobuf.nano.MessageNano {
    public com.google.android.finsky.protos.Dependency[] dependency;
    public boolean hasInstallLocation;
    public boolean hasSize;
    public boolean hasTargetSdkVersion;
    public int installLocation;
    public long size;
    public int targetSdkVersion;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p7)
    {
        if ((this.installLocation != 0) || (this.hasInstallLocation)) {
            p7.writeInt32(1, this.installLocation);
        }
        if ((this.hasSize) || (this.size != 0)) {
            p7.writeInt64(2, this.size);
        }
        if ((this.dependency != null) && (this.dependency.length > 0)) {
            int v1 = 0;
            while (v1 < this.dependency.length) {
                com.google.android.finsky.protos.Dependency v0 = this.dependency[v1];
                if (v0 != null) {
                    p7.writeMessage(3, v0);
                }
                v1++;
            }
        }
        if ((this.hasTargetSdkVersion) || (this.targetSdkVersion != 0)) {
            p7.writeInt32(4, this.targetSdkVersion);
        }
        super.writeTo(p7);
        return;
    }
}
