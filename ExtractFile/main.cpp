//#include<iostream>
//#include<fstream>
//#include<Windows.h>
//#include "extract_file_from_stream.h"
//using namespace std;
//
//
//int main()
//{
//
//   cout<<"11111"<<endl;
//
//   ifstream infile;
////   infile.open("../../hex/png.hex", ios::binary);
//   infile.open("../jpg.hex", ios::binary);
//   // infile.open("i:\\new_word.hex", ios::binary);
//   if (!infile)
//       return 0;
//   infile.seekg(0, ios::end);
//   const unsigned long file_size = infile.tellg();
//   const auto buffer = new BYTE[file_size];
//   infile.seekg(0, ios::beg);
//   infile.read(reinterpret_cast<char*>(buffer), file_size);
//   infile.close();
//
//   const ExtractFileFromStream effs;
//   ULONG out_size = 0;
//   int out_type = 0;
//   const auto content = effs.GetFile(buffer, file_size, &out_size,&out_type);
//   delete[] content;
//   delete[] buffer;
//
//   return 0;
//}