/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apisix

import (
	"github.com/kubernetes-sigs/ingress2gateway/pkg/i2gw"
	"github.com/kubernetes-sigs/ingress2gateway/pkg/i2gw/providers/common"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func httpToHTTPSFeature(ingressResources i2gw.InputResources, gatewayResources *i2gw.GatewayResources) field.ErrorList {
	var errs field.ErrorList
	httpToHTTPSAnnotation := apisixAnnotation("http-to-https")
	ruleGroups := common.GetRuleGroups(ingressResources.Ingresses)
	for _, rg := range ruleGroups {
		for _, rule := range rg.Rules {
			if val := rule.Ingress.Annotations[httpToHTTPSAnnotation]; val == "true" {
				if rule.Ingress.Spec.Rules == nil {
					continue
				}
				key := types.NamespacedName{Namespace: rule.Ingress.Namespace, Name: common.RouteName(rg.Name, rg.Host)}
				httpRoute, ok := gatewayResources.HTTPRoutes[key]
				if !ok {
					errs = append(errs, field.NotFound(field.NewPath("HTTPRoute"), key))
				}

				for i, rule := range httpRoute.Spec.Rules {
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type: gatewayv1.HTTPRouteFilterRequestRedirect,
						RequestRedirect: &gatewayv1.HTTPRequestRedirectFilter{
							Scheme:     ptr.To("https"),
							StatusCode: ptr.To(int(301)),
						},
					})
					httpRoute.Spec.Rules[i] = rule
				}
			}
		}
	}
	return errs
}
