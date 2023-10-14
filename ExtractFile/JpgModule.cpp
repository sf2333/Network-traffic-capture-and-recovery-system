#include "JpgModule.h"



JpgModule::JpgModule() {
    //0xffe1 0xffe0~0xffef 0xffdb 0xffc0 0xffc4 0xffdd 0xffda 0xffd9
    for (int i = 0; i < 16; i++)
        my_set_.insert(0xffe0 + i);
    my_set_.insert(0xffdb);
    my_set_.insert(0xffc0);
    my_set_.insert(0xffc4);
    my_set_.insert(0xffdd);
    my_set_.insert(0xffda);
    my_set_.insert(0xffd9);

}


JpgModule::~JpgModule() {
}


byte* JpgModule::GetFileHead() {
    return this->file_head_;
}

int JpgModule::GetFileSize() {
    return this->file_size_;
}


void JpgModule::SetFileSize(const int file_size) {
    this->file_size_ = file_size;
}

bool JpgModule::ParseFile(const byte* buffer, const int length) {
    if (length < 2)
        return false;
    auto identity = buffer[2] * 256 + buffer[3];

    auto len1 = 4;
    auto len2 = 5;
    auto result = false;

//    std::cout << "匹配" << std::endl;

    while ((this->my_set_.count(identity) || identity == 0xffd9) && len2 < (length - 1)) {

        result = true;
        const auto this_length = buffer[len1] * 256 + buffer[len2];



        len1 += this_length;
        len2 += this_length;
        this->file_size_ += this_length + 2;


        auto now_length = 0;
        if (identity == 0xffda) {
            identity = buffer[len1] * 256 + buffer[len2];
            //std::cout << std::hex << int(buffer[len1]);
            //std::cout << "  " << std::hex << int(buffer[len2]) << std::endl;
            len1 += 2;
            len2 += 2;
            while ((now_length + len1 < (length - 2)) && (now_length < 0xffff)) {
                if (buffer[now_length + len1] == 0xff) {
                    if (buffer[now_length + len2] == 0xda) {
                        this->file_size_ += now_length + 2;
                        len1 += now_length + 2;
                        len2 += now_length + 2;
                    } else if (buffer[now_length + len2] == 0xd9) {
                        this->file_size_ += now_length;
                        return result;
                    }
                }
                now_length++;
            }
        } else {
            identity = buffer[len1] * 256 + buffer[len2];
            //std::cout << std::hex << int(buffer[len1]);
            //std::cout << "  " << std::hex << int(buffer[len2]) << std::endl;
            len1 += 2;
            len2 += 2;
        }
    }
    return result;
}
