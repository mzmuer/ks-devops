/*
Copyright 2022 The KubeSphere Authors.

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

package v1alpha3

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"kubesphere.io/devops/pkg/api/devops/v1alpha1"
	"kubesphere.io/devops/pkg/api/devops/v1alpha3"
	fakeclientset "kubesphere.io/devops/pkg/client/clientset/versioned/fake"
	fakedevops "kubesphere.io/devops/pkg/client/devops/fake"
	"kubesphere.io/devops/pkg/client/k8s"
	"kubesphere.io/devops/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAPIsExist(t *testing.T) {
	schema, err := v1alpha1.SchemeBuilder.Register().Build()
	assert.Nil(t, err)

	err = v1.SchemeBuilder.AddToScheme(schema)
	assert.Nil(t, err)

	container := restful.NewContainer()
	AddToContainer(container, fakedevops.NewFakeDevops(nil),
		k8s.NewFakeClientSets(k8sfake.NewSimpleClientset(&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fake", Namespace: "fake",
			},
		}), nil, nil, "", nil,
			fakeclientset.NewSimpleClientset(&v1alpha3.DevOpsProject{
				ObjectMeta: metav1.ObjectMeta{Name: "fake"},
				Status:     v1alpha3.DevOpsProjectStatus{AdminNamespace: "fake"},
			}, &v1alpha3.Pipeline{
				ObjectMeta: metav1.ObjectMeta{Namespace: "fake", Name: "fake"},
			})),
		fake.NewFakeClientWithScheme(schema, &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fake", Namespace: "fake",
			},
		}))

	type args struct {
		method string
		uri    string
	}
	tests := []struct {
		name string
		args args
	}{{
		name: "credential list",
		args: args{
			method: http.MethodGet,
			uri:    "/devops/fake/credentials",
		},
	}, {
		name: "create a credential",
		args: args{
			method: http.MethodPost,
			uri:    "/devops/fake/credentials",
		},
	}, {
		name: "get a credential",
		args: args{
			method: http.MethodGet,
			uri:    "/devops/fake/credentials/fake",
		},
	}, {
		name: "update a credential",
		args: args{
			method: http.MethodPut,
			uri:    "/devops/fake/credentials/fake",
		},
	}, {
		name: "delete a credential",
		args: args{
			method: http.MethodDelete,
			uri:    "/devops/fake/credentials/fake",
		},
	}, {
		name: "get pipeline list",
		args: args{
			method: http.MethodGet,
			uri:    "/devops/fake/pipelines",
		},
	}, {
		name: "create a pipeline",
		args: args{
			method: http.MethodPost,
			uri:    "/devops/fake/pipelines",
		},
	}, {
		name: "get a pipeline",
		args: args{
			method: http.MethodGet,
			uri:    "/devops/fake/pipelines/fake",
		},
	}, {
		name: "build pipeline parameters",
		args: args{
			method: http.MethodGet,
			uri:    "/devops/fake/pipelines/fake/parameters",
		},
	}, {
		name: "update a pipeline",
		args: args{
			method: http.MethodPut,
			uri:    "/devops/fake/pipelines/fake",
		},
	}, {
		name: "delete a pipeline",
		args: args{
			method: http.MethodDelete,
			uri:    "/devops/fake/pipelines/fake",
		},
	}, {
		name: "get devops list",
		args: args{
			method: http.MethodGet,
			uri:    "/workspaces/fake/devops",
		},
	}, {
		name: "create a devops",
		args: args{
			method: http.MethodPost,
			uri:    "/workspaces/fake/devops",
		},
	}, {
		name: "get a devops",
		args: args{
			method: http.MethodGet,
			uri:    "/workspaces/fake/devops/fake",
		},
	}, {
		name: "update a devops",
		args: args{
			method: http.MethodPut,
			uri:    "/workspaces/fake/devops/fake",
		},
	}, {
		name: "delete a devops",
		args: args{
			method: http.MethodDelete,
			uri:    "/workspaces/fake/devops/fake",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpRequest, _ := http.NewRequest(tt.args.method,
				"http://fake.com/kapis/devops.kubesphere.io/v1alpha3"+tt.args.uri, nil)
			httpRequest = httpRequest.WithContext(context.WithValue(context.TODO(), constants.K8SToken, constants.ContextKeyK8SToken("")))

			httpWriter := httptest.NewRecorder()
			container.Dispatch(httpWriter, httpRequest)
			assert.NotEqual(t, 404, httpWriter.Code)
		})
	}
}

func TestGetDevOpsProject(t *testing.T) {
	schema, err := v1alpha1.SchemeBuilder.Register().Build()
	assert.Nil(t, err)
	container := restful.NewContainer()

	AddToContainer(container, fakedevops.NewFakeDevops(nil),
		k8s.NewFakeClientSets(k8sfake.NewSimpleClientset(), nil, nil, "", nil,
			fakeclientset.NewSimpleClientset(&v1alpha3.DevOpsProject{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "fake",
					Name:         "generated-fake",
					Labels: map[string]string{
						constants.WorkspaceLabelKey: "ws",
					},
				},
			})),
		fake.NewFakeClientWithScheme(schema))

	type args struct {
		method string
		uri    string
	}
	tests := []struct {
		name       string
		args       args
		expectCode int
	}{{
		name: "normal case",
		args: args{
			method: http.MethodGet,
			uri:    "/workspaces/ws/devops/generated-fake",
		},
		expectCode: http.StatusOK,
	}, {
		name: "find by a generateName",
		args: args{
			method: http.MethodGet,
			uri:    "/workspaces/ws/devops/fake?generateName=true",
		},
		expectCode: http.StatusOK,
	}, {
		name: "wrong workspace name",
		args: args{
			method: http.MethodGet,
			uri:    "/workspaces/fake/devops/fake?generateName=true",
		},
		expectCode: http.StatusBadRequest,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpRequest, _ := http.NewRequest(tt.args.method,
				"http://fake.com/kapis/devops.kubesphere.io/v1alpha3"+tt.args.uri, nil)
			httpRequest = httpRequest.WithContext(context.WithValue(context.TODO(), constants.K8SToken, constants.ContextKeyK8SToken("")))

			httpWriter := httptest.NewRecorder()
			container.Dispatch(httpWriter, httpRequest)
			assert.Equal(t, tt.expectCode, httpWriter.Code)
		})
	}
}
