/*
 * Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"flag"
	"os"
	"reflect"

	"github.com/gardener/gardener/test/testmachinery/extensions/generator"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/gardener/gardener-extension-provider-vsphere/pkg/apis/vsphere/v1alpha1"
)

var (
	infrastructureProviderConfigPath = flag.String("infrastructure-provider-config-filepath", "", "filepath to the provider specific infrastructure config")
	controlplaneProviderConfigPath   = flag.String("controlplane-provider-config-filepath", "", "filepath to the provider specific controlplane config")

	defaultLoadBalancerClass = flag.String("default-load-balancer-class", "default", "default load balancer class")
)

func main() {
	log.SetLogger(zap.New(zap.UseDevMode(false)))
	logger := log.Log.WithName("vsphere-generator")
	flag.Parse()
	if err := validate(); err != nil {
		logger.Error(err, "error validating input flags")
		os.Exit(1)
	}

	infra := v1alpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       reflect.TypeOf(v1alpha1.InfrastructureConfig{}).Name(),
		},
	}

	cp := v1alpha1.ControlPlaneConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       reflect.TypeOf(v1alpha1.ControlPlaneConfig{}).Name(),
		},
		LoadBalancerClasses: []v1alpha1.CPLoadBalancerClass{{Name: *defaultLoadBalancerClass}},
	}

	if err := generator.MarshalAndWriteConfig(*infrastructureProviderConfigPath, infra); err != nil {
		logger.Error(err, "unable to write infrastructure config")
		os.Exit(1)
	}
	if err := generator.MarshalAndWriteConfig(*controlplaneProviderConfigPath, cp); err != nil {
		logger.Error(err, "unable to write infrastructure config")
		os.Exit(1)
	}
	logger.Info("successfully written aws provider configuration", "infra", *infrastructureProviderConfigPath, "controlplane", *controlplaneProviderConfigPath)
}

func validate() error {
	if err := generator.ValidateString(infrastructureProviderConfigPath); err != nil {
		return errors.Wrap(err, "error validating infrastructure provider config path")
	}
	if err := generator.ValidateString(controlplaneProviderConfigPath); err != nil {
		return errors.Wrap(err, "error validating controlplane provider config path")
	}
	if err := generator.ValidateString(defaultLoadBalancerClass); err != nil {
		return errors.Wrap(err, "error validating default load balancer class")
	}
	return nil
}
