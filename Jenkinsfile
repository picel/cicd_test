/* groovylint-disable NestedBlockDepth */
pipeline {
    agent any

    tools {
      dockerTool 'docker'
    }

    environment {
        DOCKER_HUB_CREDENTIALS = credentials('docker-hub-credentials-id')
    }

    stages {
        stage('Prepare Workspace') {
            steps {
                git 'https://github.com/picel/cicd_test.git'
                withCredentials([usernamePassword(credentialsId: 'docker-hub-credentials-id', usernameVariable: 'DOCKER_HUB_USERNAME', passwordVariable: 'DOCKER_HUB_PASSWORD')]) {
                    sh "docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}"
                }
            }
        }

        stage('Deploy Infrastructure and Services') {
            parallel {
                stage('Apply K8s YAML Files in /infra/k8s/') {
                    steps {
                        script {
                            def yamlFiles = findFiles(glob: 'infra/k8s/**/*.yaml')
                            yamlFiles.each { file ->
                                sh "kubectl apply -f ${file.path}"
                            }
                        }
                    }
                }

                stage('Build and Deploy BFF') {
                    steps {
                        dir('infra/bff') {
                            script {
                                def imageName = 'tkdqja9573/bff-server:latest'
                                sh "docker build -t ${imageName} ."
                                sh "docker push ${imageName}"

                                sh "kubectl apply -f k8s/deployment.yaml"
                            }
                        }
                    }
                }

                stage('Build and Deploy Frontend') {
                    steps {
                        dir('frontend') {
                            script {
                                def imageName = 'tkdqja9573/frontend:latest'
                                sh "docker build -t ${imageName} ."
                                sh "docker push ${imageName}"

                                sh "kubectl apply -f k8s/deployment.yaml"
                            }
                        }
                    }
                }

                stage('Build and Deploy Backend Services') {
                    steps {
                        script {
                            def backendDirs = findFiles(glob: 'backend/*/Dockerfile').collect { it.directory }.unique()
                            backendDirs.each { dir ->
                                dir = dir.replace('Dockerfile', '')
                                dir(dir) {
                                    script {
                                        def imageName = "tkdqja9573/${dir.split('/').last()}:latest"
                                        sh "docker build -t ${imageName} ."
                                        sh "docker push ${imageName}"

                                        sh "kubectl apply -f k8s/deployment.yaml"
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }

        stage('Clean Workspace') {
            steps {
                cleanWs()
            }
        }
    }

    post {
        always {
            cleanWs()
        }
    }
}
