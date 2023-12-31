pipeline {
    agent any
    environment {
        GIT_URL = 'https://github.com/ajiinisti/livecode-cicd.git'
        BRANCH = 'main'
        IMAGE = 'livecode-cicd'
        CONTAINER = 'livecode-cicd-backend'
        DOCKER_APP = 'docker'
        DB_HOST = 'livecode-cicd-db'
        DB_USER = 'postgres'
        DB_NAME = 'book_management_system'
        DB_PASSWORD = 'password'
        DB_PORT = '5432'
        API_PORT = '8000'
        SLACK_WEBHOOK_URL = 'https://hooks.slack.com/services/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'
        EMAIL_RECIPIENT = 'ajiinisti@gmail.com'
    }
    stages {
        stage("Cleaning up") {
            steps {
                echo 'Cleaning up'
                sh "${DOCKER_APP} rmi -f ${IMAGE} || true"
                sh "${DOCKER_APP} rm -f ${CONTAINER} || true"
                sh "${DOCKER_APP} rm -f ${DB_HOST} || true"
            }
        }

        stage("Clone") {
            steps {
                echo 'Clone'
                git branch: "${BRANCH}", url: "${GIT_URL}"
            }
        }

        stage("Build and Run") {
            steps {
                echo 'Build and Run'
                sh "DB_HOST=${DB_HOST} DB_PORT=${DB_PORT} DB_NAME=${DB_NAME} DB_USER=${DB_USER} DB_PASSWORD=${DB_PASSWORD} API_PORT=${API_PORT} ${DOCKER_APP} compose up -d"
                sleep 8
                script {
                    def dbExists = sh(script: "docker exec ${DB_HOST} psql -U ${DB_USER} -d ${DB_NAME} -c 'SELECT 1' 2>/dev/null", returnStatus: true)
                    if (dbExists == 0) {
                        echo "Database '${DB_NAME}' already exists"
                    } else {
                        sh "docker exec ${DB_HOST} createdb -U ${DB_USER} ${DB_NAME}"
                        echo "Database '${DB_NAME}' created"
                    }
                }
            }
        }
    }
    post {
        always {
            echo 'This will always run'
            postNotify(currentBuild.result)
        }
        success {
            echo 'This will run only if successful'
            postNotify(currentBuild.result)
        }
        failure {
            echo 'This will run only if failed'
            postNotify(currentBuild.result)
        }
    }
}

def postNotify(buildStatus) {
    notifyEmail(buildStatus)
    notifySlack(buildStatus)
}

def notifyEmail(buildStatus) {
    def subject
    def body

    if (buildStatus == 'SUCCESS') {
        subject = "Build Successful"
        body = "The build has succeeded."
    } else {
        subject = "Build Failed"
        body = "The build has failed."
    }

    emailext body: body, subject: subject, to: "${EMAIL_RECIPIENT}"
}

def notifySlack(buildStatus) {
    def color
    def statusText

    if (buildStatus == 'SUCCESS') {
        color = 'good'
        statusText = 'succeeded'
    } else {
        color = 'danger'
        statusText = 'failed'
    }

    def payload = [
        'attachments': [[
            'color': color,
            'title': "Build ${statusText}",
            'text': "The build has ${statusText}.",
            'footer': "Jenkins"
        ]]
    ]

    sh """
    curl -X POST \
        -H 'Content-type: application/json' \
        --data-urlencode "${payload}" \
        ${SLACK_WEBHOOK_URL}
    """
}