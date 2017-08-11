all: project build

build:
	@echo "\nAPP_PLATFORM := android-$(API)" >> jni/Application.mk
	ndk-build
	ant -q debug

install:
	ndk-build
	ant -q -Dadb.device.arg=$(ADFL) installd

project:
	# (required) go get github.com/xlab/android-go/cmd/android-project
	android-project update --target android-$(API) --name NuklearActivity --path .

clean:
	@printf "APP_ABI :=" > jni/Application.mk
	@rm -rf libs obj bin gen jni/arm* jni/x86* proguard-project.txt project.properties local.properties build.xml
