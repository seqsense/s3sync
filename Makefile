# Copyright 2019 SEQSENSE, Inc.
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

.PHONY: test
test:
	go test . -v

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out .
	go tool cover -html=cover.out -o report.html

.PHONY: s3
s3:
	docker compose -f ./test/docker-compose.e2e.yml up

.PHONY: s3-bg
s3-bg:
	docker compose -f ./test/docker-compose.e2e.yml up -d --wait

.PHONY: s3-down
s3-down:
	docker compose -f ./test/docker-compose.e2e.yml down
