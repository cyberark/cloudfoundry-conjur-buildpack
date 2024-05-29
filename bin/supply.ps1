
$buildDir=$args[0]
$depsDir=$args[2]
$indexDir=$args[3]

# Validate that secret.yml exists
if (![System.IO.File]::Exists("$buildDir\secrets.yml"))
{
    echo "Unable to find a secrets.yml...exiting"
    exit 1
}

# Validate that VCAP_SERVICES contains 'cyberark-conjur'
$vcapJson = echo $env:VCAP_SERVICES | ConvertFrom-Json

if ("true" -ne $Env:CONJUR_BUILDPACK_BYPASS_SERVICE_CHECK )
{
    if( !("cyberark-conjur" -in $vcapJson.PSobject.Properties.Name) )
    {
        echo "No credentials for cyberark-conjur service found in VCAP_SERVICES... exit"
        exit 1
    }
}

pushd $depsDir\$indexDir
  mkdir profile.d | Out-Null
  copy $PSScriptRoot\..\lib\0001_retrieve-secrets.bat .\profile.d\
popd

pushd $buildDir
  mkdir .conjur | Out-Null
  copy $PSScriptRoot\..\vendor\conjur-win-env.exe .\.conjur\
popd