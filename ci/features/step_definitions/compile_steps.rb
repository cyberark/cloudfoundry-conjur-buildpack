Then(/^summon is installed$/) do
  `#{@BUILD_DIR}/summon -v`
  expect($?.exitstatus).to eq (0)
end

Then(/^summon-conjur is installed$/) do
  `#{@BUILD_DIR}/summon-conjur -v`
  expect($?.exitstatus).to eq (0)
end

Then(/^the retrieve secrets \.profile\.d script is installed$/) do
  expect(File.exist?("#{@BUILD_DIR}/.profile.d/0000_retrieve-secrets.sh")).to be_truthy
end
