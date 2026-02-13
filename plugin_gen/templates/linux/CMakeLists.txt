# The Flutter tooling requires that developers have CMake 3.10 or later
# installed. You should not increase this version, as doing so will cause
# the plugin to fail to compile for some customers of the plugin.
cmake_minimum_required(VERSION 3.10)

# Project-level configuration.
set(PROJECT_NAME "{{.ProjectName}}")
project(${PROJECT_NAME} LANGUAGES CXX)

# List of absolute paths to libraries that should be bundled with the plugin.
set({{.ProjectName}}_bundled_libraries
  "${CMAKE_CURRENT_SOURCE_DIR}/lib{{.LibName}}.so"
  PARENT_SCOPE
)
