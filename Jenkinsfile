pipeline {
   agent {
     node {
       label 'jenkins-agent'
     }
   }


  stages {
    stage('build & push') {
      steps {
        container ('jdk') {
          // sh 'git clone https://github.com/yuswift/devops-go-sample.git'
          sh 'docker build --network host -t cicd-template .'
        }
      }
    }



  }
}
