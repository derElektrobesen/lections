project(pass)
cmake_minimum_required (VERSION 2.8)

set(SOURCES main.c pass.c)
set(CMAKE_C_FLAGS "-g -O3 -Wall -Werror -Wextra")
add_executable(binary ${SOURCES})
