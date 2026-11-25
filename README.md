# cloudflare-krm-extensions

Kubernetes kinds for the Cloudflare configuration that Terraform's
`cloudflare/cloudflare` provider manages and Cloudflare's operators do not. Cloudflare publishes no Kubernetes operator for its configuration. Community operators (StringKe's, adyanth's) cover Zero Trust, tunnels, DNS records, R2 and Pages.

By the numbers: `cloudflare/cloudflare` is the #21 Terraform provider by downloads
(395,685,520); version 5.29.0 has 289 resources. Operators
cover 25. This repository adds **258 kinds** standing in for
**260 more**. Each gap kind is listed below with its Terraform equivalent.

**Readiness: most kinds are proven against a fake Cloudflare API only.** Every
kind passes create, update, drift repair and delete against a fake built from
the same descriptors (`TestEveryKindOnFake`). Those marked `live` in
`zz_generated.report.md` passed the same against the real API too; treat the
others (`fake`) as untested there until they have. Many stand in for
Enterprise-only or paid features (Magic Transit, BYOIP, Logpush, DNS
firewall, Spectrum), which have only run on the fake.

## How it works

The kinds are generated, not hand-written. `hack/gen/vendor.yaml` pins
Cloudflare's OpenAPI spec; [krm-rest-runtime](https://github.com/highfivelabs/krm-rest-runtime)'s
generator finds every resource the API exposes in a create/read/update/delete
shape, keeps those that stand in for a Terraform resource in the gap list
(`hack/gen/terraform-gaps.json`), and writes:

- `types.yaml`: per kind, its API calls, path parameters and Terraform equivalents
- `config/crds/`: one CRD per kind, its spec the API's request body
- `zz_generated.report.md`: the full list

One generic manager (`cmd/manager`) runs every kind on krm-rest-runtime's
reconciler: it creates the resource, reads it back every 10 minutes and
puts back anything changed outside Kubernetes, and deletes it with the
object. Spec fields keep Cloudflare's own names, so the API reference
applies field for field. Path parameters are spec fields too.

Hand corrections to the automatic matching (a Terraform type matched to
the wrong API path, or not matched) are in `vendor.yaml`'s `overrides`.

## Install

```sh
kubectl kustomize config/default | KO_DOCKER_REPO=<registry> ko resolve -f - | kubectl apply --server-side -f -
```

Then, in each namespace that manages Cloudflare, a Secret with
`apiToken` and a Connection named `default`
(`config/samples/connection.yaml`). Objects use it unless
`spec.connectionRef.name` names another. Secret values (tokens, passwords,
secret values) go in a Kubernetes Secret, not the spec:
`spec.secretRefs: {<field>: {name, key}}` fills that body field.

To run only some API groups, install those groups' CRDs from `config/crds`
and pass `--groups` to the manager.

## Kinds

