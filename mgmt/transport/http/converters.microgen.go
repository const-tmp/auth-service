// Code generated by microgen 1.0.5. DO NOT EDIT.

// Please, do not change functions names!
package transporthttp

import (
	transport "auth/mgmt/transport"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
)

func CommonHTTPRequestEncoder(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func CommonHTTPResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func _Decode_CreateService_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateServiceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_GetAllServices_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.GetAllServicesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_GetService_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.GetServiceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_CreateAccount_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_CreateAccountWithName_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateAccountWithNameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_GetAllAccounts_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.GetAllAccountsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_GetAccount_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.GetAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_UpdateAccount_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.UpdateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_AttachAccountToService_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.AttachAccountToServiceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_RemoveAccountFromService_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.RemoveAccountFromServiceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_CreatePermission_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreatePermissionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_GetPermission_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.GetPermissionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_GetAllPermission_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.GetAllPermissionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_GetFilteredPermissions_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.GetFilteredPermissionsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_DeletePermission_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.DeletePermissionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_GetUserPermissions_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.GetUserPermissionsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_AddUserPermission_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.AddUserPermissionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_RemoveUserPermission_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.RemoveUserPermissionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func _Decode_CreateService_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.CreateServiceResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_GetAllServices_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GetAllServicesResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_GetService_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GetServiceResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_CreateAccount_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.CreateAccountResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_CreateAccountWithName_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.CreateAccountWithNameResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_GetAllAccounts_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GetAllAccountsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_GetAccount_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GetAccountResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_UpdateAccount_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.UpdateAccountResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_AttachAccountToService_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.AttachAccountToServiceResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_RemoveAccountFromService_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.RemoveAccountFromServiceResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_CreatePermission_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.CreatePermissionResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_GetPermission_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GetPermissionResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_GetAllPermission_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GetAllPermissionResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_GetFilteredPermissions_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GetFilteredPermissionsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_DeletePermission_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.DeletePermissionResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_GetUserPermissions_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GetUserPermissionsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_AddUserPermission_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.AddUserPermissionResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Decode_RemoveUserPermission_Response(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.RemoveUserPermissionResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func _Encode_CreateService_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "create-service")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_GetAllServices_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "get-all-services")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_GetService_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "get-service")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_CreateAccount_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "create-account")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_CreateAccountWithName_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "create-account-with-name")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_GetAllAccounts_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "get-all-accounts")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_GetAccount_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "get-account")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_UpdateAccount_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "update-account")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_AttachAccountToService_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "attach-account-toservice")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_RemoveAccountFromService_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "remove-account-from-service")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_CreatePermission_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "create-permission")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_GetPermission_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "get-permission")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_GetAllPermission_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "get-all-permission")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_GetFilteredPermissions_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "get-filtered-permissions")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_DeletePermission_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "delete-permission")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_GetUserPermissions_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "get-user-permissions")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_AddUserPermission_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "add-user-permission")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_RemoveUserPermission_Request(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = path.Join(r.URL.Path, "remove-user-permission")
	return CommonHTTPRequestEncoder(ctx, r, request)
}

func _Encode_CreateService_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_GetAllServices_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_GetService_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_CreateAccount_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_CreateAccountWithName_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_GetAllAccounts_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_GetAccount_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_UpdateAccount_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_AttachAccountToService_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_RemoveAccountFromService_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_CreatePermission_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_GetPermission_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_GetAllPermission_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_GetFilteredPermissions_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_DeletePermission_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_GetUserPermissions_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_AddUserPermission_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}

func _Encode_RemoveUserPermission_Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return CommonHTTPResponseEncoder(ctx, w, response)
}