g++ -c AhoCorasick.cpp
g++ -c PngModule.cpp
g++ -c extract_file.cpp
g++ -c extract_file_from_stream.cpp
g++ -c JpgModule.cpp
g++ -c ZipModule.cpp


ar -crs libextract.a AhoCorasick.o PngModule.o extract_file_from_stream.o extract_file.o JpgModule.o ZipModule.o


