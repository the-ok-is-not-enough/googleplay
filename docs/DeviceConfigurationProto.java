public final class DeviceConfiguration$DeviceConfigurationProto extends com.google.protobuf.nano.MessageNano {
    public int glEsVersion;
    public String[] glExtension;
    public boolean hasFiveWayNavigation;
    public boolean hasGlEsVersion;
    public boolean hasHardKeyboard;
    public boolean hasHasFiveWayNavigation;
    public boolean hasHasHardKeyboard;
    public boolean hasKeyboard;
    public boolean hasLowRamDevice;
    public boolean hasMaxApkDownloadSizeMb;
    public boolean hasMaxNumOfCpuCores;
    public boolean hasNavigation;
    public boolean hasScreenDensity;
    public boolean hasScreenHeight;
    public boolean hasScreenLayout;
    public boolean hasScreenWidth;
    public boolean hasSmallestScreenWidthDp;
    public boolean hasTotalMemoryBytes;
    public boolean hasTouchScreen;
    public int keyboard;
    public boolean lowRamDevice;
    public int maxApkDownloadSizeMb;
    public int maxNumOfCpuCores;
    public String[] nativePlatform;
    public int navigation;
    public int screenDensity;
    public int screenHeight;
    public int screenLayout;
    public int screenWidth;
    public int smallestScreenWidthDp;
    public String[] systemAvailableFeature;
    public String[] systemSharedLibrary;
    public String[] systemSupportedLocale;
    public long totalMemoryBytes;
    public int touchScreen;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p7)
    {
        if ((this.touchScreen != 0) || (this.hasTouchScreen)) {
            p7.writeInt32(1, this.touchScreen);
        }
        if ((this.keyboard != 0) || (this.hasKeyboard)) {
            p7.writeInt32(2, this.keyboard);
        }
        if ((this.navigation != 0) || (this.hasNavigation)) {
            p7.writeInt32(3, this.navigation);
        }
        if ((this.screenLayout != 0) || (this.hasScreenLayout)) {
            p7.writeInt32(4, this.screenLayout);
        }
        if ((this.hasHasHardKeyboard) || (this.hasHardKeyboard)) {
            p7.writeBool(5, this.hasHardKeyboard);
        }
        if ((this.hasHasFiveWayNavigation) || (this.hasFiveWayNavigation)) {
            p7.writeBool(6, this.hasFiveWayNavigation);
        }
        if ((this.hasScreenDensity) || (this.screenDensity != 0)) {
            p7.writeInt32(7, this.screenDensity);
        }
        if ((this.hasGlEsVersion) || (this.glEsVersion != 0)) {
            p7.writeInt32(8, this.glEsVersion);
        }
        if ((this.systemSharedLibrary != null) && (this.systemSharedLibrary.length > 0)) {
            int v1_4 = 0;
            while (v1_4 < this.systemSharedLibrary.length) {
                String v0_4 = this.systemSharedLibrary[v1_4];
                if (v0_4 != null) {
                    p7.writeString(9, v0_4);
                }
                v1_4++;
            }
        }
        if ((this.systemAvailableFeature != null) && (this.systemAvailableFeature.length > 0)) {
            int v1_3 = 0;
            while (v1_3 < this.systemAvailableFeature.length) {
                String v0_3 = this.systemAvailableFeature[v1_3];
                if (v0_3 != null) {
                    p7.writeString(10, v0_3);
                }
                v1_3++;
            }
        }
        if ((this.nativePlatform != null) && (this.nativePlatform.length > 0)) {
            int v1_2 = 0;
            while (v1_2 < this.nativePlatform.length) {
                String v0_2 = this.nativePlatform[v1_2];
                if (v0_2 != null) {
                    p7.writeString(11, v0_2);
                }
                v1_2++;
            }
        }
        if ((this.hasScreenWidth) || (this.screenWidth != 0)) {
            p7.writeInt32(12, this.screenWidth);
        }
        if ((this.hasScreenHeight) || (this.screenHeight != 0)) {
            p7.writeInt32(13, this.screenHeight);
        }
        if ((this.systemSupportedLocale != null) && (this.systemSupportedLocale.length > 0)) {
            int v1_0 = 0;
            while (v1_0 < this.systemSupportedLocale.length) {
                String v0_1 = this.systemSupportedLocale[v1_0];
                if (v0_1 != null) {
                    p7.writeString(14, v0_1);
                }
                v1_0++;
            }
        }
        if ((this.glExtension != null) && (this.glExtension.length > 0)) {
            int v1_1 = 0;
            while (v1_1 < this.glExtension.length) {
                String v0_0 = this.glExtension[v1_1];
                if (v0_0 != null) {
                    p7.writeString(15, v0_0);
                }
                v1_1++;
            }
        }
        if ((this.hasMaxApkDownloadSizeMb) || (this.maxApkDownloadSizeMb != 50)) {
            p7.writeInt32(17, this.maxApkDownloadSizeMb);
        }
        if ((this.hasSmallestScreenWidthDp) || (this.smallestScreenWidthDp != 0)) {
            p7.writeInt32(18, this.smallestScreenWidthDp);
        }
        if ((this.hasLowRamDevice) || (this.lowRamDevice)) {
            p7.writeBool(19, this.lowRamDevice);
        }
        if ((this.hasTotalMemoryBytes) || (this.totalMemoryBytes != 0)) {
            p7.writeInt64(20, this.totalMemoryBytes);
        }
        if ((this.hasMaxNumOfCpuCores) || (this.maxNumOfCpuCores != 0)) {
            p7.writeInt32(21, this.maxNumOfCpuCores);
        }
        super.writeTo(p7);
        return;
    }
}
