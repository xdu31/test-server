@Library('jenkins.shared.library') _

pipeline {
    agent {
        label 'ubuntu_docker_label'
    }

    tools {
        go "Go 1.14.4"
    }

    environment {
        PROJECT                = "src/github.com/xdu31/test-server"
        GOPATH                 = "$WORKSPACE"
        PATH                   = "$PATH:$GOPATH/bin"
        image_repo             = "xdu31/test-server"
        upstream_github_repo   = "xdu31/test-server"
        upstream_execution_url = "$BUILD_URL"
        upstream_job_name      = "$JOB_BASE_NAME"
        upstream_build_id      = "$BUILD_ID"
        email_recipients       = "eaglerose31@gmail.com"
        GIT_BRANCH             = rewrite_image_tag(env.GIT_BRANCH)
    }

    options {
        checkoutToSubdirectory("src/github.com/xdu31/test-server")
    }

    stages {
        stage('Build image') {
            steps {
                sh "cd $PROJECT && make GO_CACHE= build && make image"
            }
        }
        stage("Check pull request") {
            when {
                expression {
                    return env.BRANCH_NAME.startsWith('PR')
                }
            }

            steps {
                sh """
                sudo curl -L https://github.com/docker/compose/releases/download/1.26.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
                sudo chmod +x /usr/local/bin/docker-compose
                docker-compose --version
                cd $PROJECT && make itest
                """
            }
        }
        stage("Push chart") {
            when {
                expression { !isPrBuild() }
            }
            steps {
                withAWS(region:'us-east-1', credentials:'CICD_HELM') {
                    sh "cd $PROJECT && make push-chart"
                }

                dir("${WORKSPACE}/${PROJECT}") {
                    archiveArtifacts artifacts: 'helm/*.tgz'
                    archiveArtifacts artifacts: 'helm/build.properties'
                }
            }
        }
         stage("Cleaning") {
            steps {
                 sh "cd $PROJECT && make clean"
            }
        }
    }

    post {
        success {
            echo 'calling CVE scanning job now . . .'
            build job: 'run_docker_image_cve_scan',
            parameters: [
                string(name: 'repo', value: env.image_repo),
                string(name: 'upstream_github_repo', value: env.upstream_github_repo),
                string(name: 'upstream_execution_url', value: env.upstream_execution_url),
                string(name: 'upstream_job_name', value: env.upstream_job_name),
                string(name: 'upstream_build_id', value: env.upstream_build_id),
                string(name: 'email_recipients', value: env.email_recipients)
            ],

            wait: false
        }
    }
}

def rewrite_image_tag(tag){
    tag=tag.replace("/", "-")
    return tag
}
