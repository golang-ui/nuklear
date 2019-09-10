#include "nk.h"

extern void igClipboardPaste(nk_handle handle, struct nk_text_edit *text_edit);
extern void igClipboardCopy(nk_handle handle, const char *content, int len);

void nk_register_clipboard(struct nk_context *ctx)
{
    ctx->clip.copy = igClipboardCopy;
    ctx->clip.paste = igClipboardPaste;
}