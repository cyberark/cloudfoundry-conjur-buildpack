#!/bin/bash -e

export COMPOSE_FILE=ci/docker-compose.yml

docker-compose up -d conjur

# wait for conjur
docker-compose run --rm tester bash <<'EOF'
printf "Waiting for Conjur to be ready... "
while [[ ! $(curl -o /dev/null -s -w '%{http_code}\n' conjur) == 200 ]]
do
  echo "."
  sleep 1
done
echo "Done."
EOF

# setup conjur
export API_KEY=$(docker-compose exec -T conjur rails r "print Credentials['cucumber:user:admin'].api_key")

docker-compose run --rm -v $(pwd)/ci/policy:/policy \
 -e CONJUR_APPLIANCE_URL=http://conjur \
 -e CONJUR_ACCOUNT=cucumber \
 -e CONJUR_AUTHN_API_KEY=$API_KEY \
 -e CONJUR_AUTHN_LOGIN=admin \
 --entrypoint bash tester -c "
conjur policy load root /policy/root.yml
conjur variable values add conjur_secret_id 'a conjur secret'
"

CONJUR_CREDENTIALS_JSON=$(cat << EOL
{
  "appliance_url": "http://conjur",
  "authn_api_key": "$API_KEY",
  "authn_login": "admin",
  "account": "cucumber"
}
EOL
)
