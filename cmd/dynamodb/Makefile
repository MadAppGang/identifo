

list_tables:
	aws dynamodb list-tables --endpoint-url http://localhost:8000 --region=us-west-2

apps:
	aws dynamodb scan --table-name Applications --endpoint-url http://localhost:8000 --region=us-west-2


users:
	aws dynamodb scan --table-name Users --endpoint-url http://localhost:8000 --region=us-west-2

fusers:
	aws dynamodb scan --table-name UsersByFederatedIDpo --endpoint-url http://localhost:8000 --region=us-west-2
