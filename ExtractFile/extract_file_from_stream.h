#pragma once
#include "AhoCorasick.h"
#define byte unsigned char

class ExtractFileFromStream {
    public:
        ExtractFileFromStream();
        ~ExtractFileFromStream();
        // BYTE* GetFile(BYTE* source_data, ULONG file_size);
        byte* GetFile(byte *source_data, const ULONG in_size, ULONG *out_size, int *out_type) const;
        int SearchZipIdentity(byte* buffer, byte* dest_data, const int length, AcNode* root, int *out_type) const;


};

