HEROKU:=`which heroku`
APP_NAME="attendance-slackapp"

logs:
	${HEROKU} logs -a ${APP_NAME} --tail
