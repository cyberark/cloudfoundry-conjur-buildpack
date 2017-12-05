Given(/^the "([^"]*)" script is run$/) do |script|
  @commands ||= []
  output = `
#{ @commands.join("\n") }
../bin/#{script} #{@BUILD_DIR}
  `
  @result = $?
  puts output
end

Then(/^the result should have a non\-zero exit status$/) do
  expect(@result.exitstatus).not_to eq(0)
end

Then(/^the result should have a 0 exit status$/) do
  expect(@result.exitstatus).to eq(0)
end

Then(/^the result should have a 1 exit status$/) do
  expect(@result.exitstatus).to eq(1)
end

Given(/^VCAP_SERVICES has a cyberark\-conjur key$/) do
  @commands ||= []
  @commands << <<eos
export VCAP_SERVICES='
{
  "cyberark-conjur": []
}
'
eos
end

Given(/^VCAP_SERVICES does not have a cyberark\-conjur key$/) do
  @commands ||= []
  @commands << <<eos
export VCAP_SERVICES='
{
}
'
eos
end

And(/^the build directory has a secrets\.yml file$/) do
  `touch #{@BUILD_DIR}/secrets.yml`
end

Then(/^summon is installed$/) do
  puts `ls -la #{@BUILD_DIR}/summon-conjur`
  `#{@BUILD_DIR}/summon -v`
  expect($?.exitstatus).to eq (0)
end

Then(/^summon-conjur is installed$/) do
  expect(File.exist?("#{@BUILD_DIR}/summon-conjur")).to be_truthy
  `#{@BUILD_DIR}/summon-conjur`
  # expect($?.exitstatus).to eq (0)
end

Then(/^the retrieve secrets \.profile\.d script is installed$/) do
  expect(File.exist?("#{@BUILD_DIR}/.profile.d/0000_retrieve-secrets.sh")).to be_truthy
end
