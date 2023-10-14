#include "extract_file.h"
#include "extract_file_from_stream.h"

unsigned char *extractFile(unsigned char *data, unsigned long in_size, unsigned long *out_size, int *out_type)
{

    ExtractFileFromStream effs;
    auto content = effs.GetFile(data, in_size, out_size, out_type);

    return content;
}

void deleteData(unsigned char *result)
{
    delete[] result;
}
