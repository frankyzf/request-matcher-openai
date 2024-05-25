pipeline {
    agent {
        label 'gobin'
    }
    tools {
        go 'go'
    }
    environment {
        GO111MODULE = 'on'
        PATH="/var/data/jenkins/node/bin/:${env.PATH}:/usr/local/go/bin/go:/home/homeplus/workspace/gopath/bin/:${env.HOME}/gopath/bin"
        PKG_CONFIG_PATH="${env.HOME}/anaconda3/envs/py37/lib/pkgconfig"
        GOPATH="${env.HOME}/gopath"
        GOROOT="${env.HOME}/gopath/go"
        DOCKER_SERVER="docker.autom-tech.com:5500"
        DOCKER_USER="homeplus"
        DOCKER_PASSWORD="admin123!"

        NOTIFY_RECEIVERS="dev@coloso.io,27504490@qq.com,etwo@coloso.io,abby@coloso.io,276644176@qq.com"
        NAME="request-matcher-openai"
    }

    stages {
        stage('mod init') {
            steps {
                script {
                    sh "cp ~/go.work ."
                }
            }
        }

        stage('go build') {
            steps {
                sh "make tidy"
                sh "make build"
                sh "git checkout ."
                // sh "docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD}  ${DOCKER_SERVER}"
                // sh "docker build -t $NAME:latest ."
                // sh "docker tag $NAME:latest ${DOCKER_SERVER}/$NAME:latest"
            }
        }
        // stage('Push') {
        //     steps {
        //         sh "docker push ${DOCKER_SERVER}/$NAME:latest"
        //     }
        // }
        stage('deploy') {
            steps {
                script {
                        sh "echo branch: ${GIT_BRANCH}"
                        if (env.GIT_BRANCH.startsWith("dev") || env.GIT_BRANCH.startsWith("v")) {
                            sh "make uploaduat;"
                            sh 'ssh  work4  "source ~/.bashrc; cd ~/service/request-matcher-openai-prod; ./upgrade.sh; sleep 5; ./start.sh&" '

                        } else {
                            //do nothing
                        }
                }
            }
        }
    }
    post {
        success {
            archiveArtifacts artifacts: 'bin/**,config/**,upgrade_bin.sh,upgrade.sh', onlyIfSuccessful: true
            sh 'bash /home/homeplus/service/helper/compile_notify.sh "successfully compile $NAME[$GIT_BRANCH]" '
            script {
                if (env.GIT_BRANCH.startsWith("dev") || env.GIT_BRANCH.startsWith("v")) {
                    sh "touch test_summary.log; rm test_summary.log; touch test.log; rm test.log; echo 'auto deploy develop to uat' "
                }
            }
        }
        unsuccessful {
            // sh 'touch test_summary.log; touch test.log; python3 /home/homeplus/service/helper/sendEmail.py "test.log" "$NOTIFY_RECEIVERS" "failed to compile $NAME[$GIT_BRANCH]" "result: $(cat test_summary.log)"  '
            sh 'bash /home/homeplus/service/helper/compile_notify.sh "failed to test/compile $NAME[$GIT_BRANCH]" '
        }
    }
}
