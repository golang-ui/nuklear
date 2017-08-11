ARCH ?= arm
API ?= 22
TOOLS := toolchain-$(API)/$(ARCH)
TOOLS_DIR ?= $(shell pwd)/$(TOOLS)
ANDROID_SYSROOT = $(NDK)/platforms/android-$(API)/arch-$(ARCH)
GLES_VER ?= gles3 # x86 emulator only supported gles2 as of june'17

ifeq ($(ARCH),arm)
  ABI  := armeabi-v7a
  GOCC := arm-linux-androideabi
  GOAR := arm
  GOFL := -march=armv7-a -std=c99
else ifeq ($(ARCH),arm64)
  ABI  := arm64-v8a
  GOCC := aarch64-linux-android
  GOAR := arm64
  GOFL := -march=armv8-a -std=c99
else ifeq ($(ARCH),x86)
  ABI  := x86
  GOCC := i686-linux-android
  GOAR := 386
  GOFL := -march=bonnell
else ifeq ($(ARCH),x86_64)
  ABI  := x86_64
  GOCC := x86_64-linux-android
  GOAR := amd64
  GOFL := -march=bonnell
endif

all: $(TOOLS) build

$(TOOLS):
	$(NDK)/build/tools/make_standalone_toolchain.py --force \
		--api=$(API) --install-dir=$(TOOLS_DIR) \
		--arch=$(ARCH) --stl libc++

build:
	@mkdir -p android/jni/$(ABI)
	@printf " $(ABI)" >> android/jni/Application.mk
	CC="$(TOOLS_DIR)/bin/$(GOCC)-gcc" \
	CXX="$(TOOLS_DIR)/bin/$(GOCC)-g++" \
	CGO_CFLAGS="$(GOFL)" \
	GOOS=android \
	GOARCH=$(GOAR) \
	GOARM=7 \
	CGO_ENABLED=1 \
	go build -tags $(GLES_VER) -buildmode=c-shared -o android/jni/$(ABI)/libnkactivity.so
 
apk:
	cd android && make API=$(API)

clean:
	cd android && make clean

purge: clean
	rm -rf toolchain*

install%:
	cd android && make ADFL=-$* install

listen%:
	@adb -$* logcat -c
	@adb -$* logcat '*:E' NuklearActivity
