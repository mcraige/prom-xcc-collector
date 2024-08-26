package main

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

type RedfishClient struct {
    BaseURL    string
    Username   string
    Password   string
    HTTPClient *http.Client
}

type Response struct {
    IndicatorLED string `json:"IndicatorLED"`
}

func NewRedfishClient(baseURL, username, password string) *RedfishClient {
    return &RedfishClient{
        BaseURL:  baseURL,
        Username: username,
        Password: password,
        HTTPClient: &http.Client{
            Timeout: time.Second * 10,
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
            },
        },
    }
}

func (c *RedfishClient) Login() error {
    // Implement login logic if needed
    return nil
}

func (c *RedfishClient) Logout() {
    // Implement logout logic if needed
}

func (c *RedfishClient) GetChassisIndicatorLED() (string, error) {
    req, err := http.NewRequest("GET", c.BaseURL+"/redfish/v1/Chassis/1", nil)
    if err != nil {
        return "", err
    }
    req.SetBasicAuth(c.Username, c.Password)

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("error code %d", resp.StatusCode)
    }

    var response Response
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", err
    }

    return response.IndicatorLED, nil
}

func main() {
    if len(os.Args) != 4 {
        log.Fatalf("Usage: %s <BMC IP> <Username> <Password>", os.Args[0])
    }

    ip := os.Args[1]
    username := os.Args[2]
    password := os.Args[3]

    client := NewRedfishClient("https://"+ip, username, password)
    if err := client.Login(); err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    defer client.Logout()

    led, err := client.GetChassisIndicatorLED()
    if err != nil {
        log.Fatalf("Failed to get chassis indicator LED: %v", err)
    }

    result := map[string]interface{}{
        "ret": true,
        "msg": led,
    }
    jsonResult, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(jsonResult))
}
