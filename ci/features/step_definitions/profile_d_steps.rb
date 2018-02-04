When(/^the retrieve secrets \.profile\.d script is sourced$/) do
  @commands ||= []
  @commands << <<EOL
. #{@BUILD_DIR}/.profile.d/0001_retrieve-secrets.sh #{@BUILD_DIR}
EOL
end

Then(/^the environment contains$/) do |text|
  expect(@output).to include(text)
end

When(/^the \.profile\.d scripts are sourced$/) do
  @commands ||= []
  @commands << <<EOL
. #{@BUILD_DIR}/.profile.d/0001_retrieve-secrets.sh #{@BUILD_DIR}
EOL
end
