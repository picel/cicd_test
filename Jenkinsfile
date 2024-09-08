/* groovylint-disable NestedBlockDepth */
pipeline {
    agent any

    environment {
        DOCKER_HUB_CREDENTIALS = credentials('docker-hub-credentials-id')
    }

    stages {
        stage('Prepare Workspace') {
            steps {
                sh 'git clone https://github.com/picel/cicd_test.git'
                withCredentials([usernamePassword(credentialsId: 'docker-hub-credentials-id', usernameVariable: 'DOCKER_HUB_USERNAME', passwordVariable: 'DOCKER_HUB_PASSWORD')]) {
                    sh "docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}"
                }
            }
        }

        stage('Apply K8s YAML Files in /infra/k8s/') {
            steps {
                script {
                  def scriptPaths = sh(
                    script: 'ls cicd_test/infra/k8s/*.yaml',
                    returnStdout: true
                  ).trim().split('\n')

                  scriptPaths.each { scriptPath ->
                    sh "kubectl apply -f ${scriptPath}"
                  }
                }
            }
        }

        stage('Deploy Services') {
            parallel {
                stage('Build and Deploy BFF') {
                    steps {
                        dir('cicd_test/infra/bff') {
                            script {
                                def imageName = 'tkdqja9573/bff-server:latest'
                                sh "docker build -t ${imageName} ."
                                sh "docker push ${imageName}"
                                
                                sh "kubectl apply -f k8s/deployment.yaml"
                                sh "kubectl apply -f k8s/hpa.yaml"
                            }
                        }
                    }
                }

                stage('Build and Deploy Frontend') {
                    steps {
                        dir('cicd_test/frontend') {
                            script {
                                def imageName = 'tkdqja9573/frontend:latest'
                                sh "docker build -t ${imageName} ."
                                sh "docker push ${imageName}"

                                sh "kubectl apply -f k8s/deployment.yaml"
                                sh "kubectl apply -f k8s/hpa.yaml"
                            }
                        }
                    }
                }

                stage('Build and Deploy Backend Services') {
                    steps {
                        script {
                            def backendDirs = sh(
                                script: 'ls -d cicd_test/backend/*/',
                                returnStdout: true
                            ).trim().split('\n')
                            backendDirs.each { directory ->
                                def imageName = "tkdqja9573/${directory.split('/').last()}:latest"
                                sh "echo ${directory}"
                                dir(directory) {
                                    script {
                                        sh "echo ${directory}"
                                        sh "docker build -t ${imageName} ."
                                        sh "docker push ${imageName}"

                                        sh "kubectl apply -f k8s/deployment.yaml"
                                        sh "kubectl apply -f k8s/hpa.yaml"
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
    }
}
