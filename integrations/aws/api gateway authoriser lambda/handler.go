package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

//authHandler is a lambda handler to validate JWT token
//parse token
//check for signature
//check claims
//parse claims
//return claims as JSON struct

// https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html
// https://github.com/aws/aws-lambda-go/blob/master/events/README_ApiGatewayCustomAuthorizer.md
//Handler main entry point for Lambda
func Handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	fmt.Printf("Processing incoming event: %v", event)
	alg := model.TokenServiceAlgorithmRS256

	token := event.Headers["Authorization"]
	if len(token) < 50 {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid token in the header incoming data")
	}

	appID := event.Headers["X-Identifo-ClientID"]
	if len(appID) == 0 {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid application ID")
	}

	tstr := jwt.ExtractTokenFromBearerHeader(token)
	if tstr == nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid token in header")
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if len(jwtIssuer) == 0 {
		jwtIssuer = "cc.creatorconnect.link"
	}

	// publicKey, err := jwt.LoadPublicKeyFromString(os.Getenv("PUBLIC_KEY"))
	// if err != nil {
	// 	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid public key: " + err.Error())
	// }
	//if you want to use file instead of env variable, use this code insted
	publicKey, err := jwt.LoadPublicKeyFromPEM("./public.pem", alg)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid public key: " + err.Error())
	}

	v := jwt.NewValidator(appID, jwtIssuer, "", "access")
	tokenV, err := jwt.ParseTokenWithPublicKey(string(tstr), publicKey)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error parsing token: " + err.Error())
	}
	if err := v.Validate(tokenV); err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid token:" + err.Error())
	}

	return CreatePolicy(tokenV.UserID(), "Allow", event.MethodArn, nil), nil
}

//CreatePolicy is a help function to generate an IAM policy
func CreatePolicy(principalID, effect, resource string, context map[string]interface{}) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	// Optional output with custom properties of the String, Number or Boolean type.
	if context != nil {
		authResponse.Context = context
	}
	return authResponse
}
