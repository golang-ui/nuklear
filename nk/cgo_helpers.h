// THE AUTOGENERATED LICENSE. ALL THE RIGHTS ARE RESERVED BY ROBOTS.

// WARNING: This file has automatically been generated on Mon, 02 Sep 2019 19:08:38 CST.
// Code generated by https://git.io/c-for-go. DO NOT EDIT.

#include "nuklear.h"
#include <stdlib.h>
#pragma once

#define __CGOGEN 1

// nk_plugin_alloc_88237ac9 is a proxy for callback nk_plugin_alloc.
void* nk_plugin_alloc_88237ac9(nk_handle arg0, void* old, unsigned long int arg2);

// nk_plugin_free_9e32bb09 is a proxy for callback nk_plugin_free.
void nk_plugin_free_9e32bb09(nk_handle arg0, void* old);

// nk_plugin_filter_1df5f22c is a proxy for callback nk_plugin_filter.
int nk_plugin_filter_1df5f22c(struct nk_text_edit* arg0, unsigned int unicode);

// nk_plugin_paste_70e696c4 is a proxy for callback nk_plugin_paste.
void nk_plugin_paste_70e696c4(nk_handle arg0, struct nk_text_edit* arg1);

// nk_plugin_copy_9ea6c143 is a proxy for callback nk_plugin_copy.
void nk_plugin_copy_9ea6c143(nk_handle arg0, char* arg1, int len);

// nk_text_width_f_67477c0 is a proxy for callback nk_text_width_f.
float nk_text_width_f_67477c0(nk_handle arg0, float h, char* arg2, int len);

// nk_query_font_glyph_f_5ba87240 is a proxy for callback nk_query_font_glyph_f.
void nk_query_font_glyph_f_5ba87240(nk_handle handle, float font_height, struct nk_user_font_glyph* glyph, unsigned int codepoint, unsigned int next_codepoint);

// nk_command_custom_callback_d451fdb1 is a proxy for callback nk_command_custom_callback.
void nk_command_custom_callback_d451fdb1(void* canvas, short int x, short int y, unsigned short int w, unsigned short int h, nk_handle callback_data);

