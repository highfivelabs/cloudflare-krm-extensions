// Copyright 2026 High Five Labs
// SPDX-License-Identifier: Apache-2.0

package extensions

import (
	"context"

	"github.com/highfivelabs/krm-rest-runtime/descriptor"
	"github.com/highfivelabs/krm-rest-runtime/reconciler"
	"github.com/highfivelabs/krm-rest-runtime/rest"
)

// Hooks change request bodies the descriptors cannot describe.
var Hooks = map[string][]reconciler.Hook{
	// PUT /zones/{zone_id}/origin/cloud_regions/{origin_ip} wants the IP in
	// the body too ("must match the path parameter", else 400 "Invalid or
	// missing origin IP address"); a path parameter is never in the body.
	"OriginCloudRegion": {ParamInBody("origin_ip")},
}

// ParamInBody copies a path parameter into the request body.
func ParamInBody(name string) reconciler.Hook {
	return func(_ context.Context, _ *rest.Client, _ *descriptor.Kind, params map[string]string, body map[string]any) error {
		if v, ok := params[name]; ok {
			body[name] = v
		}
		return nil
	}
}
