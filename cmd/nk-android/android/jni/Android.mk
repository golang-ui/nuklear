LOCAL_PATH := $(call my-dir)

include $(CLEAR_VARS)

LOCAL_MODULE    := nk
LOCAL_SRC_FILES := $(TARGET_ARCH_ABI)/libnkactivity.so
LOCAL_LDLIBS    := -llog -landroid

include $(PREBUILT_SHARED_LIBRARY)
