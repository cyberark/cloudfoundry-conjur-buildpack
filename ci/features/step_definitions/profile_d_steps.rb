When(/^the retrieve secrets profile\.d script is sourced$/) do
  @commands ||= []

  # Set $HOME and $DEPS_DIR to mimic variables that are made available to the
  # script when run normally via profile.d.
  @commands << <<EOL
HOME=#{@BUILD_DIR} DEPS_DIR=#{@DEPS_DIR} . #{@DEPS_DIR}/#{@INDEX_DIR}/profile.d/0001_retrieve-secrets.sh
EOL
end

Then(/^the environment contains$/) do |text|
  expect(@output).to include(text)
end
