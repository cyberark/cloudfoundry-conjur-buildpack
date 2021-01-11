Then(/^the '([^"]*)' variable has a secret value$/) do |id, value|
  conjur_api.resource(make_full_id('variable', id)).add_value value
end
