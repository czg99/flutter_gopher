# The Flutter tooling requires that developers have CMake 3.10 or later
# installed. You should not increase this version, as doing so will cause
# the plugin to fail to compile for some customers of the plugin.
cmake_minimum_required(VERSION 3.10)

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

# Execute build.sh shell script
message(STATUS "Executing build.sh shell script...")
execute_process(
    COMMAND bash "${CMAKE_CURRENT_SOURCE_DIR}/build.sh"
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    RESULT_VARIABLE SH_RESULT
    OUTPUT_VARIABLE SH_OUTPUT
    ERROR_VARIABLE SH_ERROR
)

# Check shell script execution result
if(NOT SH_RESULT EQUAL 0)
    message(STATUS "Shell script output: ${SH_OUTPUT}")
    message(STATUS "Shell script error: ${SH_ERROR}")
    message(FATAL_ERROR "build.sh execution failed, return code: ${SH_RESULT}")
else()
    message(STATUS "build.sh executed successfully")
endif()

# List of absolute paths to libraries that should be bundled with the plugin.
set({{.ProjectName}}_bundled_libraries
  "${CMAKE_CURRENT_SOURCE_DIR}/lib{{.LibName}}.so"
  PARENT_SCOPE
)
