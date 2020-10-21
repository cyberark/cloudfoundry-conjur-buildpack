#!/usr/bin/env groovy

pipeline {
  agent { label 'executor-v2' }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
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
    stage('Validate') {
      parallel {
        stage('Changelog') {
          steps { sh './bin/parse-changelog.sh' }
        }
      }
    }
    
    stage('Build buildpack') {
      steps {
        sh './package.sh'

        archiveArtifacts artifacts: '*.zip', fingerprint: true
      }
    }
    stage('Test') {
      steps {
        sh 'summon ./test.sh'

        junit 'ci/features/reports/*.xml'
      }
    }

  }

  post {
    always {
      cleanupAndNotify(currentBuild.currentResult)
      // Remove this Jenkins Agent's IP from AWS security groups
      removeIPAccess()
    }
  }
}
