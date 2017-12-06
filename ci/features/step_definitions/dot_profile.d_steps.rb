Then(/^the environment contains the secret values as per secrets\.yml$/) do
  step "the 'env' command is run"
  expect(@output).to include("CONJUR_SECRET=a conjur secret")
  expect(@output).to include("LITERAL_SECRET=a literal secret")
end

When(/^the \.profile\.d script is sourced$/) do
  @commands ||= []
  @commands << <<EOL
. #{@BUILD_DIR}/.profile.d/0000_retrieve-secrets.sh #{@BUILD_DIR}
EOL
end
