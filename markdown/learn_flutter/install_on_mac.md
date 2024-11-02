## MAC安装 Flutter 环境

全文参考：https://doc.flutterchina.club/setup-macos/

### 下载Flutter sdk 1.5G

https://docs.flutter.dev/release/archive#macos

解压到指定目录：

```shell
cd ~/Downloads && mkdir -p ~/sdk
unzip -q flutter_macos_arm64_3.24.3-stable.zip -d ~/sdk/
```

设置ENV：

```shell
echo 'export PUB_HOSTED_URL=https://pub.flutter-io.cn' >> ~/.zshrc
echo 'export FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn' >> ~/.zshrc
echo 'export PATH=~/sdk/flutter/bin:$PATH' >> ~/.zshrc
source ~/.zshrc

flutter --disable-analytics
```

### 安装缺少的依赖

```shell
> flutter doctor
Flutter assets will be downloaded from https://storage.flutter-io.cn. Make sure you trust this source!
Doctor summary (to see all details, run flutter doctor -v):
[✓] Flutter (Channel stable, 3.24.3, on macOS 13.0.1 22A400 darwin-arm64, locale zh-Hans-CN)
[✗] Android toolchain - develop for Android devices
    ✗ Unable to locate Android SDK.
      Install Android Studio from: https://developer.android.com/studio/index.html
      On first launch it will assist you in installing the Android SDK components.
      (or visit https://flutter.dev/to/macos-android-setup for detailed instructions).
      If the Android SDK has been installed to a custom location, please use
      `flutter config --android-sdk` to update to that location.

[✗] Xcode - develop for iOS and macOS
    ✗ Xcode installation is incomplete; a full installation is necessary for iOS and macOS development.
      Download at: https://developer.apple.com/xcode/
      Or install Xcode via the App Store.
      Once installed, run:
        sudo xcode-select --switch /Applications/Xcode.app/Contents/Developer
        sudo xcodebuild -runFirstLaunch
    ✗ CocoaPods not installed.
        CocoaPods is a package manager for iOS or macOS platform code.
        Without CocoaPods, plugins will not work on iOS or macOS.
        For more info, see https://flutter.dev/to/platform-plugins
      For installation instructions, see https://guides.cocoapods.org/using/getting-started.html#installation
[✓] Chrome - develop for the web
[!] Android Studio (not installed)
[✓] Connected device (2 available)
[!] Network resources
    ✗ A cryptographic error occurred while checking "https://cocoapods.org/": Connection terminated during handshake
      You may be experiencing a man-in-the-middle attack, your network may be compromised, or you may have malware installed on your computer.

! Doctor found issues in 4 categories.
```

如果要为iOS开发flutter应用，则安装Xcode，否则安装Android Studio。

> 只是取决于你的调试设备是哪种平台，方便开发，代码可运行在多平台。

#### 安装 Xcode

只能在MacOS上安装Xcode。

#### 安装 Android Studio

https://developer.android.com/studio?hl=zh-cn

安装完成后，Android SDK就好了。然后打开IDE，插件里面安装 Dart、Flutter 插件，否则新建项目时看不到New Flutter Project按钮。

#### 处理 Android 工具链

opening Android Studio and going to SDK Manager, switch to the SDK Tools tab and check Android
SDK Command-line Tools (latest).

```shell
flutter doctor --android-licenses
```

### 使用注意

#### 1. 强退Android Studio可能导致flutter bin lock ，需要删除：

```shell
rm -rf ~/sdk/flutter/bin/cache/lockfile
```

#### 2. 启动main.dart报错：“Entrypoint isn't within the current project”

删除根目录下的`.idea`文件，重启IDE。
