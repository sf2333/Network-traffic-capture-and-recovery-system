#include "ZipModule.h"



ZipModule::ZipModule()
= default;


ZipModule::~ZipModule()
= default;


int ZipModule::ParseZipEnd(const byte* zip_end) {

	//解析未压缩大小
	byte length_str[4];
	length_str[2] = zip_end[20];
	length_str[3] = zip_end[21];
	const auto other_length = ParseLength2(length_str);
	if (other_length < 0) {
		return 0;
	}
	return other_length;
}

//解析2字节的长度
int ZipModule::ParseLength2(byte* length) {
	return length[0] + length[1] * 256;
}

byte* ZipModule::GetIdentifyHead() {
	return this->identify_head_;
}

byte* ZipModule::GetIdentifyCore() {
	return this->identify_core_;
}

byte* ZipModule::GetIdentifyEnd() {
	return this->identify_end_;
}
