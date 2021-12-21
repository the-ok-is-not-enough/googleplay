public final class Response$ResponseWrapper extends com.google.protobuf.nano.MessageNano {
    public com.google.android.finsky.protos.ServerCommands commands;
    public com.google.android.finsky.protos.Notification[] notification;
    public com.google.android.finsky.protos.Response$Payload payload;
    public com.google.android.finsky.protos.PreFetch[] preFetch;
    public com.google.android.finsky.protos.ServerCookies serverCookies;
    public com.google.android.finsky.protos.ServerMetadata serverMetadata;
    public com.google.android.finsky.protos.Targets targets;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p5)
    {
        if (this.payload != null) {
            p5.writeMessage(1, this.payload);
        }
        if (this.commands != null) {
            p5.writeMessage(2, this.commands);
        }
        if ((this.preFetch != null) && (this.preFetch.length > 0)) {
            int v1_0 = 0;
            while (v1_0 < this.preFetch.length) {
                com.google.android.finsky.protos.Notification v0_1 = this.preFetch[v1_0];
                if (v0_1 != null) {
                    p5.writeMessage(3, v0_1);
                }
                v1_0++;
            }
        }
        if ((this.notification != null) && (this.notification.length > 0)) {
            int v1_1 = 0;
            while (v1_1 < this.notification.length) {
                com.google.android.finsky.protos.Notification v0_0 = this.notification[v1_1];
                if (v0_0 != null) {
                    p5.writeMessage(4, v0_0);
                }
                v1_1++;
            }
        }
        if (this.serverMetadata != null) {
            p5.writeMessage(5, this.serverMetadata);
        }
        if (this.targets != null) {
            p5.writeMessage(6, this.targets);
        }
        if (this.serverCookies != null) {
            p5.writeMessage(7, this.serverCookies);
        }
        super.writeTo(p5);
        return;
    }
}
