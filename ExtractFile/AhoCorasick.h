#pragma once
#include<Windows.h>
#include<iostream>
#include<queue>
#include<map>
#include<string>
#include<fstream>

#define byte unsigned char

using namespace std;

struct AcNode {
    AcNode* fail;
    AcNode* next[256]{};
    byte str;
    int count;

    AcNode();
};

class AhoCoraSick {
    public:
        //建立时插入字符串，构建树
        static void Insert(byte* str, AcNode* root, int identity);

        //构建fail
        static void BuildAcFail(AcNode* root);
        static int SearchKeyword(const byte* buffer, int length, AcNode* root);

        //int search_keyword(byte* buffer, int length, AcNode* root);
};
