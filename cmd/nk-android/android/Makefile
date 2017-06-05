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

project-n:
	android-project update --target android-N --name NuklearActivity --path .

clean:
	@printf "APP_ABI :=" > jni/Application.mk
	@rm -rf libs obj bin gen jni/arm* jni/x86*
	@rm -f proguard-project.txt project.properties
	@rm -f local.properties build.xml
