.PHONY:compile

app = Identifo
function = Identifo
# Specify needed environment variables.
# api_gateway_endpoint = 
# aws_profile = 

compile: update_admin_panel
		dep ensure -update
		GOOS=linux go build -o $(app)

run:	compile
		sam local invoke $(app) -e test/event.json 

ship:	compile
		zip -r $(app)_lambda.zip $(app) server-config.yaml jwt admin_panel email_templates static import
		rm $(app)

deploy:	ship
		aws lambda update-function-code --function-name $(function) --zip-file fileb://$(app)_lambda.zip --publish --profile $(aws_profile)
		# aws lambda update-function-configuration --function-name $(function) --environment 
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
		cd identifo-admin-master && \
		export API_URL=$(api_gateway_endpoint) && export BASE_URL=/develop/adminpanel/ && export ASSETS_PATH=/develop/adminpanel/ && \
		npm i && npm run build && \
		rm -rf ../admin_panel/build/* && mv build ../admin_panel

		# Clean up.
		rm -f master.zip
		rm -fr identifo-admin-master