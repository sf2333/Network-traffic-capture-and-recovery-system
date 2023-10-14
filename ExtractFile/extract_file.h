#ifndef __EXTRACT_FILE_
#define __EXTRACT_FILE_

#ifdef __cplusplus
extern "C"{
#endif

unsigned char* extractFile(unsigned char* data, unsigned long in_size, unsigned long *out_size, int *out_type);

void deleteData(unsigned char* result);

#ifdef __cplusplus
}
#endif





#endif