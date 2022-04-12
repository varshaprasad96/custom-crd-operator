#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

export execDir="/Users/vnarsing/go/src/github.com/code-generator"

"${execDir}"/generate-groups.sh "deepcopy,client,informer,lister" \
  github.com/varshaprasad96/custom-crd-operator/pkg/generated github.com/varshaprasad96/custom-crd-operator/pkg/apis \
  example.com:v1alpha1 \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${execDir}"/hack/boilerplate.go.txt



"${execDir}"/generate-groups.sh "client,informer,lister" \
  github.com/example-inc/lib-go-plugin-operator/generated github.com/example-inc/lib-go-plugin-operator \
  api:v1alpha1 \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${execDir}"/hack/boilerplate.go.txt


"${execDir}"/generate-groups.sh "client,informer,lister" \
  github.com/example-inc/lib-go-plugin-operator/api/generated github.com/example-inc/lib-go-plugin-operator/api \
  cache.my.domain:v1alpha1 \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${execDir}"/hack/boilerplate.go.txt