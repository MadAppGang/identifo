.PHONY:compile

#this is examplemakefile, feel free to change values regarding your aws3 settings
app = IdentifoApiAuth
function = IdentifoApiAuth
aws_arn = arn:aws:lambda:us-east-1:832107067482:function:IdentifoApiAuth
aws_profile = cc

compile:
		GOOS=linux go build -o $(app)

run:	compile
		sam local invoke $(app) -e test/event.json 

ship:	compile
		cp ../../pem/public.pem ./
		zip $(app)_lambda.zip $(app) public.pem
		rm ./public.pem
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