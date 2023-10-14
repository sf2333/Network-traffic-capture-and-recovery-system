#include "PngModule.h"

PngModule::PngModule()
= default;

PngModule::~PngModule()
= default;

byte *PngModule::GetFileIdentity() {
  return this->file_identity_;
}

void PngModule::SetLocation(DWORD location) {
  this->location_ = location;
}

DWORD PngModule::GetLocation() {
  return this->location_;
}

bool PngModule::ParePngHead(byte *buffer, int length) {

  if (length < 24) {
    return false;
  }

  int height =
      buffer[16] * 256 * 256 * 256 + buffer[17] * 256 * 256 + buffer[18] * 256
          + buffer[19];
  int width =
      buffer[20] * 256 * 256 * 256 + buffer[21] * 256 * 256 + buffer[22] * 256
          + buffer[23];

  auto size = height * width;

  this->file_size_ = size > length ? length : size;

  return true;
}

int PngModule::GetFileSize() {
  return this->file_size_;
}
