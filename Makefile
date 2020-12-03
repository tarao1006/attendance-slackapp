HEROKU:=`which heroku`
APP_FLAG=-a attendance-slackapp

logs:
	${HEROKU} ${@F} ${APP_FLAG} --tail
