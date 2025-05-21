pipeline {
  agent any
  environment {
    DOCKER_IMAGE = "rpc-demo:${env.BUILD_ID}"
  }
  stages {
    stage('拉取代码') {
      steps {
        checkout([$class: 'GitSCM',
                 branches: [[name: '*/main']],
                 userRemoteConfigs: [[url: 'https://gitee.com/starry1213/useDemo.git', credentialsId: 'gitee-credentials']]]
        )
      }
    }
    stage('构建Docker镜像') {
      steps {
        script {
          // 直接使用 Dockerfile 中的编译步骤
          docker.build("${env.DOCKER_IMAGE}", "-f rpc-demo.dockerfile .")
        }
      }
    }
    stage('部署容器') {
      steps {
        script {
          sh "docker stop rpc-demo || true"
          sh "docker rm rpc-demo || true"
          sh "docker run -d --name rpc-demo -p 9080:9080 ${env.DOCKER_IMAGE}"
        }
      }
    }
  }
}