.PHONY:compile

app = Identifo
function = Identifo
# Specify your lambda arn and AWS profile.
# aws_arn = 
# aws_profile =

compile:
		GOOS=linux go build -o $(app)

run:	compile
		sam local invoke $(app) -e test/event.json 

ship:	compile
		zip $(app)_lambda.zip $(app) 
		rm $(app)

deploy:	ship
		aws lambda update-function-code --function-name $(function) --zip-file fileb://$(app)_lambda.zip --publish --profile $(aws_profile)
		rm $(app)_lambda.zip

update:	
		dep ensure -update

test:
		go test

debug:
		sam local start-api