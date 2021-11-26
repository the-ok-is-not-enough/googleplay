public final class AppDetails extends com.google.protobuf.nano.MessageNano {
    public String[] appCategory;
    public String appType;
    public String[] autoAcquireFreeAppIfHigherVersionAvailableTag;
    public String[] certificateHash;
    public com.google.android.finsky.protos.CertificateSet[] certificateSet;
    public int contentRating;
    public boolean declaresIab;
    public String developerEmail;
    public String developerName;
    public String developerWebsite;
    public boolean everExternallyHosted;
    public boolean externallyHosted;
    public com.google.android.finsky.protos.FileMetadata[] file;
    public boolean gamepadRequired;
    public boolean hasAppType;
    public boolean hasContentRating;
    public boolean hasDeclaresIab;
    public boolean hasDeveloperEmail;
    public boolean hasDeveloperName;
    public boolean hasDeveloperWebsite;
    public boolean hasEverExternallyHosted;
    public boolean hasExternallyHosted;
    public boolean hasGamepadRequired;
    public boolean hasInstallLocation;
    public boolean hasInstallNotes;
    public boolean hasInstallationSize;
    public boolean hasMajorVersionNumber;
    public boolean hasNumDownloads;
    public boolean hasPackageName;
    public boolean hasPreregistrationPromoCode;
    public boolean hasRecentChangesHtml;
    public boolean hasTargetSdkVersion;
    public boolean hasTitle;
    public boolean hasUploadDate;
    public boolean hasVariesByAccount;
    public boolean hasVersionCode;
    public boolean hasVersionString;
    public com.google.android.finsky.protos.InstallDetails installDetails;
    public int installLocation;
    public String installNotes;
    public long installationSize;
    public int majorVersionNumber;
    public String numDownloads;
    public String packageName;
    public String[] permission;
    public String preregistrationPromoCode;
    public String recentChangesHtml;
    public String[] splitId;
    public int targetSdkVersion;
    public String title;
    public String uploadDate;
    public boolean variesByAccount;
    public int versionCode;
    public String versionString;
    public final void writeTo(com.google.protobuf.nano.CodedOutputByteBufferNano p8)
    {
        if ((this.hasDeveloperName) || (!this.developerName.equals(""))) {
            p8.writeString(1, this.developerName);
        }
        if ((this.hasMajorVersionNumber) || (this.majorVersionNumber != 0)) {
            p8.writeInt32(2, this.majorVersionNumber);
        }
        if ((this.hasVersionCode) || (this.versionCode != 0)) {
            p8.writeInt32(3, this.versionCode);
        }
        if ((this.hasVersionString) || (!this.versionString.equals(""))) {
            p8.writeString(4, this.versionString);
        }
        if ((this.hasTitle) || (!this.title.equals(""))) {
            p8.writeString(5, this.title);
        }
        if ((this.appCategory != null) && (this.appCategory.length > 0)) {
            int v1_6 = 0;
            while (v1_6 < this.appCategory.length) {
                String v0_6 = this.appCategory[v1_6];
                if (v0_6 != null) {
                    p8.writeString(7, v0_6);
                }
                v1_6++;
            }
        }
        if ((this.hasContentRating) || (this.contentRating != 0)) {
            p8.writeInt32(8, this.contentRating);
        }
        if ((this.hasInstallationSize) || (this.installationSize != 0)) {
            p8.writeInt64(9, this.installationSize);
        }
        if ((this.permission != null) && (this.permission.length > 0)) {
            int v1_5 = 0;
            while (v1_5 < this.permission.length) {
                String v0_5 = this.permission[v1_5];
                if (v0_5 != null) {
                    p8.writeString(10, v0_5);
                }
                v1_5++;
            }
        }
        if ((this.hasDeveloperEmail) || (!this.developerEmail.equals(""))) {
            p8.writeString(11, this.developerEmail);
        }
        if ((this.hasDeveloperWebsite) || (!this.developerWebsite.equals(""))) {
            p8.writeString(12, this.developerWebsite);
        }
        if ((this.hasNumDownloads) || (!this.numDownloads.equals(""))) {
            p8.writeString(13, this.numDownloads);
        }
        if ((this.hasPackageName) || (!this.packageName.equals(""))) {
            p8.writeString(14, this.packageName);
        }
        if ((this.hasRecentChangesHtml) || (!this.recentChangesHtml.equals(""))) {
            p8.writeString(15, this.recentChangesHtml);
        }
        if ((this.hasUploadDate) || (!this.uploadDate.equals(""))) {
            p8.writeString(16, this.uploadDate);
        }
        if ((this.file != null) && (this.file.length > 0)) {
            int v1_0 = 0;
            while (v1_0 < this.file.length) {
                String v0_4 = this.file[v1_0];
                if (v0_4 != null) {
                    p8.writeMessage(17, v0_4);
                }
                v1_0++;
            }
        }
        if ((this.hasAppType) || (!this.appType.equals(""))) {
            p8.writeString(18, this.appType);
        }
        if ((this.certificateHash != null) && (this.certificateHash.length > 0)) {
            int v1_1 = 0;
            while (v1_1 < this.certificateHash.length) {
                String v0_3 = this.certificateHash[v1_1];
                if (v0_3 != null) {
                    p8.writeString(19, v0_3);
                }
                v1_1++;
            }
        }
        if ((this.hasVariesByAccount) || (this.variesByAccount != 1)) {
            p8.writeBool(21, this.variesByAccount);
        }
        if ((this.certificateSet != null) && (this.certificateSet.length > 0)) {
            int v1_2 = 0;
            while (v1_2 < this.certificateSet.length) {
                String v0_2 = this.certificateSet[v1_2];
                if (v0_2 != null) {
                    p8.writeMessage(22, v0_2);
                }
                v1_2++;
            }
        }
        if ((this.autoAcquireFreeAppIfHigherVersionAvailableTag != null) && (this.autoAcquireFreeAppIfHigherVersionAvailableTag.length > 0)) {
            int v1_3 = 0;
            while (v1_3 < this.autoAcquireFreeAppIfHigherVersionAvailableTag.length) {
                String v0_1 = this.autoAcquireFreeAppIfHigherVersionAvailableTag[v1_3];
                if (v0_1 != null) {
                    p8.writeString(23, v0_1);
                }
                v1_3++;
            }
        }
        if ((this.hasDeclaresIab) || (this.declaresIab)) {
            p8.writeBool(24, this.declaresIab);
        }
        if ((this.splitId != null) && (this.splitId.length > 0)) {
            int v1_4 = 0;
            while (v1_4 < this.splitId.length) {
                String v0_0 = this.splitId[v1_4];
                if (v0_0 != null) {
                    p8.writeString(25, v0_0);
                }
                v1_4++;
            }
        }
        if ((this.hasGamepadRequired) || (this.gamepadRequired)) {
            p8.writeBool(26, this.gamepadRequired);
        }
        if ((this.hasExternallyHosted) || (this.externallyHosted)) {
            p8.writeBool(27, this.externallyHosted);
        }
        if ((this.hasEverExternallyHosted) || (this.everExternallyHosted)) {
            p8.writeBool(28, this.everExternallyHosted);
        }
        if ((this.hasInstallNotes) || (!this.installNotes.equals(""))) {
            p8.writeString(30, this.installNotes);
        }
        if ((this.installLocation != 0) || (this.hasInstallLocation)) {
            p8.writeInt32(31, this.installLocation);
        }
        if ((this.hasTargetSdkVersion) || (this.targetSdkVersion != 0)) {
            p8.writeInt32(32, this.targetSdkVersion);
        }
        if ((this.hasPreregistrationPromoCode) || (!this.preregistrationPromoCode.equals(""))) {
            p8.writeString(33, this.preregistrationPromoCode);
        }
        if (this.installDetails != null) {
            p8.writeMessage(34, this.installDetails);
        }
        super.writeTo(p8);
        return;
    }
}
