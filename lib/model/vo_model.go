package model

import (
	"github.com/Authing/authing-go-sdk/lib/enum"
)

type ListMemberRequest struct {
	NodeId               string `json:"nodeId"`
	Page                 int    `json:"page"`
	Limit                int    `json:"limit"`
	IncludeChildrenNodes bool   `json:"includeChildrenNodes"`
}

type UserDetailData struct {
	ThirdPartyIdentity User `json:"thirdPartyIdentity"`
}

type UserDetailResponse struct {
	Message string `json:"message"`
	Code    int64  `json:"code"`
	Data    User   `json:"data"`
}

type ExportAllOrganizationResponse struct {
	Message string `json:"message"`
	Code    int64  `json:"code"`
	Data    []Node `json:"data"`
}

type NodeByIdDetail struct {
	NodeById Node `json:"nodeById"`
}

type NodeByIdResponse struct {
	Data NodeByIdDetail `json:"data"`
}

type QueryListRequest struct {
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
	SortBy enum.SortByEnum `json:"sortBy"`
}

type Users struct {
	Users PaginatedUsers `json:"users"`
}
type ListUserResponse struct {
	Data Users `json:"data"`
}

type ListOrganizationResponse struct {
	Message string        `json:"message"`
	Code    int64         `json:"code"`
	Data    PaginatedOrgs `json:"data"`
}

type GetOrganizationByIdData struct {
	Org Org `json:"org"`
}

type GetOrganizationByIdResponse struct {
	Data GetOrganizationByIdData `json:"data"`
}

type GetRoleListRequest struct {
	Page      int             `json:"page"`
	Limit     int             `json:"limit"`
	SortBy    enum.SortByEnum `json:"sortBy"`
	Namespace string          `json:"namespace"`
}

type Roles struct {
	Roles PaginatedRoles `json:"roles"`
}
type GetRoleListResponse struct {
	Data Roles `json:"data"`
}

type GetRoleUserListRequest struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Code      string `json:"code"`
	Namespace string `json:"namespace"`
}

type ValidateTokenRequest struct {
	AccessToken      string    `json:"accessToken"`
	IdToken     string    `json:"idToken"`
}

type ClientCredentialInput struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type GetAccessTokenByClientCredentialsRequest struct {
	Scope string `json:"scope"`
	ClientCredentialInput *ClientCredentialInput `json:"client_credential_input"`
}

