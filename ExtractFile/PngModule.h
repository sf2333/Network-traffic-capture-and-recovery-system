#pragma once
#include<Windows.h>
#define byte unsigned char

class PngModule {
    private:
        byte file_identity_[9] = { 0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A };
        DWORD location_ = 0;//记录在磁盘中的偏移位置
        int file_size_ = 0;

    public:
        PngModule();
        ~PngModule();

        const static int identity = 1; //在ac自动机算法中标识
        const static int head_size = 8;

        byte* GetFileIdentity();


        void SetLocation(DWORD location);
        DWORD GetLocation();

        bool ParePngHead(byte* buffer, int length);

        int GetFileSize();

};

