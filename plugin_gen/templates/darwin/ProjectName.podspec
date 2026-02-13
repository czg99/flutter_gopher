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

  s.subspec 'NoARC' do |sp|
    sp.source_files = 'Classes/protos/**/*.{h,m}'
    sp.requires_arc = false
  end

  s.ios.dependency 'Flutter'
  s.osx.dependency 'FlutterMacOS'

  s.dependency 'Protobuf'

  s.ios.deployment_target = '11.0'
  s.osx.deployment_target = '10.11'

  s.pod_target_xcconfig = { 'DEFINES_MODULE' => 'YES' }

  s.ios.vendored_frameworks = '{{.LibName}}.xcframework'
  s.osx.vendored_libraries = 'lib{{.LibName}}.dylib'
end
