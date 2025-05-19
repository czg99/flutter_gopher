# The Flutter tooling requires that developers have a version of Visual Studio
# installed that includes CMake 3.14 or later. You should not increase this
# version, as doing so will cause the plugin to fail to compile for some
# customers of the plugin.
cmake_minimum_required(VERSION 3.14)

# Project-level configuration.
set(PROJECT_NAME "{{.ProjectName}}")
project(${PROJECT_NAME} LANGUAGES CXX)

# Check if Go is installed
find_program(GO_EXECUTABLE go)
if(NOT GO_EXECUTABLE)
    message(FATAL_ERROR "Go compiler not found. Please install Go.")
endif()

# Check for Zig compiler
find_program(ZIG_EXECUTABLE zig)
if(NOT ZIG_EXECUTABLE)
    message(FATAL_ERROR "Zig compiler not found. Will use system default compiler.")
endif()

# Execute build.ps1 PowerShell script
message(STATUS "Executing build.ps1 PowerShell script...")
execute_process(
    COMMAND powershell.exe -ExecutionPolicy Bypass -File "${CMAKE_CURRENT_SOURCE_DIR}/build.ps1"
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    RESULT_VARIABLE CMD_RESULT
    OUTPUT_VARIABLE CMD_OUTPUT
    ERROR_VARIABLE CMD_ERROR
)

# Check script execution result
if(NOT CMD_RESULT EQUAL 0)
    message(STATUS "PowerShell output: ${CMD_OUTPUT}")
    message(STATUS "PowerShell error: ${CMD_ERROR}")
    message(FATAL_ERROR "build.ps1 execution failed, return code: ${CMD_RESULT}")
else()
    message(STATUS "build.ps1 executed successfully")
endif()

# List of absolute paths to libraries that should be bundled with the plugin.
set({{.ProjectName}}_bundled_libraries
  "${CMAKE_CURRENT_SOURCE_DIR}/{{.LibName}}.dll"
  PARENT_SCOPE
)
