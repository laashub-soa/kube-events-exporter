/*
Copyright 2020 Red Hat, Inc. All rights reserved.

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

package framework

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// GetDeployment gets the given deployment.
func (f *Framework) GetDeployment(ns, name string) (*appsv1.Deployment, error) {
	deployment, err := f.KubeClient.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "get deployment %s", name)
	}
	return deployment, nil
}

// CreateDeployment creates the given deployment.
func (f *Framework) CreateDeployment(deployment *appsv1.Deployment, ns string) (*appsv1.Deployment, error) {
	deployment.Namespace = ns
	deployment, err := f.KubeClient.AppsV1().Deployments(ns).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "create deployment %s", deployment.Name)
	}
	return deployment, nil
}

// MakeDeployment creates a deployment object from yaml manifest.
func MakeDeployment(manifestPath string) (*appsv1.Deployment, error) {
	manifest, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, errors.Wrapf(err, "read deployment manifest %s", manifestPath)
	}

	deployment := appsv1.Deployment{}
	err = yaml.Unmarshal(manifest, &deployment)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal deployment manifest %s", manifestPath)
	}

	return &deployment, nil
}

// DeleteDeployment deletes the given deployment.
func (f *Framework) DeleteDeployment(ns, name string) error {
	return f.KubeClient.AppsV1().Deployments(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// WaitUntilDeploymentReady waits until given deployment is ready.
func (f *Framework) WaitUntilDeploymentReady(ns, name string) error {
	err := wait.Poll(5*time.Second, 2*time.Minute, func() (bool, error) {
		deployment, err := f.GetDeployment(ns, name)
		if err != nil {
			return false, nil
		}
		return deployment.Status.ReadyReplicas == deployment.Status.Replicas, nil
	})
	if err != nil {
		return errors.Wrapf(err, "deployment %s pods are not ready", name)
	}

	return nil
}
