public final class ServerCommands extends com.google.protobuf.nano.MessageNano {
    public boolean clearCache;
    public String displayErrorMessage;
    public boolean hasClearCache;
    public boolean hasDisplayErrorMessage;
    public boolean hasLogErrorStacktrace;
    public String logErrorStacktrace;
    public com.google.android.finsky.protos.UserSettingDirtyData[] userSettingDirtyData;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p5)
    {
        if ((this.hasClearCache) || (this.clearCache)) {
            p5.writeBool(1, this.clearCache);
        }
        if ((this.hasDisplayErrorMessage) || (!this.displayErrorMessage.equals(""))) {
            p5.writeString(2, this.displayErrorMessage);
        }
        if ((this.hasLogErrorStacktrace) || (!this.logErrorStacktrace.equals(""))) {
            p5.writeString(3, this.logErrorStacktrace);
        }
        if ((this.userSettingDirtyData != null) && (this.userSettingDirtyData.length > 0)) {
            int v1 = 0;
            while (v1 < this.userSettingDirtyData.length) {
                com.google.android.finsky.protos.UserSettingDirtyData v0 = this.userSettingDirtyData[v1];
                if (v0 != null) {
                    p5.writeMessage(4, v0);
                }
                v1++;
            }
        }
        super.writeTo(p5);
        return;
    }
}
