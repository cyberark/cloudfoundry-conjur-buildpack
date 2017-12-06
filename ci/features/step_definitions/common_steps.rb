require 'tempfile'

Given(/^the '([^"]*)' command is run$/) do |script_command|
  @commands ||= []

  # all this fuss to use a bash shell as opposed to sh shell

  f = Tempfile.open('command.sh')
  file_contents = <<EOS
#!/bin/bash -e
#{ @commands.join("\n") }
#{ script_command }
EOS
  f.write(file_contents)
  f.close

  `chmod +x #{f.path}`

  @output = `bash #{f.path}`
  @result = $?

  f.unlink
end

Given(/^the '([^"]*)' script is run$/) do |script|
  step "the '/buildpack-build/bin/#{script} #{@BUILD_DIR}' command is run"
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

Given(/^VCAP_SERVICES contains cybark\-conjur credentials$/) do
  @commands ||= []
  @commands << <<eos
export VCAP_SERVICES='
{
 "cyberark-conjur": [{
   "credentials": #{ENV['CONJUR_CREDENTIALS_JSON']}
 }]
}
'
eos
end

Given(/^VCAP_SERVICES has a cyberark\-conjur key$/) do
  @commands ||= []
  @commands << <<EOS
export VCAP_SERVICES='
{
 "cyberark-conjur": []
}
'
EOS
end

Given(/^VCAP_SERVICES does not have a cyberark\-conjur key$/) do
  @commands ||= []
  @commands << <<EOS
export VCAP_SERVICES='
{
}
'
EOS
end

And(/^the build directory has a secrets\.yml file$/) do
  secretsyml = <<EOS
CONJUR_SECRET: !var conjur_secret_id
LITERAL_SECRET: a literal secret
EOS
  File.open("#{@BUILD_DIR}/secrets.yml", 'w') { |file| file.write(secretsyml) }
end
