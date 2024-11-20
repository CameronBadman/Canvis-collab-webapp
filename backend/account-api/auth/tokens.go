package auth

import (
	"account-api/config"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"math/big"
	"net/http"
)

// GetCognitoPublicKey retrieves the public key from Cognito's JWKS endpoint
func getCognitoPublicKey(kid string) (*rsa.PublicKey, error) {
	// Cognito's JWKS endpoint URL
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", config.AwsRegion, config.UserPoolID)

	// Fetch the JWKS from the Cognito endpoint
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()

	// Read the JWKS response body
	jwksData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response: %v", err)
	}

	// Parse the JWKS
	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}

	if err := json.Unmarshal(jwksData, &jwks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWKS: %v", err)
	}

	// Iterate over the keys to find the one with the matching key ID
	var publicKey *rsa.PublicKey
	for _, key := range jwks.Keys {
		if key["kid"] == kid {
			// Assuming the key is an RSA key, if it's EC, handle it separately.
			// Extract the modulus and exponent to create the RSA public key
			modulus := key["n"].(string)
			exponent := key["e"].(string)

			// Convert the modulus and exponent to bytes
			nBytes, err := base64.RawURLEncoding.DecodeString(modulus)
			if err != nil {
				return nil, fmt.Errorf("failed to decode modulus: %v", err)
			}

			eBytes, err := base64.RawURLEncoding.DecodeString(exponent)
			if err != nil {
				return nil, fmt.Errorf("failed to decode exponent: %v", err)
			}

			// Construct the RSA public key from the modulus and exponent
			pubKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: int(new(big.Int).SetBytes(eBytes).Int64()),
			}

			publicKey = pubKey
			break
		}
	}

	if publicKey == nil {
		return nil, fmt.Errorf("no matching public key found for kid: %s", kid)
	}

	// Return the public key
	return publicKey, nil
}

// ExtractSubFromIDToken validates the ID token and extracts the user ID (sub)
func ExtractSubFromIDToken(idToken string) (string, error) {
	// Parse the token to extract the kid from the header
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Extract the kid from the token's header
	kid, ok := parsedToken.Header["kid"].(string)
	if !ok || kid == "" {
		return "", errors.New("token does not contain a valid kid")
	}

	// Get the public key from Cognito
	publicKey, err := getCognitoPublicKey(kid)
	if err != nil {
		return "", fmt.Errorf("failed to fetch public key: %v", err)
	}

	// Parse and validate the token using the public key
	parsedToken, err = jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is using the correct signing method
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		// Return the public key for verification
		return publicKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to validate token: %v", err)
	}

	// Extract the claims from the parsed token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return "", fmt.Errorf("invalid token claims or token not valid")
	}

	// Extract the user ID (sub) from the claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("user ID (sub) not found in token")
	}

	// Return the user ID (sub)
	return userID, nil
}
