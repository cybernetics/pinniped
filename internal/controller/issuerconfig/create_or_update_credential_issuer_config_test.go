/*
Copyright 2020 VMware, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package issuerconfig

import (
	"context"
	"fmt"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	coretesting "k8s.io/client-go/testing"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/suzerain-io/pinniped/generated/1.19/apis/crdpinniped/v1alpha1"
	crdpinnipedv1alpha1 "github.com/suzerain-io/pinniped/generated/1.19/apis/crdpinniped/v1alpha1"
	pinnipedfake "github.com/suzerain-io/pinniped/generated/1.19/client/clientset/versioned/fake"
)

func TestCreateOrUpdateCredentialIssuerConfig(t *testing.T) {
	spec.Run(t, "specs", func(t *testing.T, when spec.G, it spec.S) {
		var r *require.Assertions
		var ctx context.Context
		var pinnipedAPIClient *pinnipedfake.Clientset
		var credentialIssuerConfigGVR schema.GroupVersionResource
		const installationNamespace = "some-namespace"
		const configName = "pinniped-config"

		it.Before(func() {
			r = require.New(t)
			ctx = context.Background()
			pinnipedAPIClient = pinnipedfake.NewSimpleClientset()
			credentialIssuerConfigGVR = schema.GroupVersionResource{
				Group:    crdpinnipedv1alpha1.GroupName,
				Version:  crdpinnipedv1alpha1.SchemeGroupVersion.Version,
				Resource: "credentialissuerconfigs",
			}
		})

		when("the config does not exist", func() {
			it("creates a new config which includes only the updates made by the func parameter", func() {
				err := CreateOrUpdateCredentialIssuerConfig(ctx, installationNamespace, pinnipedAPIClient,
					func(configToUpdate *v1alpha1.CredentialIssuerConfig) {
						configToUpdate.Status.KubeConfigInfo = &crdpinnipedv1alpha1.CredentialIssuerConfigKubeConfigInfo{
							CertificateAuthorityData: "some-ca-value",
						}
					},
				)
				r.NoError(err)

				expectedGetAction := coretesting.NewGetAction(credentialIssuerConfigGVR, installationNamespace, configName)

				expectedCreateAction := coretesting.NewCreateAction(
					credentialIssuerConfigGVR,
					installationNamespace,
					&crdpinnipedv1alpha1.CredentialIssuerConfig{
						TypeMeta: metav1.TypeMeta{},
						ObjectMeta: metav1.ObjectMeta{
							Name:      configName,
							Namespace: installationNamespace,
						},
						Status: crdpinnipedv1alpha1.CredentialIssuerConfigStatus{
							Strategies: []crdpinnipedv1alpha1.CredentialIssuerConfigStrategy{},
							KubeConfigInfo: &crdpinnipedv1alpha1.CredentialIssuerConfigKubeConfigInfo{
								Server:                   "",
								CertificateAuthorityData: "some-ca-value",
							},
						},
					},
				)

				r.Equal([]coretesting.Action{expectedGetAction, expectedCreateAction}, pinnipedAPIClient.Actions())
			})

			when("there is an unexpected error while creating the existing object", func() {
				it.Before(func() {
					pinnipedAPIClient.PrependReactor("create", "credentialissuerconfigs", func(_ coretesting.Action) (bool, runtime.Object, error) {
						return true, nil, fmt.Errorf("error on create")
					})
				})

				it("returns an error", func() {
					err := CreateOrUpdateCredentialIssuerConfig(ctx, installationNamespace, pinnipedAPIClient,
						func(configToUpdate *v1alpha1.CredentialIssuerConfig) {},
					)
					r.EqualError(err, "could not create or update credentialissuerconfig: create failed: error on create")
				})
			})
		})

		when("the config already exists", func() {
			var existingConfig *crdpinnipedv1alpha1.CredentialIssuerConfig

			it.Before(func() {
				existingConfig = &crdpinnipedv1alpha1.CredentialIssuerConfig{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Name:      configName,
						Namespace: installationNamespace,
					},
					Status: crdpinnipedv1alpha1.CredentialIssuerConfigStatus{
						Strategies: []crdpinnipedv1alpha1.CredentialIssuerConfigStrategy{
							{
								Type:           crdpinnipedv1alpha1.KubeClusterSigningCertificateStrategyType,
								Status:         crdpinnipedv1alpha1.SuccessStrategyStatus,
								Reason:         crdpinnipedv1alpha1.FetchedKeyStrategyReason,
								Message:        "initial-message",
								LastUpdateTime: metav1.Now(),
							},
						},
						KubeConfigInfo: &crdpinnipedv1alpha1.CredentialIssuerConfigKubeConfigInfo{
							Server:                   "initial-server-value",
							CertificateAuthorityData: "initial-ca-value",
						},
					},
				}
				r.NoError(pinnipedAPIClient.Tracker().Add(existingConfig))
			})

			it("updates the existing config to only apply the updates made by the func parameter", func() {
				err := CreateOrUpdateCredentialIssuerConfig(ctx, installationNamespace, pinnipedAPIClient,
					func(configToUpdate *v1alpha1.CredentialIssuerConfig) {
						configToUpdate.Status.KubeConfigInfo.CertificateAuthorityData = "new-ca-value"
					},
				)
				r.NoError(err)

				expectedGetAction := coretesting.NewGetAction(credentialIssuerConfigGVR, installationNamespace, configName)

				// Only the edited field should be changed.
				expectedUpdatedConfig := existingConfig.DeepCopy()
				expectedUpdatedConfig.Status.KubeConfigInfo.CertificateAuthorityData = "new-ca-value"
				expectedUpdateAction := coretesting.NewUpdateAction(credentialIssuerConfigGVR, installationNamespace, expectedUpdatedConfig)

				r.Equal([]coretesting.Action{expectedGetAction, expectedUpdateAction}, pinnipedAPIClient.Actions())
			})

			it("avoids the cost of an update if the local updates made by the func parameter did not actually change anything", func() {
				err := CreateOrUpdateCredentialIssuerConfig(ctx, installationNamespace, pinnipedAPIClient,
					func(configToUpdate *v1alpha1.CredentialIssuerConfig) {
						configToUpdate.Status.KubeConfigInfo.CertificateAuthorityData = "initial-ca-value"
					},
				)
				r.NoError(err)

				expectedGetAction := coretesting.NewGetAction(credentialIssuerConfigGVR, installationNamespace, configName)
				r.Equal([]coretesting.Action{expectedGetAction}, pinnipedAPIClient.Actions())
			})

			when("there is an unexpected error while getting the existing object", func() {
				it.Before(func() {
					pinnipedAPIClient.PrependReactor("get", "credentialissuerconfigs", func(_ coretesting.Action) (bool, runtime.Object, error) {
						return true, nil, fmt.Errorf("error on get")
					})
				})

				it("returns an error", func() {
					err := CreateOrUpdateCredentialIssuerConfig(ctx, installationNamespace, pinnipedAPIClient,
						func(configToUpdate *v1alpha1.CredentialIssuerConfig) {},
					)
					r.EqualError(err, "could not create or update credentialissuerconfig: get failed: error on get")
				})
			})

			when("there is an unexpected error while updating the existing object", func() {
				it.Before(func() {
					pinnipedAPIClient.PrependReactor("update", "credentialissuerconfigs", func(_ coretesting.Action) (bool, runtime.Object, error) {
						return true, nil, fmt.Errorf("error on update")
					})
				})

				it("returns an error", func() {
					err := CreateOrUpdateCredentialIssuerConfig(ctx, installationNamespace, pinnipedAPIClient,
						func(configToUpdate *v1alpha1.CredentialIssuerConfig) {
							configToUpdate.Status.KubeConfigInfo.CertificateAuthorityData = "new-ca-value"
						},
					)
					r.EqualError(err, "could not create or update credentialissuerconfig: error on update")
				})
			})

			when("there is a conflict error while updating the existing object on the first try and the next try succeeds", func() {
				var slightlyDifferentExistingConfig *crdpinnipedv1alpha1.CredentialIssuerConfig

				it.Before(func() {
					hit := false
					slightlyDifferentExistingConfig = existingConfig.DeepCopy()
					slightlyDifferentExistingConfig.Status.KubeConfigInfo.Server = "some-other-server-value-from-conflicting-update"

					pinnipedAPIClient.PrependReactor("update", "credentialissuerconfigs", func(_ coretesting.Action) (bool, runtime.Object, error) {
						// Return an error on the first call, then fall through to the default (successful) response.
						if !hit {
							// Before the update fails, also change the object that will be returned by the next Get(),
							// to make sure that the production code does a fresh Get() after detecting a conflict.
							r.NoError(pinnipedAPIClient.Tracker().Update(credentialIssuerConfigGVR, slightlyDifferentExistingConfig, installationNamespace))
							hit = true
							return true, nil, apierrors.NewConflict(schema.GroupResource{
								Group:    apiregistrationv1.GroupName,
								Resource: "credentialissuerconfigs",
							}, "alphav1.pinniped.dev", fmt.Errorf("there was a conflict"))
						}
						return false, nil, nil
					})
				})

				it("retries updates on conflict", func() {
					err := CreateOrUpdateCredentialIssuerConfig(ctx, installationNamespace, pinnipedAPIClient,
						func(configToUpdate *v1alpha1.CredentialIssuerConfig) {
							configToUpdate.Status.KubeConfigInfo.CertificateAuthorityData = "new-ca-value"
						},
					)
					r.NoError(err)

					expectedGetAction := coretesting.NewGetAction(credentialIssuerConfigGVR, installationNamespace, configName)

					// The first attempted update only includes its own edits.
					firstExpectedUpdatedConfig := existingConfig.DeepCopy()
					firstExpectedUpdatedConfig.Status.KubeConfigInfo.CertificateAuthorityData = "new-ca-value"
					firstExpectedUpdateAction := coretesting.NewUpdateAction(credentialIssuerConfigGVR, installationNamespace, firstExpectedUpdatedConfig)

					// Both the edits made by this update and the edits made by the conflicting update should be included.
					secondExpectedUpdatedConfig := existingConfig.DeepCopy()
					secondExpectedUpdatedConfig.Status.KubeConfigInfo.Server = "some-other-server-value-from-conflicting-update"
					secondExpectedUpdatedConfig.Status.KubeConfigInfo.CertificateAuthorityData = "new-ca-value"
					secondExpectedUpdateAction := coretesting.NewUpdateAction(credentialIssuerConfigGVR, installationNamespace, secondExpectedUpdatedConfig)

					expectedActions := []coretesting.Action{
						expectedGetAction,
						firstExpectedUpdateAction,
						expectedGetAction,
						secondExpectedUpdateAction,
					}
					r.Equal(expectedActions, pinnipedAPIClient.Actions())
				})
			})
		})
	}, spec.Parallel(), spec.Report(report.Terminal{}))
}
