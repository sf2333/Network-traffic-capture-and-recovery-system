#include "extract_file_from_stream.h"
#include "AhoCorasick.h"
#include "JpgModule.h"
#include "ZipModule.h"
#include "PngModule.h"
#include <iostream>
#include <fstream>
#pragma execution_character_set("utf-8")

ExtractFileFromStream::ExtractFileFromStream() {
//    std::cout << "ExtractFileFromStream create " << std::endl;
}

ExtractFileFromStream::~ExtractFileFromStream() {
//    std::cout << "ExtractFileFromStream delete " << std::endl;
}

BYTE *ExtractFileFromStream::GetFile(byte *source_data, const ULONG in_size,
                                     ULONG *out_size, int *out_type) const {

  //初始化AC自动机
  const auto root = new AcNode();

  //初始化zip文件格式自动机
  ZipModule word_header;
  auto str1 = new byte[ZipModule::head_size + 1];
  memcpy(str1, word_header.GetIdentifyHead(), ZipModule::head_size);
  str1[ZipModule::head_size] = '\0';
  AhoCoraSick::Insert(str1, root, ZipModule::head_identity);
  delete[] str1;

  //// 中间部分，可以不用
  // ZipModule word_core;
  // str1 = new byte[ZipModule::core_size + 1];
  // memcpy(str1, word_core.getIdentifyCore(), ZipModule::core_size);
  // str1[ZipModule::core_size] = '\0';
  // AhoCoraSick::Insert(str1, root, ZipModule::core_identity);
  // delete[] str1;

  // word 结尾
  ZipModule word_end;
  str1 = new byte[ZipModule::end_size + 1];
  memcpy(str1, word_end.GetIdentifyEnd(), ZipModule::end_size);
  str1[ZipModule::end_size] = '\0';
  AhoCoraSick::Insert(str1, root, ZipModule::end_identity);
  delete[] str1;

  //初始化jpg文件格式自动机
  JpgModule jpg_module;
  str1 = new byte[JpgModule::head_size + 1];
  memcpy(str1, jpg_module.GetFileHead(), JpgModule::head_size);
  str1[JpgModule::head_size] = '\0';
  AhoCoraSick::Insert(str1, root, JpgModule::identity);
  delete[] str1;

  //初始化png文件格式自动机
  PngModule png_module;
  str1 = new byte[PngModule::head_size + 1];
  memcpy(str1, png_module.GetFileIdentity(), PngModule::head_size);
  str1[PngModule::head_size] = '\0';
  AhoCoraSick::Insert(str1, root, PngModule::identity);
  delete[] str1;

  //构建ac自动机fail跳转
  AhoCoraSick::BuildAcFail(root);

  const auto dest_data = new byte[in_size];
  const auto extract_file_size =
      SearchZipIdentity(source_data, dest_data, in_size, root, out_type);
  *out_size = extract_file_size;

  // 测试写出的文件是否正确，你自己注释掉，不用的话
//   fstream file("I:\\old_word.docx", ios::out | ios::binary);
//
//   // fstream file();
//   if (!file)
//   {
//       cout << "Error opening file." << endl;
//       return dest_data;
//   }
//  fstream file;
//  file.open("../old_jpg.jpg", ios::out | ios::binary);
//  file.write(reinterpret_cast<char *>(dest_data), extract_file_size);
//  file.close();

  return dest_data;
}

int ExtractFileFromStream::SearchZipIdentity(byte *buffer,
                                             byte *dest_data,
                                             const int length,
                                             AcNode *root,
                                             int *out_type) const {
  auto ac_node = root;
  auto number = 0;
  long long int index = 9999999999999;
  for (auto i = 0; i < length; i++) {
    const auto str = buffer[i];

    while ((ac_node->next[str] == nullptr) && (ac_node != root)) //当前字符不为根节点且不匹配
      ac_node = ac_node->fail;//去当前节点失败指针所指向的字符继续匹配

    if (ac_node->next[str] != nullptr) {
      ac_node = ac_node->next[str];
      DWORD location = i;

      switch (ac_node->count) {
        case PngModule::identity: {
          location = location - PngModule::head_size + 1;
          PngModule png_module;
          png_module.SetLocation(location);
          //解析png文件长度
          if (png_module.ParePngHead(buffer + i - PngModule::head_size + 1,
                                     length - i + PngModule::head_size)) {
//                    std::cout << to_string(png_module.GetFileSize()) << std::endl;
            std::cout << "Png识别头的位置为：" << location << endl;
            *out_type = 3;
            number = png_module.GetFileSize();
            memcpy_s(dest_data,
                     length,
                     buffer + i - PngModule::head_size + 1,
                     number);
            return number;
          }
          break;
        }
        case JpgModule::identity: {
          location = location - JpgModule::head_size + 1;

          JpgModule jpg_module;
          // jpg_module.setLocation(location);
          //解析png文件长度
          if (jpg_module.ParseFile(buffer + i - JpgModule::head_size + 1,
                                   length - i + JpgModule::head_size)) {
//                    std::cout << to_string(jpg_module.GetFileSize()) << std::endl;
            std::cout << "jpg头的位置为：" << location << std::endl;

//            number = jpg_module.GetFileSize();
            number = length - i + JpgModule::head_size;
//            memcpy_s(dest_data, length, buffer + i - 1, length);
//            memcpy_s(dest_data, length, buffer + i - JpgModule::head_size + 1, number);
            memcpy_s(dest_data, length, buffer + i - JpgModule::head_size + 1, length - i + JpgModule::head_size);
            *out_type = 1;
            return number;
          }
          break;
        }
        case ZipModule::head_identity: {
          location = location - ZipModule::head_size + 1;
          if (location < index) {
            index = location;
            std::cout << "Word头的位置为：" << location << std::endl;
          }
          break;
        }
        case ZipModule::core_identity: {
          break;
        }
        case ZipModule::end_identity: {
          location = location - ZipModule::end_size + 1;
          byte zip_end[30];
          memcpy_s(zip_end, 30, buffer + location, 22);
          const auto len = ZipModule::ParseZipEnd(zip_end);
          auto end_position = 0;
          if (len) {
            end_position = location + 22 + end_position;
            std::cout << "Word尾的位置为：" << location << std::endl;
          } else
            end_position = location + 22;
          const auto f_size = end_position - index;
          memcpy_s(dest_data, f_size, buffer + index + 1, f_size);
          *out_type = 2;

          return f_size;
          break;
        }
//            default:
//                std::cout << "0x504b0506位置为：" << i << std::endl;
      }
    }
  }
  return number;
}
