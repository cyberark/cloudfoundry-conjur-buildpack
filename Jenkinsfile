#!/usr/bin/env groovy
@Library(['product-pipelines-shared-library', 'conjur-enterprise-sharedlib']) _

// Automated release, promotion and dependencies
properties([
  // Include the automated release parameters for the build
  release.addParams(),
  // Dependencies of the project that should trigger builds
  dependencies([])
])

// Performs release promotion.  No other stages will be run
if (params.MODE == "PROMOTE") {
  release.promote(params.VERSION_TO_PROMOTE) { infrapool, sourceVersion, targetVersion, assetDirectory ->
    // Any assets from sourceVersion Github release are available in assetDirectory
    // Any version number updates from sourceVersion to targetVersion occur here
    // Any publishing of targetVersion artifacts occur here
    // Anything added to assetDirectory will be attached to the Github Release

    //Note: assetDirectory is on the infrapool agent, not the local Jenkins agent.
  }
  release.copyEnterpriseRelease(params.VERSION_TO_PROMOTE)
  return
}

pipeline {
  agent { label 'conjur-enterprise-common-agent' }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
    lock resource: "tas-infra"
  }

  environment {
    // Sets the MODE to the specified or autocalculated value as appropriate
    MODE = release.canonicalizeMode()
  }

  triggers {
    cron(getDailyCronString())
    parameterizedCron(getWeeklyCronString("H(1-5)","%MODE=RELEASE"))
  }

  stages {
    // Aborts any builds triggered by another project that wouldn't include any changes
    stage ("Skip build if triggering job didn't create a release") {
      when {
        expression {
          MODE == "SKIP"
        }
      }
      steps {
        script {
          currentBuild.result = 'ABORTED'
          error("Aborting build because this build was triggered from upstream, but no release was built")
        }
      }
    }

    stage('Get InfraPool ExecutorV2 Agent(s)') {
      steps {
        script {
          // Request ExecutorV2 agents for 1 hour
          infrapool = getInfraPoolAgent.connected(type: "ExecutorV2", quantity: 1, duration: 1)[0]
        }
      }
    }

    // Generates a VERSION file based on the current build number and latest version in CHANGELOG.md
    stage('Validate Changelog and set version') {
      steps {
        script {
          updateVersion(infrapool, "CHANGELOG.md", "${BUILD_NUMBER}")
        }
      }
    }

    stage('Get latest upstream dependencies') {
      steps {
        script {
          updatePrivateGoDependencies("${WORKSPACE}/conjur-env/go.mod")
          infrapool.agentPut from: "conjur-env/vendor", to: "${WORKSPACE}/conjur-env/"
          infrapool.agentPut from: "conjur-env/go.*", to: "${WORKSPACE}/conjur-env/"
          infrapool.agentPut from: "/root/go", to: "/var/lib/jenkins/"
        }
      }
    }

    stage('Package') {
      steps {
        script {
          infrapool.agentSh './package.sh --skip-gomod-download && ./unpack.sh'
          infrapool.agentArchiveArtifacts(artifacts: 'conjur_buildpack*.zip')
        }
      }
    }

    stage('Test') {
      parallel {
        stage('Integration Tests') {
          steps {
            script {
              env.INFRAPOOL_JENKINS_HOME = env.JENKINS_HOME
              infrapool.agentSh './ci/test_integration'
              sh 'mkdir -p tests/integration/reports/integration'
              infrapool.agentGet(
                from: 'tests/integration/reports/integration/*.xml',
                to: 'tests/integration/reports/integration/'
              )
              junit 'tests/integration/reports/integration/*.xml'
            }
          }
        }

        // Disabled as allocateTas isn't currently available.
        // stage('End To End Tests') {
        //   steps {
        //     script {
        //       env.INFRAPOOL_JENKINS_HOME = env.JENKINS_HOME
        //       allocateTas.withTas(infrapool, 'isv_ci_tas_srt_5_0'){
        //         infrapool.agentSh 'summon -f ./ci/secrets.yml ./ci/test_e2e'
        //         sh 'mkdir -p tests/integration/reports/e2e'
        //         infrapool.agentGet(
        //           from: 'tests/integration/reports/e2e/*.xml',
        //           to: 'tests/integration/reports/e2e/'
        //         )
        //         junit 'tests/integration/reports/e2e/*.xml'
        //       }
        //     }
        //   }
        // }

        stage('Unit Tests') {
          stages {
            stage("Secret Retrieval Script Tests") {
              steps {
                script {
                  infrapool.agentSh './tests/retrieve-secrets/start'
                  infrapool.agentGet(
                    from: 'TestReport-test.xml',
                    to: 'TestReport-test.xml'
                  )
                  junit 'TestReport-test.xml'
                }
              }
            }

            stage("Conjur-Env Unit Tests") {
              steps {
                script {
                  infrapool.agentSh './ci/test_conjur-env'
                  infrapool.agentGet(
                    from: 'conjur-env/output/*.xml',
                    to: 'conjur-env/output/'
                  )
                  infrapool.agentGet(
                    from: 'conjur-env/output/c.out',
                    to: 'conjur-env/output/'
                  )
                }
              }

              post {
                always {
                  junit 'conjur-env/output/*.xml'
                  cobertura autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: 'conjur-env/output/coverage.xml', conditionalCoverageTargets: '30, 0, 0', failUnhealthy: false, failUnstable: false, lineCoverageTargets: '30, 0, 0', maxNumberOfBuilds: 0, methodCoverageTargets: '30, 0, 0', onlyStable: false, sourceEncoding: 'ASCII', zoomCoverageChart: false

                  // Don't fail builds if we can't upload coverage information to Codacy
                  catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    retry (2) {
                      codacy action: 'reportCoverage', language: 'Go', filePath: 'conjur-env/output/c.out', extraArgs: '--force-coverage-parser go'
                    }
                  }
                }
              }
            }
          }
        }
      }
    }

    stage('Release') {
      when {
        expression {
          MODE == "RELEASE"
        }
      }

      steps {
        script {
          release(infrapool, { billOfMaterialsDirectory, assetDirectory ->
            /* Publish release artifacts to all the appropriate locations
               Copy any artifacts to assetDirectory on the infrapool node
               to attach them to the Github release.

               If your assets are on the infrapool node in the target
               directory, use a copy like this:
                  infrapool.agentSh "cp target/* ${assetDirectory}"
               Note That this will fail if there are no assets, add :||
               if you want the release to succeed with no assets.

               If your assets are in target on the main Jenkins agent, use:
                 infrapool.agentPut(from: 'target/', to: assetDirectory)
            */
            infrapool.agentSh "cp conjur_buildpack*.zip ${assetDirectory}"
          })
        }
      }
    }
  }
  post {
    always {
      releaseInfraPoolAgent()
    }
  }
}
