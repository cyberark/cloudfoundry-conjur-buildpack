Then(/^conjur-env is installed$/) do
  `ls #{@BIN_DIR}/conjur-env`
  expect($?.exitstatus).to eq (0)
end

Then(/^the retrieve secrets \.profile\.d script is installed$/) do
  expect(File.exist?("#{@BUILD_DIR}/.profile.d/0001_retrieve-secrets.sh")).to be_truthy
end
