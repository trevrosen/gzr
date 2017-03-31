node('default') {
  stage('Test CLI') {
    deleteDir()
    checkout scm
    // won't need this key after repo is public
    withCredentials([file(credentialsId: 'bypass-cicd-sshkey', variable: 'SSHKEY')]) {
      sh 'sudo cp -u ${SSHKEY} /var/jenkins_home/.ssh/id_rsa'
      sh 'sudo cp -u ${SSHKEY} gozer-web/id_rsa'
      sh 'sudo chmod go-rwx /var/jenkins_home/.ssh/id_rsa && sudo chown jenkins: /var/jenkins_home/.ssh/id_rsa'
      sh 'sudo chmod go-rwx gozer-web/id_rsa && sudo chown jenkins: gozer-web/id_rsa'
    }
    sh './ci/test-cli.sh'
  }
  stage('Test Web') {
    sh './ci/test-web.sh'
  }
  stage('Build CLI') {
    sh './ci/build.sh'
  }
  stage('Release CLI to Github') {
    withCredentials([[$class: 'StringBinding', credentialsId: 'c2c752a0-cbc7-44ed-9ba7-5467bc9cd2ec', variable: 'GITHUB_TOKEN']]) {
      sh './ci/release-github.sh'
    }
  }
  stage('Release Web to Dockerhub') {
    withCredentials([file(credentialsId: 'docker-hub-registry-config', variable: 'DOCKERCFG')]) {
      sh 'cp -u ${DOCKERCFG} /var/jenkins_home/.dockercfg'
      sh 'chmod 600 /var/jenkins_home/.dockercfg'
    }
    sh './ci/release-docker.sh'
  }
}
