// Copyright 2026 High Five Labs
// SPDX-License-Identifier: Apache-2.0

package extensions

import (
	"encoding/json"
	"net/http"

	"github.com/highfivelabs/krm-rest-runtime/rest"
)

// Vendor is how the manager reaches Cloudflare.
var Vendor = &rest.Vendor{
	Name:           "cloudflare",
	DefaultBaseURL: "https://api.cloudflare.com/client/v4",
	CredentialKeys: []string{"apiToken"},
	Authorize:      func(r *http.Request, c map[string]string) { r.Header.Set("Authorization", "Bearer "+c["apiToken"]) },
	NotFound:       notFound,
}

// notFoundCodes are Cloudflare error codes that mean "no such object" but
// come with a status other than 404 (found by the live run).
var notFoundCodes = map[int]bool{
	1472: true, // GET a deleted mTLS certificate: HTTP 400 "Certificate not found."
}

func notFound(status int, body []byte) bool {
	if status != http.StatusBadRequest {
		return false
	}
	var r struct {
		Errors []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	if json.Unmarshal(body, &r) != nil {
		return false
	}
	for _, e := range r.Errors {
		if notFoundCodes[e.Code] || e.Message == "requested snippet not found" {
			return true
		}
	}
	return false
}
