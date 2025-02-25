package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Register(cognitoClient *cognitoidentityprovider.Client, userPoolClientID string, clientID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}

		userAttres := []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(req.Email),
			},
			{
				Name:  aws.String("username"),
				Value: aws.String(req.Username),
			},
			{
				Name:  aws.String("password"),
				Value: aws.String(req.Password),
			},
		}

		input := &cognitoidentityprovider.SignUpInput{
			ClientId: aws.String(userPoolClientID),
			Username: aws.String(userPoolClientID),
			Password: aws.String(userPoolClientID),
			UserAttributes: []types.AttributeType{
				{
					Name:  aws.String("email"),
					Value: aws.String(req.Email),
				},
			},
		}

		result, err := cognitoClient.SignUp(c, input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"userSub":       *result.UserSub,
			"userConfirmed": result.UserConfirmed,
		})
	}
}
