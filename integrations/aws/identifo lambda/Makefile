.PHONY:compile

app = Identifo
function = Identifo
# Specify your lambda arn, API Gateway endpoint, S3 bucket with admin panel, and AWS profile.
# aws_arn = 
# api_gateway_endpoint = 
# admin_panel_bucket = 
# aws_profile =

compile:
		GOOS=linux go build -o $(app)

run:	compile
		sam local invoke $(app) -e test/event.json 

ship:	compile
		zip -r $(app)_lambda.zip $(app) server-config.yaml jwt email_templates static import
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

update_admin_panel:
		# Fetch and build source code.
		wget https://github.com/MadAppGang/identifo-admin/archive/master.zip
		tar xvf master.zip
		cd identifo-admin-master && export API_URL=$(api_gateway_endpoint) && \
		npm i && npm run build && \
		rm -rf ../admin_panel/build/* && mv build ../admin_panel

		# Clean up.
		rm -f master.zip
		rm -fr identifo-admin-master

deploy_admin_panel: update_admin_panel
		# Empty bucket.
		aws s3 rm s3://$(admin_panel_bucket) --recursive
		# Upload new contents.
		aws s3 sync admin_panel/build s3://$(admin_panel_bucket) --acl public-read