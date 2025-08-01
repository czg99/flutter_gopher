Pod::Spec.new do |s|
  s.name             = '{{.ProjectName}}'
  s.version          = '0.0.1'
  s.summary          = 'A new Flutter project.'
  s.description      = <<-DESC
A new Flutter project.
                       DESC
  s.homepage         = 'http://example.com'
  s.license          = { :file => '../LICENSE' }
  s.author           = { 'Your Company' => 'email@example.com' }

  s.source           = { :path => '.' }
  s.source_files = 'Classes/**/*'
  s.public_header_files = 'Classes/**/*.h'
  s.dependency 'FlutterMacOS'
  s.dependency 'Protobuf'

  s.platform = :osx, '10.11'
  s.pod_target_xcconfig = { 'DEFINES_MODULE' => 'YES' }

  s.script_phases = [
    {
      :name => 'Run Pre-Build Script',
      :script => "sh '#{__dir__}/build_dylib.sh'",
      :execution_position => 'before_compile'
    }
  ]

  s.prepare_command = <<-CMD
    sh ./build_dylib.sh
  CMD

  s.vendored_libraries = 'lib{{.LibName}}.dylib'
end
