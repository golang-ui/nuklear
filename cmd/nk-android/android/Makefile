all: project build 

build:
	ndk-build
	ant debug

install:
	ndk-build
	ant debug install

project:
	# (optional) android list targets
	android update project --target android-23 --name NuklearActivity --path .

project-n:
	android update project --target android-N --name NuklearActivity --path .

clean:
	rm -rf libs obj bin gen jni/lib
	rm -f proguard-project.txt project.properties
	rm -f local.properties build.xml
