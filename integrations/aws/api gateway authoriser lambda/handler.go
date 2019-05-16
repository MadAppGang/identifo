package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/madappgang/identifo/jwt"
)

// Handler is a main entry point for Lambda. Handler validates JWT tokens and returns IAM roles.
// https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html
// https://github.com/aws/aws-lambda-go/blob/master/events/README_ApiGatewayCustomAuthorizer.md
func Handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	fmt.Printf("Processing incoming event: %v", event)
	alg := jwt.TokenServiceAlgorithmRS256

	// Field names must be case-insensitive.
	// https://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html#sec4.2.
	for k, v := range event.Headers {
		lowercasedKey := strings.ToLower(k)
		event.Headers[lowercasedKey] = v
	}

	token := event.Headers["authorization"]
	if len(token) < 50 {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid token in the header incoming data")
	}

	appID := event.Headers["x-identifo-clientid"]
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
	// If you want to use file instead of env variable, use this code instead.
	publicKey, err := jwt.LoadPublicKeyFromPEM("./public.pem", alg)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid public key: " + err.Error())
	}

	v := jwt.NewDefaultValidator(appID, jwtIssuer, "", "access")
	tokenV, err := jwt.ParseTokenWithPublicKey(string(tstr), publicKey)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error parsing token: " + err.Error())
	}
	if err := v.Validate(tokenV); err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Invalid token:" + err.Error())
	}

	return CreatePolicy(tokenV.UserID(), "Allow", event.MethodArn, nil), nil
}

// CreatePolicy is a helper function for generating an IAM policy.
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
