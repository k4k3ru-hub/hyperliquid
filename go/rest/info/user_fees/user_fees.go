//
// user_fees.go
//
package user_fees

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    myConstant "github.com/k4k3ru-hub/hyperliquid/go/constant"
    myRestDTO "github.com/k4k3ru-hub/hyperliquid/go/rest/dto"
)


const (
    TypeValue = "userFees"
)


//
// Client.
//
type Client struct {
    parent      ParentClient
    endpointURL string
    httpMethod  string
    reqBody     *myRestDTO.RequestBody
    httpHeader  http.Header
}

type UserFees struct {
    DailyUserVlm                []DailyUserVlmEntry    `json:"dailyUserVlm"`
    FeeSchedule                 FeeSchedule            `json:"feeSchedule"`
    UserCrossRate               string                 `json:"userCrossRate"`
    UserAddRate                 string                 `json:"userAddRate"`
    UserSpotCrossRate           string                 `json:"userSpotCrossRate"`
    UserSpotAddRate             string                 `json:"userSpotAddRate"`
    ActiveReferralDiscount      string                 `json:"activeReferralDiscount"`
    Trial                       any                    `json:"trial"`
    FeeTrialReward              string                 `json:"feeTrialReward"`
    NextTrialAvailableTimestamp any                    `json:"nextTrialAvailableTimestamp"`
    StakingLink                 StakingLink            `json:"stakingLink"`
    ActiveStakingDiscount       ActiveStakingDiscount  `json:"activeStakingDiscount"`
}

type DailyUserVlmEntry struct {
    Date      string `json:"date"`
    UserCross string `json:"userCross"`
    UserAdd   string `json:"userAdd"`
    Exchange  string `json:"exchange"`
}

type FeeSchedule struct {
    Cross                 string                 `json:"cross"`
    Add                   string                 `json:"add"`
    SpotCross             string                 `json:"spotCross"`
    SpotAdd               string                 `json:"spotAdd"`
    Tiers                 FeeScheduleTiers       `json:"tiers"`
    ReferralDiscount      string                 `json:"referralDiscount"`
    StakingDiscountTiers  []StakingDiscountTier  `json:"stakingDiscountTiers"`
}

type FeeScheduleTiers struct {
    VIP []VIPTier `json:"vip"`
    MM  []MMTier  `json:"mm"`
}

type VIPTier struct {
    NtlCutoff string `json:"ntlCutoff"`
    Cross     string `json:"cross"`
    Add       string `json:"add"`
    SpotCross string `json:"spotCross"`
    SpotAdd   string `json:"spotAdd"`
}

type MMTier struct {
    MakerFractionCutoff string `json:"makerFractionCutoff"`
    Add                 string `json:"add"`
}

type StakingDiscountTier struct {
    BpsOfMaxSupply string `json:"bpsOfMaxSupply"`
    Discount       string `json:"discount"`
}

type StakingLink struct {
    Type        string `json:"type"`
    StakingUser string `json:"stakingUser"`
}

type ActiveStakingDiscount struct {
    BpsOfMaxSupply string `json:"bpsOfMaxSupply"`
    Discount       string `json:"discount"`
}

type ParentClient interface {
    SetEndpointURL(endpointURL string)
    SetHttpMethod(method string)
    SetHttpHeader(header http.Header)
    SetBody(body *myRestDTO.RequestBody)
    Send(context.Context) ([]byte, error)
}


//
// New Client.
//
// Version:
//   - 2026-04-05: Added.
//
func NewClient(parentClient ParentClient, user string) (*Client, error) {
    // Guard.
    if parentClient == nil {
        return nil, fmt.Errorf("failed to create user fees client: missing required value: parent_client=null.")
    }

    // Create request body.
    reqBody := &myRestDTO.RequestBody{
        Type: TypeValue,
        User: user,
    }

    // Create http header.
    httpHeader := http.Header{
        "Content-Type": {myConstant.ContentTypeJson},
    }

    return &Client{
        parent: parentClient,
        endpointURL: myConstant.BaseUrlRest + myConstant.ApiEndpointInfo,
        httpMethod: http.MethodPost,
        reqBody: reqBody,
        httpHeader: httpHeader,
    }, nil
}


//
// Send a request.
//
func (c *Client) Send(ctx context.Context) (*UserFees, error) {
    // Set date to parent client.
    c.parent.SetEndpointURL(c.endpointURL)
    c.parent.SetHttpMethod(c.httpMethod)
    c.parent.SetHttpHeader(c.httpHeader)
    c.parent.SetBody(c.reqBody)

    // Send a request.
    resBody, err := c.parent.Send(ctx)
    if err != nil {
        return nil, err
    }

    // Parse JSON data.
    result := &UserFees{}
    if err := json.Unmarshal(resBody, result); err != nil {
        return nil, fmt.Errorf("failed to send: %w: response=%q", err, string(resBody))
    }

    return result, nil
}

