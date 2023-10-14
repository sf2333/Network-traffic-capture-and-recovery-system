#pragma once
#include<Windows.h>
#include<set>
#include<iostream>
#define byte unsigned char

class JpgModule {
    private:
        byte file_head_[3] = { 0xff, 0xd8 };
        std::set<int> my_set_;
        int file_size_ = 0;

    public:
        JpgModule();
        ~JpgModule();

        const static int identity = 2;

        const static int head_size = 2;

        // void setLocation(DWORD location);

        byte* GetFileHead();

        int GetFileSize();
        void SetFileSize(int file_size);

        bool ParseFile(const byte* buffer, int length);
};

