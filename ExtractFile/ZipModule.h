#pragma once
#include<iostream>
#include<Windows.h>
#include<string>
using namespace std;
#define byte unsigned char
/*
用来匹配office文件的zip对象
*/
class ZipModule {
    private:
        byte identify_head_[5] = { 0x50, 0x4b, 0x03, 0x04, 0x00 }; //0-3bit
        byte identify_end_[5] = { 0x50, 0x4b, 0x05, 0x06, 0x00 };
        byte identify_core_[5] = { 0x50, 0x4b, 0x01, 0x02, 0x00 };

    public:
        const static int head_size = 5;
        const static int core_size = 5;
        const static int end_size = 5;
		const static int head_identity = 5; //在ac自动机算法中标识
		const static int core_identity = 6; //在ac自动机算法中标识
		const static int end_identity = 7; //在ac自动机算法中标识

    public:
        ZipModule();
        ~ZipModule();

        //解析zip文件头部

        static int ParseZipEnd(const byte* zip_end);
        static int ParseLength2(byte* length);

        byte* GetIdentifyHead();
        byte* GetIdentifyCore();
        byte* GetIdentifyEnd();

};

