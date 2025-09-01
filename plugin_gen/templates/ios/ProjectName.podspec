Pod::Spec.new do |s|
  s.name             = '{{.ProjectName}}'
  s.version          = '0.0.1'
  s.summary          = 'A new Flutter plugin project.'
  s.description      = <<-DESC
A new Flutter plugin project.
                       DESC
  s.homepage         = 'http://example.com'
  s.license          = { :file => '../LICENSE' }
  s.author           = { 'Your Company' => 'email@example.com' }
  s.source           = { :path => '.' }
  s.source_files = 'Classes/**/*'
  s.public_header_files = 'Classes/**/*.h'

  s.subspec 'NoARC' do |sp|
    sp.source_files = 'Classes/protos/**/*.{h,m}'
    sp.requires_arc = false
  end

  s.dependency 'Flutter'
  s.dependency 'Protobuf'

  s.platform = :ios, '11.0'
  # Flutter.framework does not contain a i386 slice.
  s.pod_target_xcconfig = { 'DEFINES_MODULE' => 'YES', 'EXCLUDED_ARCHS[sdk=iphonesimulator*]' => 'i386' }

  s.script_phases = [
    {
      :name => 'Run Pre-Build Script',
      :script => "sh '#{__dir__}/build.sh'",
      :execution_position => 'before_compile'
    }
  ]

  s.prepare_command = <<-CMD
    sh ./build.sh
  CMD

  s.vendored_frameworks = '{{.LibName}}.xcframework'
end
