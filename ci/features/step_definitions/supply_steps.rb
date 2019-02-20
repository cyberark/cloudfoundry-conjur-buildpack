Then(/^conjur-env is installed$/) do
  `ls #{@VENDOR_DIR}/conjur-env`
  expect($?.exitstatus).to eq (0)
end

Then(/^the retrieve secrets profile\.d script is installed$/) do
  expect(File.exist?("#{@DEPS_DIR}/#{@INDEX_DIR}/profile.d/0001_retrieve-secrets.sh")).to be_truthy
end
