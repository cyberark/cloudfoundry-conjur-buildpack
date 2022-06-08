#!/usr/bin/env groovy

pipeline {
  agent { label 'executor-v2' }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
    lock resource: "tas-infra"
  }

  triggers {
    cron(getDailyCronString())
  }

  stages {
    stage('Grant IP Access') {
      steps {
        // Grant access to this Jenkins agent's IP to AWS security groups
        grantIPAccess()
      }
    }

    stage('Validate Changelog') {
      steps {
        sh './ci/parse-changelog.sh'
      }
    }

    stage('Package') {
      steps {
        sh './package.sh && ./unpack.sh'
        archiveArtifacts('conjur_buildpack*.zip')
      }
    }

    stage('Test') {
      parallel {
        stage('Integration Tests') {
          steps {
            sh './ci/test_integration'
            junit 'tests/integration/reports/integration/*.xml'
          }
        }

        stage('End To End Tests') {
          steps {
            allocateTas()
            sh 'summon -f ./ci/secrets.yml ./ci/test_e2e'
            junit 'tests/integration/reports/e2e/*.xml'
          }
          post {
            always {
              destroyTas()
            }
          }
        }

        stage('Unit Tests') {
          stages {
            stage("Secret Retrieval Script Tests") {
              steps {
                sh './tests/retrieve-secrets/start'
                junit 'TestReport-test.xml'
              }
            }

            stage("Conjur-Env Unit Tests") {
              steps {
                sh './ci/test_conjur-env'
                junit 'conjur-env/output/*.xml'
              }
            }
          }
        }
      }
    }
  }

  post {
    always {
      cleanupAndNotify(currentBuild.currentResult)
    }
  }
}
