Given('I create an org and space') do
  login_to_pcf
  cf_ci_org
  cf_ci_space
end

Given('I install the buildpack') do
  cf_auth('admin', ENV['CF_ADMIN_PASSWORD'])

  Dir.chdir('..') do
    ShellSession.execute('../upload.sh', 'BUILDPACK_NAME' => cf_ci_buildpack_name)
  end
end

When('I push a {string} app with the {string} buildpack') do |language, buildpack_type|
  login_to_pcf
  cf_target(cf_ci_org, cf_ci_space)

  @app_name = "#{buildpack_type}-#{language}-app"

  Dir.chdir("apps/#{language}") do
    if buildpack_type == 'online'
      create_online_app_manifest
    else
      create_offline_app_manifest
    end
    ShellSession.execute("cf push #{@app_name} --random-route")
  end
end

Then('the secrets.yml values are available in the app') do
  page_content = cf_app_content
  expect(page_content).to match(/Database Username: space_username/)
  expect(page_content).to match(/Database Password: space_password/)
end
