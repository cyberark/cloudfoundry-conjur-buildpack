require 'rspec'
require 'tmpdir'

RSpec.configure do |config|
  config.expect_with :rspec do |c|
    c.syntax = :expect
  end
end

Before('@BUILD_DIR') do
  @BUILD_DIR = Dir.mktmpdir
end

After('@BUILD_DIR') do
  FileUtils.remove_entry_secure @BUILD_DIR
end
