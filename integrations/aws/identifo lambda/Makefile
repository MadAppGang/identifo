.PHONY:compile
include .env

app = Identifo
function = Identifo
# Check if needed environment variables are defined.
# The precedence of reading environment variables: Makefile -> .env -> environment.
# api_endpoint = 
# aws_profile = 

compile: update_admin_panel
		dep ensure -update
		GOOS=linux go build -o $(app)

run:	compile
		sam local invoke $(app) -e test/event.json 

ship:	compile
		zip -r $(app)_lambda.zip $(app) server-config.yaml jwt static import
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
		$(call check_defined, api_endpoint, URL address of the Identifo server)
		$(call check_defined, aws_profile, Name of the AWS profile)
		# Fetch and build source code.
		wget https://github.com/MadAppGang/identifo-admin/archive/master.zip
		tar xvf master.zip
		cd identifo-admin-master && \
		export API_URL=$(api_endpoint) && export BASE_URL=/adminpanel/ && export ASSETS_PATH=/adminpanel/ && \
		npm i && npm run build && \
		rm -rf ../static/admin_panel/build/* && mv build ../static/admin_panel

		# Clean up.
		rm -f master.zip
		rm -fr identifo-admin-master

check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))