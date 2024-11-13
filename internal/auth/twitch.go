package auth

import (
	"bot/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    TokenType   string `json:"token_type"`
}

type UserResponse struct {
	UserId string `json:"id"`
	ViewCount int `json:"view_count"`
}



func GetTwitchUserToken(config *config.Config, authorizationCode, redirectURI string) (string, error) {
    url := "https://id.twitch.tv/oauth2/token"
    
    data := map[string]string{
        "client_id":     config.ClientId,
        "client_secret": config.ClientSecret,
        "code":          authorizationCode,
        "grant_type":    "authorization_code",
        "redirect_uri":  redirectURI,
    }
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    var tokenResponse TokenResponse
    err = json.Unmarshal(body, &tokenResponse)
    if err != nil {
        return "", err
    }
    
    return tokenResponse.AccessToken, nil
}




func GetTwitchUserId(config *config.Config, accessToken string) (string, error) {
	url := "https://api.twitch.tv/helix/users?login=joindev"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Client-Id", config.ClientId)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var userResponse struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return "", err
	}

	if len(userResponse.Data) == 0 {
		return "", fmt.Errorf("user not found")
	}


	return userResponse.Data[0].ID, nil
}

func GetTwitchToken(config *config.Config) (string, error) {
	url:= "https://id.twitch.tv/oauth2/token"

	data := map[string]string{
		"client_id": config.ClientId,
		"client_secret": config.ClientSecret,
		"grant_type": "client_credentials",
	}
	jsonData, err := json.Marshal(data)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))



    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var tokenResponse TokenResponse
    err = json.Unmarshal(body, &tokenResponse)
    if err != nil {
        return "", err
    }


    return tokenResponse.AccessToken, nil
}