| Kind | Terraform equivalent | Shape | Read |
| --- | --- | --- | --- |
| `access.cloudflare.hfl.dev/v1alpha1` `AccessRule` | `cloudflare_access_rule` | crud | `GET /accounts/{account_id}/firewall/access_rules/rules/{rule_id}` |
| `account.cloudflare.hfl.dev/v1alpha1` `Account` | `cloudflare_account` | crud | `GET /accounts/{account_id}` |
| `account.cloudflare.hfl.dev/v1alpha1` `AccountDnsSettings` | `cloudflare_account_dns_settings` | settings | `GET /accounts/{account_id}/dns_settings` |
| `account.cloudflare.hfl.dev/v1alpha1` `AccountDnsSettingsInternalView` | `cloudflare_account_dns_settings_internal_view` | crud | `GET /accounts/{account_id}/dns_settings/views/{view_id}` |
| `account.cloudflare.hfl.dev/v1alpha1` `AccountMember` | `cloudflare_account_member` | crud | `GET /accounts/{account_id}/members/{member_id}` |
| `account.cloudflare.hfl.dev/v1alpha1` `AccountSubscription` | `cloudflare_account_subscription` | crud | `GET /accounts/{account_id}/subscriptions/{subscription_identifier}` |
| `account.cloudflare.hfl.dev/v1alpha1` `AccountToken` | `cloudflare_account_token` | crud | `GET /accounts/{account_id}/tokens/{token_id}` |
| `address.cloudflare.hfl.dev/v1alpha1` `AddressMap` | `cloudflare_address_map` | crud | `GET /accounts/{account_id}/addressing/address_maps/{address_map_id}` |
| `ai.cloudflare.hfl.dev/v1alpha1` `AiGateway` | `cloudflare_ai_gateway` | crud | `GET /accounts/{account_id}/ai-gateway/gateways/{id}` |
| `ai.cloudflare.hfl.dev/v1alpha1` `AiGatewayDynamicRouting` | `cloudflare_ai_gateway_dynamic_routing` | crud | `GET /accounts/{account_id}/ai-gateway/gateways/{gateway_id}/routes/{id}` |
| `ai.cloudflare.hfl.dev/v1alpha1` `AiSearchInstance` | `cloudflare_ai_search_instance` | crud | `GET /accounts/{account_id}/ai-search/instances/{id}` |
| `ai.cloudflare.hfl.dev/v1alpha1` `AiSearchNamespace` | `cloudflare_ai_search_namespace` | crud | `GET /accounts/{account_id}/ai-search/namespaces/{name}` |
| `ai.cloudflare.hfl.dev/v1alpha1` `AiSearchToken` | `cloudflare_ai_search_token` | crud | `GET /accounts/{account_id}/ai-search/tokens/{id}` |
| `ai.cloudflare.hfl.dev/v1alpha1` `AiSecurity` | `cloudflare_ai_security` | settings | `GET /zones/{zone_id}/ai-security/settings` |
| `ai.cloudflare.hfl.dev/v1alpha1` `AiSecurityCustomTopics` | `cloudflare_ai_security_custom_topics` | settings | `GET /zones/{zone_id}/ai-security/custom-topics` |
| `api.cloudflare.hfl.dev/v1alpha1` `ApiShield` | `cloudflare_api_shield` | settings | `GET /zones/{zone_id}/api_gateway/configuration` |
| `api.cloudflare.hfl.dev/v1alpha1` `ApiShieldDiscoveryOperation` | `cloudflare_api_shield_discovery_operation` | settings | `GET /zones/{zone_id}/api_gateway/discovery/operations/{discovery_id}` |
| `api.cloudflare.hfl.dev/v1alpha1` `ApiShieldOperation` | `cloudflare_api_shield_operation` | crud | `GET /zones/{zone_id}/api_gateway/operations/{operation_id}` |
| `api.cloudflare.hfl.dev/v1alpha1` `ApiShieldOperationSchemaValidationSettings` | `cloudflare_api_shield_operation_schema_validation_settings` | settings | `GET /zones/{zone_id}/api_gateway/operations/{operation_id}/schema_validation` |
| `api.cloudflare.hfl.dev/v1alpha1` `ApiShieldSchema` | `cloudflare_api_shield_schema` | crud | `GET /zones/{zone_id}/api_gateway/user_schemas/{schema_id}` |
| `api.cloudflare.hfl.dev/v1alpha1` `ApiShieldSchemaValidationSettings` | `cloudflare_api_shield_schema_validation_settings` | settings | `GET /zones/{zone_id}/api_gateway/settings/schema_validation` |
| `api.cloudflare.hfl.dev/v1alpha1` `ApiToken` | `cloudflare_api_token` | crud | `GET /user/tokens/{token_id}` |
| `argo.cloudflare.hfl.dev/v1alpha1` `ArgoSmartRouting` | `cloudflare_argo_smart_routing` | settings | `GET /zones/{zone_id}/argo/smart_routing` |
| `argo.cloudflare.hfl.dev/v1alpha1` `ArgoTieredCaching` | `cloudflare_argo_tiered_caching` | settings | `GET /zones/{zone_id}/argo/tiered_caching` |
| `authenticated.cloudflare.hfl.dev/v1alpha1` `AuthenticatedOriginPulls` | `cloudflare_authenticated_origin_pulls` | upsert | `GET /zones/{zone_id}/origin_tls_client_auth/hostnames/{hostname}` |
| `authenticated.cloudflare.hfl.dev/v1alpha1` `AuthenticatedOriginPullsCertificate` | `cloudflare_authenticated_origin_pulls_certificate` | crud | `GET /zones/{zone_id}/origin_tls_client_auth/{certificate_id}` |
| `authenticated.cloudflare.hfl.dev/v1alpha1` `AuthenticatedOriginPullsHostnameCertificate` | `cloudflare_authenticated_origin_pulls_hostname_certificate` | crud | `GET /zones/{zone_id}/origin_tls_client_auth/hostnames/certificates/{certificate_id}` |
| `authenticated.cloudflare.hfl.dev/v1alpha1` `AuthenticatedOriginPullsSettings` | `cloudflare_authenticated_origin_pulls_settings` | settings | `GET /zones/{zone_id}/origin_tls_client_auth/settings` |
| `bot.cloudflare.hfl.dev/v1alpha1` `BotManagement` | `cloudflare_bot_management` | settings | `GET /zones/{zone_id}/bot_management` |
| `byo.cloudflare.hfl.dev/v1alpha1` `ByoIpPrefix` | `cloudflare_byo_ip_prefix` | crud | `GET /accounts/{account_id}/addressing/prefixes/{prefix_id}` |
| `calls.cloudflare.hfl.dev/v1alpha1` `CallsSfuApp` | `cloudflare_calls_sfu_app` | crud | `GET /accounts/{account_id}/calls/apps/{app_id}` |
| `calls.cloudflare.hfl.dev/v1alpha1` `CallsTurnApp` | `cloudflare_calls_turn_app` | crud | `GET /accounts/{account_id}/calls/turn_keys/{key_id}` |
| `certificate.cloudflare.hfl.dev/v1alpha1` `CertificateAuthoritiesHostnameAssociations` | `cloudflare_certificate_authorities_hostname_associations` | settings | `GET /zones/{zone_id}/certificate_authorities/hostname_associations` |
| `certificate.cloudflare.hfl.dev/v1alpha1` `CertificatePack` | `cloudflare_certificate_pack` | crud | `GET /zones/{zone_id}/ssl/certificate_packs/{certificate_pack_id}` |
| `client.cloudflare.hfl.dev/v1alpha1` `ClientCertificate` | `cloudflare_client_certificate` | crud | `GET /certificates/{certificate_id}` |
| `cloud.cloudflare.hfl.dev/v1alpha1` `CloudConnectorRules` | `cloudflare_cloud_connector_rules` | settings | `GET /zones/{zone_id}/cloud_connector/rules` |
| `cloudforce.cloudflare.hfl.dev/v1alpha1` `CloudforceOneRequest` | `cloudflare_cloudforce_one_request` | crud | `GET /accounts/{account_id}/cloudforce-one/requests/{request_id}` |
| `cloudforce.cloudflare.hfl.dev/v1alpha1` `CloudforceOneRequestAsset` | `cloudflare_cloudforce_one_request_asset` | crud | `GET /accounts/{account_id}/cloudforce-one/requests/{request_id}/asset/{asset_id}` |
| `cloudforce.cloudflare.hfl.dev/v1alpha1` `CloudforceOneRequestMessage` | `cloudflare_cloudforce_one_request_message` | settings | `GET /accounts/{account_id}/cloudforce-one/v2/requests/{project_type}/{request_id}/messages/{message_id}` |
| `cloudforce.cloudflare.hfl.dev/v1alpha1` `CloudforceOneRequestPriority` | `cloudflare_cloudforce_one_request_priority` | crud | `GET /accounts/{account_id}/cloudforce-one/requests/priority/{priority_id}` |
| `connectivity.cloudflare.hfl.dev/v1alpha1` `ConnectivityDirectoryService` | `cloudflare_connectivity_directory_service` | crud | `GET /accounts/{account_id}/connectivity/directory/services/{service_id}` |
| `content.cloudflare.hfl.dev/v1alpha1` `ContentScanning` | `cloudflare_content_scanning` | settings | `GET /zones/{zone_id}/content-upload-scan/settings` |
| `content.cloudflare.hfl.dev/v1alpha1` `ContentScanningExpression` | `cloudflare_content_scanning_expression` | crud | `GET /zones/{zone_id}/content-upload-scan/payloads` |
| `ct.cloudflare.hfl.dev/v1alpha1` `CtAlerting` | `cloudflare_ct_alerting` | settings | `GET /zones/{zone_id}/ct/alerting` |
| `custom.cloudflare.hfl.dev/v1alpha1` `CustomCsr` | `cloudflare_custom_csr` | crud | `GET /zones/{zone_id}/custom_csrs/{custom_csr_id}` |
| `custom.cloudflare.hfl.dev/v1alpha1` `CustomHostname` | `cloudflare_custom_hostname` | crud | `GET /zones/{zone_id}/custom_hostnames/{custom_hostname_id}` |
| `custom.cloudflare.hfl.dev/v1alpha1` `CustomHostnameFallbackOrigin` | `cloudflare_custom_hostname_fallback_origin` | upsert | `GET /zones/{zone_id}/custom_hostnames/fallback_origin` |
| `custom.cloudflare.hfl.dev/v1alpha1` `CustomOriginTrustStore` | `cloudflare_custom_origin_trust_store` | crud | `GET /zones/{zone_id}/acm/custom_trust_store/{custom_origin_trust_store_id}` |
| `custom.cloudflare.hfl.dev/v1alpha1` `CustomPageAsset` | `cloudflare_custom_page_asset` | crud | `GET /zones/{zone_id}/custom_pages/assets/{asset_name}` |
| `custom.cloudflare.hfl.dev/v1alpha1` `CustomPages` | `cloudflare_custom_pages` | settings | `GET /zones/{zone_id}/custom_pages/{identifier}` |
| `custom.cloudflare.hfl.dev/v1alpha1` `CustomSsl` | `cloudflare_custom_ssl` | crud | `GET /zones/{zone_id}/custom_certificates/{custom_certificate_id}` |
| `d1.cloudflare.hfl.dev/v1alpha1` `D1Database` | `cloudflare_d1_database` | crud | `GET /accounts/{account_id}/d1/database/{database_id}` |
| `dls.cloudflare.hfl.dev/v1alpha1` `DlsPrefixBinding` | `cloudflare_dls_prefix_binding` | crud | `GET /accounts/{account_id}/dls/regional_services/prefix_bindings/{binding_id}` |
| `dns.cloudflare.hfl.dev/v1alpha1` `DnsFirewall` | `cloudflare_dns_firewall` | crud | `GET /accounts/{account_id}/dns_firewall/{dns_firewall_id}` |
| `dns.cloudflare.hfl.dev/v1alpha1` `DnsZoneTransfersAcl` | `cloudflare_dns_zone_transfers_acl` | crud | `GET /accounts/{account_id}/secondary_dns/acls/{acl_id}` |
| `dns.cloudflare.hfl.dev/v1alpha1` `DnsZoneTransfersIncoming` | `cloudflare_dns_zone_transfers_incoming` | upsert | `GET /zones/{zone_id}/secondary_dns/incoming` |
| `dns.cloudflare.hfl.dev/v1alpha1` `DnsZoneTransfersOutgoing` | `cloudflare_dns_zone_transfers_outgoing` | upsert | `GET /zones/{zone_id}/secondary_dns/outgoing` |
| `dns.cloudflare.hfl.dev/v1alpha1` `DnsZoneTransfersPeer` | `cloudflare_dns_zone_transfers_peer` | crud | `GET /accounts/{account_id}/secondary_dns/peers/{peer_id}` |
| `dns.cloudflare.hfl.dev/v1alpha1` `DnsZoneTransfersTsig` | `cloudflare_dns_zone_transfers_tsig` | crud | `GET /accounts/{account_id}/secondary_dns/tsigs/{tsig_id}` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailRoutingAddress` | `cloudflare_email_routing_address` | crud | `GET /accounts/{account_id}/email/routing/addresses/{destination_address_identifier}` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailRoutingCatchAll` | `cloudflare_email_routing_catch_all` | settings | `GET /zones/{zone_id}/email/routing/rules/catch_all` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailRoutingDns` | `cloudflare_email_routing_dns` | upsert | `GET /zones/{zone_id}/email/routing` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailRoutingRule` | `cloudflare_email_routing_rule` | crud | `GET /zones/{zone_id}/email/routing/rules/{rule_identifier}` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailRoutingSettings` | `cloudflare_email_routing_settings` | settings | `GET /zones/{zone_id}/email/routing` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailSecurityAllowPolicy` | `cloudflare_email_security_allow_policy` | crud | `GET /accounts/{account_id}/email-security/settings/allow_policies/{policy_id}` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailSecurityBlockSender` | `cloudflare_email_security_block_sender` | crud | `GET /accounts/{account_id}/email-security/settings/block_senders/{pattern_id}` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailSecurityDomain` | `cloudflare_email_security_domain` | crud | `GET /accounts/{account_id}/email-security/settings/domains/{domain_id}` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailSecurityImpersonationRegistry` | `cloudflare_email_security_impersonation_registry` | crud | `GET /accounts/{account_id}/email-security/settings/impersonation_registry/{impersonation_registry_id}` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailSecurityTrustedDomains` | `cloudflare_email_security_trusted_domains` | crud | `GET /accounts/{account_id}/email-security/settings/trusted_domains/{trusted_domain_id}` |
| `email.cloudflare.hfl.dev/v1alpha1` `EmailSendingSubdomain` | `cloudflare_email_sending_subdomain` | crud | `GET /zones/{zone_id}/email/sending/subdomains/{subdomain_id}` |
| `field.cloudflare.hfl.dev/v1alpha1` `FieldExtractor` | `cloudflare_field_extractor` | upsert | `GET /accounts/{account_id}/field_extractors/{extractor}` |
| `flagship.cloudflare.hfl.dev/v1alpha1` `FlagshipApp` | `cloudflare_flagship_app` | crud | `GET /accounts/{account_id}/flagship/apps/{app_id}` |
| `flagship.cloudflare.hfl.dev/v1alpha1` `FlagshipFlag` | `cloudflare_flagship_flag` | crud | `GET /accounts/{account_id}/flagship/apps/{app_id}/flags/{flag_key}` |
| `google.cloudflare.hfl.dev/v1alpha1` `GoogleTagGateway` | `cloudflare_google_tag_gateway` | settings | `GET /zones/{zone_id}/settings/google-tag-gateway/config` |
| `healthcheck.cloudflare.hfl.dev/v1alpha1` `Healthcheck` | `cloudflare_healthcheck` | crud | `GET /zones/{zone_id}/healthchecks/{healthcheck_id}` |
| `hostname.cloudflare.hfl.dev/v1alpha1` `HostnameTlsSetting` | `cloudflare_hostname_tls_setting` | upsert | `GET /zones/{zone_id}/hostnames/settings/{setting_id}/{hostname}` |
| `hyperdrive.cloudflare.hfl.dev/v1alpha1` `HyperdriveConfig` | `cloudflare_hyperdrive_config` | crud | `GET /accounts/{account_id}/hyperdrive/configs/{hyperdrive_id}` |
| `image.cloudflare.hfl.dev/v1alpha1` `Image` | `cloudflare_image` | crud | `GET /accounts/{account_id}/images/v1/{image_id}` |
| `image.cloudflare.hfl.dev/v1alpha1` `ImageVariant` | `cloudflare_image_variant` | crud | `GET /accounts/{account_id}/images/v1/variants/{variant_id}` |
| `keyless.cloudflare.hfl.dev/v1alpha1` `KeylessCertificate` | `cloudflare_keyless_certificate` | crud | `GET /zones/{zone_id}/keyless_certificates/{keyless_certificate_id}` |
| `leaked.cloudflare.hfl.dev/v1alpha1` `LeakedCredentialCheck` | `cloudflare_leaked_credential_check` | settings | `GET /zones/{zone_id}/leaked-credential-checks` |
| `leaked.cloudflare.hfl.dev/v1alpha1` `LeakedCredentialCheckRule` | `cloudflare_leaked_credential_check_rule` | crud | `GET /zones/{zone_id}/leaked-credential-checks/detections/{detection_id}` |
| `list.cloudflare.hfl.dev/v1alpha1` `List` | `cloudflare_list` | crud | `GET /accounts/{account_id}/rules/lists/{list_id}` |
| `list.cloudflare.hfl.dev/v1alpha1` `ListItem` | `cloudflare_list_item` | crud | `GET /accounts/{account_id}/rules/lists/{list_id}/items/{id}` |
| `load.cloudflare.hfl.dev/v1alpha1` `LoadBalancer` | `cloudflare_load_balancer` | crud | `GET /zones/{zone_id}/load_balancers/{load_balancer_id}` |
| `load.cloudflare.hfl.dev/v1alpha1` `LoadBalancerMonitor` | `cloudflare_load_balancer_monitor` | crud | `GET /accounts/{account_id}/load_balancers/monitors/{monitor_id}` |
| `load.cloudflare.hfl.dev/v1alpha1` `LoadBalancerMonitorGroup` | `cloudflare_load_balancer_monitor_group` | crud | `GET /accounts/{account_id}/load_balancers/monitor_groups/{monitor_group_id}` |
| `load.cloudflare.hfl.dev/v1alpha1` `LoadBalancerPool` | `cloudflare_load_balancer_pool` | crud | `GET /accounts/{account_id}/load_balancers/pools/{pool_id}` |
| `log.cloudflare.hfl.dev/v1alpha1` `LogExplorerDataset` | `cloudflare_log_explorer_dataset` | crud | `GET /zones/{zone_id}/logs/explorer/datasets/{dataset_id}` |
| `logpull.cloudflare.hfl.dev/v1alpha1` `LogpullRetention` | `cloudflare_logpull_retention` | settings | `GET /zones/{zone_id}/logs/control/retention/flag` |
| `logpush.cloudflare.hfl.dev/v1alpha1` `LogpushJob` | `cloudflare_logpush_job` | crud | `GET /zones/{zone_id}/logpush/jobs/{job_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicNetworkMonitoringConfiguration` | `cloudflare_magic_network_monitoring_configuration` | upsert | `GET /accounts/{account_id}/mnm/config` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicNetworkMonitoringRule` | `cloudflare_magic_network_monitoring_rule` | crud | `GET /accounts/{account_id}/mnm/rules/{rule_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicTransitCf1Site` | `cloudflare_magic_transit_cf1_site` | crud | `GET /accounts/{account_id}/magic/cf1_sites/{cf1_site_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicTransitConnector` | `cloudflare_magic_transit_connector` | crud | `GET /accounts/{account_id}/magic/connectors/{connector_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicTransitSite` | `cloudflare_magic_transit_site` | crud | `GET /accounts/{account_id}/magic/sites/{site_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicTransitSiteAcl` | `cloudflare_magic_transit_site_acl` | crud | `GET /accounts/{account_id}/magic/sites/{site_id}/acls/{acl_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicTransitSiteLan` | `cloudflare_magic_transit_site_lan` | crud | `GET /accounts/{account_id}/magic/sites/{site_id}/lans/{lan_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicTransitSiteWan` | `cloudflare_magic_transit_site_wan` | crud | `GET /accounts/{account_id}/magic/sites/{site_id}/wans/{wan_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicWanBgpFilterProfile` | `cloudflare_magic_wan_bgp_filter_profile` | crud | `GET /accounts/{account_id}/magic/bgp/filter_profiles/{profile_id}` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicWanGreTunnel` | `cloudflare_magic_wan_gre_tunnel` | crud | `GET /accounts/{account_id}/magic/gre_tunnels` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicWanIpsecTunnel` | `cloudflare_magic_wan_ipsec_tunnel` | crud | `GET /accounts/{account_id}/magic/ipsec_tunnels` |
| `magic.cloudflare.hfl.dev/v1alpha1` `MagicWanStaticRoute` | `cloudflare_magic_wan_static_route` | crud | `GET /accounts/{account_id}/magic/routes` |
| `managed.cloudflare.hfl.dev/v1alpha1` `ManagedTransforms` | `cloudflare_managed_transforms` | upsert | `GET /zones/{zone_id}/managed_headers` |
| `monetization.cloudflare.hfl.dev/v1alpha1` `MonetizationRules` | `cloudflare_monetization_rules` | upsert | `GET /zones/{zone_id}/monetization/rules` |
| `moq.cloudflare.hfl.dev/v1alpha1` `MoqRelay` | `cloudflare_moq_relay` | crud | `GET /accounts/{account_id}/moq/relays/{relay_id}` |
| `mtls.cloudflare.hfl.dev/v1alpha1` `MtlsCertificate` | `cloudflare_mtls_certificate` | crud | `GET /accounts/{account_id}/mtls_certificates/{mtls_certificate_id}` |
| `nel.cloudflare.hfl.dev/v1alpha1` `NelSetting` | `cloudflare_nel_setting` | settings | `GET /zones/{zone_id}/settings/nel` |
| `notification.cloudflare.hfl.dev/v1alpha1` `NotificationPolicy` | `cloudflare_notification_policy` | crud | `GET /accounts/{account_id}/alerting/v3/policies/{policy_id}` |
| `notification.cloudflare.hfl.dev/v1alpha1` `NotificationPolicyWebhooks` | `cloudflare_notification_policy_webhooks` | crud | `GET /accounts/{account_id}/alerting/v3/destinations/webhooks/{webhook_id}` |
| `oauth.cloudflare.hfl.dev/v1alpha1` `OauthClient` | `cloudflare_oauth_client` | crud | `GET /accounts/{account_id}/oauth_clients/{oauth_client_id}` |
| `observatory.cloudflare.hfl.dev/v1alpha1` `ObservatoryScheduledTest` | `cloudflare_observatory_scheduled_test` | upsert | `GET /zones/{zone_id}/speed_api/schedule/{url}` |
| `organization.cloudflare.hfl.dev/v1alpha1` `Organization` | `cloudflare_organization` | crud | `GET /organizations/{organization_id}` |
| `organization.cloudflare.hfl.dev/v1alpha1` `OrganizationProfile` | `cloudflare_organization_profile` | settings | `GET /organizations/{organization_id}/profile` |
| `origin.cloudflare.hfl.dev/v1alpha1` `OriginCloudRegion` | `cloudflare_origin_cloud_region` | upsert | `GET /zones/{zone_id}/origin/cloud_regions/{origin_ip}` |
| `origin.cloudflare.hfl.dev/v1alpha1` `OriginTlsComplianceModes` | `cloudflare_origin_tls_compliance_modes` | upsert | `GET /zones/{zone_id}/settings/origin_tls_compliance_modes` |
| `page.cloudflare.hfl.dev/v1alpha1` `PageRule` | `cloudflare_page_rule` | crud | `GET /zones/{zone_id}/pagerules/{pagerule_id}` |
| `page.cloudflare.hfl.dev/v1alpha1` `PageShieldPolicy` | `cloudflare_page_shield_policy` | crud | `GET /zones/{zone_id}/page_shield/policies/{policy_id}` |
| `pipeline.cloudflare.hfl.dev/v1alpha1` `Pipeline` | `cloudflare_pipeline` | crud | `GET /accounts/{account_id}/pipelines/{pipeline_name}` |
| `pipeline.cloudflare.hfl.dev/v1alpha1` `PipelineSink` | `cloudflare_pipeline_sink` | crud | `GET /accounts/{account_id}/pipelines/v1/sinks/{sink_id}` |
| `pipeline.cloudflare.hfl.dev/v1alpha1` `PipelineStream` | `cloudflare_pipeline_stream` | crud | `GET /accounts/{account_id}/pipelines/v1/streams/{stream_id}` |
| `precursor.cloudflare.hfl.dev/v1alpha1` `Precursor` | `cloudflare_precursor` | settings | `GET /zones/{zone_id}/precursor` |
| `queue.cloudflare.hfl.dev/v1alpha1` `Queue` | `cloudflare_queue` | crud | `GET /accounts/{account_id}/queues/{queue_id}` |
| `queue.cloudflare.hfl.dev/v1alpha1` `QueueConsumer` | `cloudflare_queue_consumer` | crud | `GET /accounts/{account_id}/queues/{queue_id}/consumers/{consumer_id}` |
| `r2.cloudflare.hfl.dev/v1alpha1` `R2BucketCors` | `cloudflare_r2_bucket_cors` | upsert | `GET /accounts/{account_id}/r2/buckets/{bucket_name}/cors` |
| `r2.cloudflare.hfl.dev/v1alpha1` `R2BucketLifecycle` | `cloudflare_r2_bucket_lifecycle` | settings | `GET /accounts/{account_id}/r2/buckets/{bucket_name}/lifecycle` |
| `r2.cloudflare.hfl.dev/v1alpha1` `R2BucketLock` | `cloudflare_r2_bucket_lock` | upsert | `GET /accounts/{account_id}/r2/buckets/{bucket_name}/lock` |
| `r2.cloudflare.hfl.dev/v1alpha1` `R2BucketSippy` | `cloudflare_r2_bucket_sippy` | upsert | `GET /accounts/{account_id}/r2/buckets/{bucket_name}/sippy` |
| `r2.cloudflare.hfl.dev/v1alpha1` `R2DataCatalog` | `cloudflare_r2_data_catalog` | upsert | `GET /accounts/{account_id}/r2-catalog/{bucket_name}` |
| `r2.cloudflare.hfl.dev/v1alpha1` `R2ManagedDomain` | `cloudflare_r2_managed_domain` | settings | `GET /accounts/{account_id}/r2/buckets/{bucket_name}/domains/managed` |
| `regional.cloudflare.hfl.dev/v1alpha1` `RegionalHostname` | `cloudflare_regional_hostname` | crud | `GET /zones/{zone_id}/addressing/regional_hostnames/{hostname}` |
| `regional.cloudflare.hfl.dev/v1alpha1` `RegionalTieredCache` | `cloudflare_regional_tiered_cache` | settings | `GET /zones/{zone_id}/cache/regional_tiered_cache` |
| `schema.cloudflare.hfl.dev/v1alpha1` `SchemaValidationOperationSettings` | `cloudflare_schema_validation_operation_settings` | settings | `GET /zones/{zone_id}/schema_validation/settings/operations` |
| `schema.cloudflare.hfl.dev/v1alpha1` `SchemaValidationSchemas` | `cloudflare_schema_validation_schemas` | crud | `GET /zones/{zone_id}/schema_validation/schemas/{schema_id}` |
| `schema.cloudflare.hfl.dev/v1alpha1` `SchemaValidationSettings` | `cloudflare_schema_validation_settings` | settings | `GET /zones/{zone_id}/schema_validation/settings` |
| `secrets.cloudflare.hfl.dev/v1alpha1` `SecretsStore` | `cloudflare_secrets_store` | crud | `GET /accounts/{account_id}/secrets_store/stores/{store_id}` |
| `secrets.cloudflare.hfl.dev/v1alpha1` `SecretsStoreSecret` | `cloudflare_secrets_store_secret` | crud | `GET /accounts/{account_id}/secrets_store/stores/{store_id}/secrets/{secret_id}` |
| `share.cloudflare.hfl.dev/v1alpha1` `Share` | `cloudflare_share` | crud | `GET /accounts/{account_id}/shares/{share_id}` |
| `share.cloudflare.hfl.dev/v1alpha1` `ShareRecipient` | `cloudflare_share_recipient` | crud | `GET /accounts/{account_id}/shares/{share_id}/recipients/{recipient_id}` |
| `share.cloudflare.hfl.dev/v1alpha1` `ShareResource` | `cloudflare_share_resource` | crud | `GET /accounts/{account_id}/shares/{share_id}/resources/{share_resource_id}` |
| `snippet.cloudflare.hfl.dev/v1alpha1` `Snippet` | `cloudflare_snippet`, `cloudflare_snippets` | upsert | `GET /zones/{zone_id}/snippets/{snippet_name}` |
| `snippet.cloudflare.hfl.dev/v1alpha1` `SnippetRules` | `cloudflare_snippet_rules` | upsert | `GET /zones/{zone_id}/snippets/snippet_rules` |
| `spectrum.cloudflare.hfl.dev/v1alpha1` `SpectrumApplication` | `cloudflare_spectrum_application` | crud | `GET /zones/{zone_id}/spectrum/apps/{app_id}` |
| `sso.cloudflare.hfl.dev/v1alpha1` `SsoConnector` | `cloudflare_sso_connector` | crud | `GET /accounts/{account_id}/sso_connectors/{sso_connector_id}` |
| `stream.cloudflare.hfl.dev/v1alpha1` `Stream` | `cloudflare_stream` | crud | `GET /accounts/{account_id}/stream/{identifier}` |
| `stream.cloudflare.hfl.dev/v1alpha1` `StreamAudioTrack` | `cloudflare_stream_audio_track` | crud | `GET /accounts/{account_id}/stream/{identifier}/audio` |
| `stream.cloudflare.hfl.dev/v1alpha1` `StreamCaptionLanguage` | `cloudflare_stream_caption_language` | upsert | `GET /accounts/{account_id}/stream/{identifier}/captions/{language}` |
| `stream.cloudflare.hfl.dev/v1alpha1` `StreamDownload` | `cloudflare_stream_download` | upsert | `GET /accounts/{account_id}/stream/{identifier}/downloads` |
| `stream.cloudflare.hfl.dev/v1alpha1` `StreamKey` | `cloudflare_stream_key` | crud | `GET /accounts/{account_id}/stream/keys` |
| `stream.cloudflare.hfl.dev/v1alpha1` `StreamLiveInput` | `cloudflare_stream_live_input` | crud | `GET /accounts/{account_id}/stream/live_inputs/{live_input_identifier}` |
| `stream.cloudflare.hfl.dev/v1alpha1` `StreamWatermark` | `cloudflare_stream_watermark` | crud | `GET /accounts/{account_id}/stream/watermarks/{identifier}` |
| `stream.cloudflare.hfl.dev/v1alpha1` `StreamWebhook` | `cloudflare_stream_webhook` | upsert | `GET /accounts/{account_id}/stream/webhook` |
| `tiered.cloudflare.hfl.dev/v1alpha1` `TieredCache` | `cloudflare_tiered_cache` | upsert | `GET /zones/{zone_id}/cache/tiered_cache_smart_topology_enable` |
| `token.cloudflare.hfl.dev/v1alpha1` `TokenValidationConfig` | `cloudflare_token_validation_config` | crud | `GET /zones/{zone_id}/token_validation/config/{config_id}` |
| `token.cloudflare.hfl.dev/v1alpha1` `TokenValidationRules` | `cloudflare_token_validation_rules` | crud | `GET /zones/{zone_id}/token_validation/rules/{rule_id}` |
| `total.cloudflare.hfl.dev/v1alpha1` `TotalTls` | `cloudflare_total_tls` | settings | `GET /zones/{zone_id}/acm/total_tls` |
| `turnstile.cloudflare.hfl.dev/v1alpha1` `TurnstileWidget` | `cloudflare_turnstile_widget` | crud | `GET /accounts/{account_id}/challenges/widgets/{sitekey}` |
| `universal.cloudflare.hfl.dev/v1alpha1` `UniversalSslSetting` | `cloudflare_universal_ssl_setting` | settings | `GET /zones/{zone_id}/ssl/universal/settings` |
| `url.cloudflare.hfl.dev/v1alpha1` `UrlNormalizationSettings` | `cloudflare_url_normalization_settings` | upsert | `GET /zones/{zone_id}/url_normalization` |
| `user.cloudflare.hfl.dev/v1alpha1` `User` | `cloudflare_user` | settings | `GET /user` |
| `user.cloudflare.hfl.dev/v1alpha1` `UserAgentBlockingRule` | `cloudflare_user_agent_blocking_rule` | crud | `GET /zones/{zone_id}/firewall/ua_rules/{ua_rule_id}` |
| `user.cloudflare.hfl.dev/v1alpha1` `UserGroup` | `cloudflare_user_group` | crud | `GET /accounts/{account_id}/iam/user_groups/{user_group_id}` |
| `user.cloudflare.hfl.dev/v1alpha1` `UserGroupMembers` | `cloudflare_user_group_members` | crud | `GET /accounts/{account_id}/iam/user_groups/{user_group_id}/members/{member_id}` |
| `vulnerability.cloudflare.hfl.dev/v1alpha1` `VulnerabilityScannerCredential` | `cloudflare_vulnerability_scanner_credential` | crud | `GET /accounts/{account_id}/vuln_scanner/credential_sets/{credential_set_id}/credentials/{credential_id}` |
| `vulnerability.cloudflare.hfl.dev/v1alpha1` `VulnerabilityScannerCredentialSet` | `cloudflare_vulnerability_scanner_credential_set` | crud | `GET /accounts/{account_id}/vuln_scanner/credential_sets/{credential_set_id}` |
| `vulnerability.cloudflare.hfl.dev/v1alpha1` `VulnerabilityScannerTargetEnvironment` | `cloudflare_vulnerability_scanner_target_environment` | crud | `GET /accounts/{account_id}/vuln_scanner/target_environments/{target_environment_id}` |
| `waiting.cloudflare.hfl.dev/v1alpha1` `WaitingRoom` | `cloudflare_waiting_room` | crud | `GET /zones/{zone_id}/waiting_rooms/{waiting_room_id}` |
| `waiting.cloudflare.hfl.dev/v1alpha1` `WaitingRoomEvent` | `cloudflare_waiting_room_event` | crud | `GET /zones/{zone_id}/waiting_rooms/{waiting_room_id}/events/{event_id}` |
| `waiting.cloudflare.hfl.dev/v1alpha1` `WaitingRoomRules` | `cloudflare_waiting_room_rules` | settings | `GET /zones/{zone_id}/waiting_rooms/{waiting_room_id}/rules` |
| `waiting.cloudflare.hfl.dev/v1alpha1` `WaitingRoomSettings` | `cloudflare_waiting_room_settings` | settings | `GET /zones/{zone_id}/waiting_rooms/settings` |
| `web.cloudflare.hfl.dev/v1alpha1` `WebAnalyticsRule` | `cloudflare_web_analytics_rule` | crud | `GET /accounts/{account_id}/rum/v2/{ruleset_id}/rules` |
| `web.cloudflare.hfl.dev/v1alpha1` `WebAnalyticsSite` | `cloudflare_web_analytics_site` | crud | `GET /accounts/{account_id}/rum/site_info/{site_id}` |
| `web3.cloudflare.hfl.dev/v1alpha1` `Web3Hostname` | `cloudflare_web3_hostname` | crud | `GET /zones/{zone_id}/web3/hostnames/{identifier}` |
| `worker.cloudflare.hfl.dev/v1alpha1` `Worker` | `cloudflare_worker` | crud | `GET /accounts/{account_id}/workers/workers/{worker_id}` |
| `worker.cloudflare.hfl.dev/v1alpha1` `WorkerVersion` | `cloudflare_worker_version` | crud | `GET /accounts/{account_id}/workers/workers/{worker_id}/versions/{version_id}` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersCronTrigger` | `cloudflare_workers_cron_trigger` | settings | `GET /accounts/{account_id}/workers/scripts/{script_name}/schedules` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersCustomDomain` | `cloudflare_workers_custom_domain` | crud | `GET /accounts/{account_id}/workers/domains/{domain_id}` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersDeployment` | `cloudflare_workers_deployment` | crud | `GET /accounts/{account_id}/workers/scripts/{script_name}/deployments/{deployment_id}` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersForPlatformsDispatchNamespace` | `cloudflare_workers_for_platforms_dispatch_namespace` | crud | `GET /accounts/{account_id}/workers/dispatch/namespaces/{dispatch_namespace}` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersKv` | `cloudflare_workers_kv` | upsert | `GET /accounts/{account_id}/storage/kv/namespaces/{namespace_id}/values/{key_name}` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersKvNamespace` | `cloudflare_workers_kv_namespace` | crud | `GET /accounts/{account_id}/storage/kv/namespaces/{namespace_id}` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersRoute` | `cloudflare_workers_route` | crud | `GET /zones/{zone_id}/workers/routes/{route_id}` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersScript` | `cloudflare_workers_script` | upsert | `GET /accounts/{account_id}/workers/scripts/{script_name}` |
| `workers.cloudflare.hfl.dev/v1alpha1` `WorkersScriptSubdomain` | `cloudflare_workers_script_subdomain` | upsert | `GET /accounts/{account_id}/workers/scripts/{script_name}/subdomain` |
| `workflow.cloudflare.hfl.dev/v1alpha1` `Workflow` | `cloudflare_workflow` | upsert | `GET /accounts/{account_id}/workflows/{workflow_name}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessAiControlsMcpPortal` | `cloudflare_zero_trust_access_ai_controls_mcp_portal` | crud | `GET /accounts/{account_id}/access/ai-controls/mcp/portals/{id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessAiControlsMcpServer` | `cloudflare_zero_trust_access_ai_controls_mcp_server` | crud | `GET /accounts/{account_id}/access/ai-controls/mcp/servers/{id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessCustomPage` | `cloudflare_zero_trust_access_custom_page` | crud | `GET /accounts/{account_id}/access/custom_pages/{custom_page_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessInfrastructureTarget` | `cloudflare_zero_trust_access_infrastructure_target` | crud | `GET /accounts/{account_id}/infrastructure/targets/{target_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessKeyConfiguration` | `cloudflare_zero_trust_access_key_configuration` | settings | `GET /accounts/{account_id}/access/keys` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessMtlsCertificate` | `cloudflare_zero_trust_access_mtls_certificate` | crud | `GET /accounts/{account_id}/access/certificates/{certificate_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessMtlsHostnameSettings` | `cloudflare_zero_trust_access_mtls_hostname_settings` | settings | `GET /accounts/{account_id}/access/certificates/settings` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessShortLivedCertificate` | `cloudflare_zero_trust_access_short_lived_certificate` | upsert | `GET /accounts/{account_id}/access/apps/{app_id}/ca` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustAccessTag` | `cloudflare_zero_trust_access_tag` | crud | `GET /accounts/{account_id}/access/tags/{tag_name}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCasbFindingType` | `cloudflare_zero_trust_casb_finding_type` | crud | `GET /accounts/{account_id}/data-security/posture/finding_types/{finding_type_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCasbIntegration` | `cloudflare_zero_trust_casb_integration` | crud | `GET /accounts/{account_id}/one/integrations/{id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCasbPolicy` | `cloudflare_zero_trust_casb_policy` | crud | `GET /accounts/{account_id}/data-security/posture/policies/{policy_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCasbWebhook` | `cloudflare_zero_trust_casb_webhook` | crud | `GET /accounts/{account_id}/data-security/posture/webhooks/{webhook_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustConnectivitySettings` | `cloudflare_zero_trust_connectivity_settings` | settings | `GET /accounts/{account_id}/zerotrust/connectivity_settings` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsBasicAuthCredential` | `cloudflare_zero_trust_cws_basic_auth_credential` | crud | `GET /accounts/{account_id}/cws/basic-auth-credentials/{credential_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsBasicAuthSecret` | `cloudflare_zero_trust_cws_basic_auth_secret` | crud | `GET /accounts/{account_id}/cws/basic-auth-credentials/{credential_id}/secrets/{secret_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsDeployment` | `cloudflare_zero_trust_cws_deployment` | crud | `GET /accounts/{account_id}/cws/deployments/{deployment_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsDomainAcl` | `cloudflare_zero_trust_cws_domain_acl` | crud | `GET /accounts/{account_id}/cws/domain-acls/{domain_acl_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsPolicy` | `cloudflare_zero_trust_cws_policy` | crud | `GET /accounts/{account_id}/cws/policies/{policy_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsPolicySet` | `cloudflare_zero_trust_cws_policy_set` | crud | `GET /accounts/{account_id}/cws/policy-sets/{policy_set_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsRegoMetadata` | `cloudflare_zero_trust_cws_rego_metadata` | crud | `GET /accounts/{account_id}/cws/rego-metadata/{rego_metadata_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsRegoModule` | `cloudflare_zero_trust_cws_rego_module` | crud | `GET /accounts/{account_id}/cws/rego-modules/{rego_module_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsTenantStructure` | `cloudflare_zero_trust_cws_tenant_structure` | crud | `GET /accounts/{account_id}/cws/tenant-structures/{structure_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsTenantStructureNode` | `cloudflare_zero_trust_cws_tenant_structure_node` | crud | `GET /accounts/{account_id}/cws/tenant-structures/{structure_id}/nodes/{node_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustCwsWorkload` | `cloudflare_zero_trust_cws_workload` | crud | `GET /accounts/{account_id}/cws/workloads/{workload_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDeviceCustomProfileLocalDomainFallback` | `cloudflare_zero_trust_device_custom_profile_local_domain_fallback` | upsert | `GET /accounts/{account_id}/devices/policy/{policy_id}/fallback_domains` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDeviceDefaultProfileCertificates` | `cloudflare_zero_trust_device_default_profile_certificates` | settings | `GET /zones/{zone_id}/devices/policy/certificates` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDeviceDefaultProfileLocalDomainFallback` | `cloudflare_zero_trust_device_default_profile_local_domain_fallback` | upsert | `GET /accounts/{account_id}/devices/policy/fallback_domains` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDeviceDeploymentGroups` | `cloudflare_zero_trust_device_deployment_groups` | crud | `GET /accounts/{account_id}/devices/deployment-groups/{group_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDeviceIpProfile` | `cloudflare_zero_trust_device_ip_profile` | crud | `GET /accounts/{account_id}/devices/ip-profiles/{profile_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDeviceManagedNetworks` | `cloudflare_zero_trust_device_managed_networks` | crud | `GET /accounts/{account_id}/devices/networks/{network_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDevicePostureIntegration` | `cloudflare_zero_trust_device_posture_integration` | crud | `GET /accounts/{account_id}/devices/posture/integration/{integration_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDeviceSettings` | `cloudflare_zero_trust_device_settings` | upsert | `GET /accounts/{account_id}/devices/settings` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDeviceSubnet` | `cloudflare_zero_trust_device_subnet` | crud | `GET /accounts/{account_id}/zerotrust/subnets/warp/{subnet_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDexRule` | `cloudflare_zero_trust_dex_rule` | crud | `GET /accounts/{account_id}/dex/rules/{rule_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDexTest` | `cloudflare_zero_trust_dex_test` | crud | `GET /accounts/{account_id}/dex/devices/dex_tests/{dex_test_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpCustomProfile` | `cloudflare_zero_trust_dlp_custom_profile` | crud | `GET /accounts/{account_id}/dlp/profiles/custom/{profile_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpDataClass` | `cloudflare_zero_trust_dlp_data_class` | crud | `GET /accounts/{account_id}/dlp/data_classes/{data_class_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpDataTag` | `cloudflare_zero_trust_dlp_data_tag` | crud | `GET /accounts/{account_id}/dlp/data_tag_categories/{category_id}/data_tags/{tag_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpDataTagCategory` | `cloudflare_zero_trust_dlp_data_tag_category` | crud | `GET /accounts/{account_id}/dlp/data_tag_categories/{category_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpDataset` | `cloudflare_zero_trust_dlp_dataset` | crud | `GET /accounts/{account_id}/dlp/datasets/{dataset_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpEntry` | `cloudflare_zero_trust_dlp_custom_entry`, `cloudflare_zero_trust_dlp_entry` | crud | `GET /accounts/{account_id}/dlp/entries/{entry_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpIntegrationEntry` | `cloudflare_zero_trust_dlp_integration_entry` | crud | `GET /accounts/{account_id}/dlp/entries` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpPredefinedEntry` | `cloudflare_zero_trust_dlp_predefined_entry` | crud | `GET /accounts/{account_id}/dlp/entries` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpPredefinedProfile` | `cloudflare_zero_trust_dlp_predefined_profile` | upsert | `GET /accounts/{account_id}/dlp/profiles/predefined/{profile_id}/config` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpProximityValidator` | `cloudflare_zero_trust_dlp_proximity_validator` | crud | `GET /accounts/{account_id}/dlp/validators/proximity_words/{validator_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpSensitivityGroup` | `cloudflare_zero_trust_dlp_sensitivity_group` | crud | `GET /accounts/{account_id}/dlp/sensitivity_groups/{sensitivity_group_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpSensitivityLevel` | `cloudflare_zero_trust_dlp_sensitivity_level` | crud | `GET /accounts/{account_id}/dlp/sensitivity_groups/{sensitivity_group_id}/levels/{sensitivity_level_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpSensitivityLevelOrder` | `cloudflare_zero_trust_dlp_sensitivity_level_order` | settings | `GET /accounts/{account_id}/dlp/sensitivity_groups/{sensitivity_group_id}/level_order` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDlpSettings` | `cloudflare_zero_trust_dlp_settings` | upsert | `GET /accounts/{account_id}/dlp/settings` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustDnsLocation` | `cloudflare_zero_trust_dns_location` | crud | `GET /accounts/{account_id}/gateway/locations/{location_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustGatewayCertificate` | `cloudflare_zero_trust_gateway_certificate` | crud | `GET /accounts/{account_id}/gateway/certificates/{certificate_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustGatewayLogging` | `cloudflare_zero_trust_gateway_logging` | settings | `GET /accounts/{account_id}/gateway/logging` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustGatewayPacfile` | `cloudflare_zero_trust_gateway_pacfile` | crud | `GET /accounts/{account_id}/gateway/pacfiles/{pacfile_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustGatewayProxyEndpoint` | `cloudflare_zero_trust_gateway_proxy_endpoint` | crud | `GET /accounts/{account_id}/gateway/proxy_endpoints/{proxy_endpoint_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustNetworkHostnameRoute` | `cloudflare_zero_trust_network_hostname_route` | crud | `GET /accounts/{account_id}/zerotrust/routes/hostname/{hostname_route_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustOrganization` | `cloudflare_zero_trust_organization` | settings | `GET /accounts/{account_id}/access/organizations` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustOrganizationTrustSettings` | `cloudflare_zero_trust_organization_trust_settings` | upsert | `GET /accounts/{account_id}/access/organizations/trust_settings` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustResourceLibraryApplication` | `cloudflare_zero_trust_resource_library_application` | crud | `GET /accounts/{account_id}/resource-library/applications/{id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustRiskBehavior` | `cloudflare_zero_trust_risk_behavior` | settings | `GET /accounts/{account_id}/zt_risk_scoring/behaviors` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustRiskScoringIntegration` | `cloudflare_zero_trust_risk_scoring_integration` | crud | `GET /accounts/{account_id}/zt_risk_scoring/integrations/{integration_id}` |
| `zerotrust.cloudflare.hfl.dev/v1alpha1` `ZeroTrustTunnelWarpConnectorConfig` | `cloudflare_zero_trust_tunnel_warp_connector_config` | settings | `GET /accounts/{account_id}/warp_connector/{tunnel_id}/configurations` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneAutoOriginTlsKex` | `cloudflare_zone_auto_origin_tls_kex` | settings | `GET /zones/{zone_id}/settings/auto_origin_tls_kex` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneCacheReserve` | `cloudflare_zone_cache_reserve` | settings | `GET /zones/{zone_id}/cache/cache_reserve` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneCacheVariants` | `cloudflare_zone_cache_variants` | upsert | `GET /zones/{zone_id}/cache/variants` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneDnsSettings` | `cloudflare_zone_dns_settings` | settings | `GET /zones/{zone_id}/dns_settings` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneDnssec` | `cloudflare_zone_dnssec` | upsert | `GET /zones/{zone_id}/dnssec` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneHold` | `cloudflare_zone_hold` | upsert | `GET /zones/{zone_id}/hold` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneLockdown` | `cloudflare_zone_lockdown` | crud | `GET /zones/{zone_id}/firewall/lockdowns/{lock_downs_id}` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneSetting` | `cloudflare_zone_setting` | settings | `GET /zones/{zone_id}/settings/{setting_id}` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneSubscription` | `cloudflare_zone_subscription` | upsert | `GET /zones/{zone_id}/subscription` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneTracing` | `cloudflare_zone_tracing` | upsert | `GET /zones/{zone_id}/observability/tracing/settings` |
| `zone.cloudflare.hfl.dev/v1alpha1` `ZoneTracingRules` | `cloudflare_zone_tracing_rules` | upsert | `GET /zones/{zone_id}/observability/tracing/rules` |

## Not covered

<!-- BEGIN generated: not covered (krmgen.py, from vendor.yaml hardBlocks and draftNotes) -->
4 of the 264 Terraform gap types have no kind: 4 are not built (`hardBlocks` in `hack/gen/vendor.yaml`).

| Terraform | Status | Why |
| --- | --- | --- |
| `cloudflare_filter` | not built | A retired API; the live run's answer is 410 (`filters.api.deprecated`; use rulesets). |
| `cloudflare_firewall_rule` | not built | A retired API; the live run's answer is 410 (`firewallrules.api.deprecated`; use rulesets). |
| `cloudflare_logpush_ownership_challenge` | not built | A one-off action (`POST .../logpush/ownership` sends a challenge file to a destination); the provider's resource has no read, update or delete. |
| `cloudflare_rate_limit` | not built | A retired API; the live run's answer is 410 (use rulesets). |
<!-- END generated: not covered -->

## Kinds with a limit worth knowing

- `ListItem` is added by an asynchronous bulk POST that answers no ID; it is then found in the list by its value (the provider's `?search=`, up to 500 a page), and looked for before it is added again. It cannot be changed: delete and recreate it.
- `WorkersCustomDomain`, `R2DataCatalog` and `StreamDownload` cannot be changed in place, as the provider's; `ZoneHold` is made by a POST (its `include_subdomains` then set by the PATCH), and reads `hold: false` as none; `R2DataCatalog` reads `status: inactive` as none.
- `EmailRoutingDns` turns Email Routing on (`POST .../email/routing/dns`, which adds its DNS records) and off (`DELETE`), and reads whether it is on as `enabled` in the zone's Email Routing settings (`GET .../email/routing`, which `EmailRoutingSettings` reads and writes too): `enabled: false` is none. It cannot be changed in place; the provider's optional `name` is not sent.
- `AuthenticatedOriginPulls` writes one hostname's association with the bulk `PUT .../hostnames`; deleting the object sends it with `enabled: null`, which voids it.

- `MagicWanGreTunnel`, `MagicWanIpsecTunnel`, `MagicWanStaticRoute` read the account's list (`.../gre_tunnels` etc.): a read of one answers it inside `result.gre_tunnel`, a create in `result`, and a kind has one envelope. An update answers `result.modified_gre_tunnel`, so `status.atProvider` shows that answer until the next read.
- `WebAnalyticsRule`, `ZeroTrustDlpPredefinedEntry` and `ZeroTrustDlpIntegrationEntry` are read from their list too (the API has no read of one rule; the provider reads an entry at `.../dlp/entries/{entry_id}`, the same object as in the list).
- `StreamAudioTrack` is added by `POST .../audio/copy` (from a URL, as the API adds tracks) and found by `uid` in the video's list at `.../audio`. The provider's resource instead edits a track that already exists.
- `ManagedTransforms` lists the managed transforms it sets; the API reads back every one the zone has, and each listed one is compared with the one of the same `id`, in any order, as the provider does. Transforms not listed are left as they are (the provider turns off those it does not list that are on). The delete resets them all.
- `SecretsStoreSecret`, `MagicTransitCf1Site` and `ContentScanningExpression` are made by a bulk `POST` of a one-element list, whose answer's first element is the new object; `ContentScanningExpression` is read from the zone's list (the API has no read of one).
- `WorkersKv` writes a key's `value` (and `metadata`) as a multipart form and reads the value back as its raw bytes; `metadata` (read at another path) is written each time the spec changes, not compared. `expiration` and `expiration_ttl` (query parameters) are not set.
- `StreamCaptionLanguage` uploads a WebVTT `file` (`{name, content}`); it is not read back. The provider can instead ask the API to generate captions (`POST .../captions/{language}/generate`), which the kind does not.
- `CloudforceOneRequestAsset` is made by a multipart `POST .../asset/new` (the spec's form has `source` only) and changed by a JSON `PUT`; its read answers a one-element list.
- `ZoneDnssec`, `ZoneCacheVariants` and `ZoneTracing` are turned on by a PATCH; deleting the object resets them (`ZoneDnssec`'s delete removes the zone's DNSSEC records).
- `ZeroTrustOrganization` is the account's organization (the zone-level one is not covered), changed by PATCH; it must exist already (made in the dashboard, or by the API's `POST`), and deleting the object leaves it, as the provider does.
- `WorkersCronTrigger` holds all of a Worker's cron schedules: `PUT` takes a bare array, and the answer is read at `result.schedules`. Deleting the object leaves them, as the provider does.
- `WaitingRoomRules` holds all of a room's rules; deleting the object sets an empty list, as the provider does. `ZeroTrustDlpPredefinedProfile`'s delete is a no-op in the API: the profile stays as last set.
- `ObservatoryScheduledTest` schedules with the API's default `region` and `frequency`: those are query parameters, which a kind cannot set yet.
- `LeakedCredentialCheck`, `TotalTls`, `LogpullRetention`, `WorkersScriptSubdomain` and `ZeroTrustAccessShortLivedCertificate` are written by a POST to the path they are read from.

## For agents (MCP)

`cmd/mcp` serves this repository's kinds over the Model Context Protocol
(stdio), from `types.yaml`, `config/crds` and the live tests' records.
`.mcp.json` registers it, so an agent session opened in this repository
has it without setup. Elsewhere, run it from the repository root:

```sh
go run ./cmd/mcp
```

Its tools: `list_kinds`, `find_kinds_for_terraform`, `describe_kind` (every
field, the REST calls, the live record), `example` (the live-tested spec
where there is one), `validate_manifest` (the CRD schema, as the API server
checks it) and `plan_on_fake` (the controller run against the fake API,
with every REST call it makes). Nothing it does reaches the real API.
krm-rest-runtime's README ("For agents") has the details.

## Regenerate

```sh
hack/gen/generate.sh   # needs python3 + pyyaml, and krm-rest-runtime beside this repo
go test ./...
```
