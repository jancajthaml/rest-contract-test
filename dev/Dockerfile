# Copyright (c) 2017-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
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

# ---------------------------------------------------------------------------- #

FROM golang:1.10-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl

# INFO temporary gcc support for RAML parsers
RUN apk add --no-cache gcc musl-dev

RUN curl -L https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 -o /usr/bin/dep && \
    chmod +x /usr/bin/dep

RUN go get -u \
        \
        github.com/golang/lint/golint \
        github.com/fzipp/gocyclo \
        github.com/client9/misspell/cmd/misspell

# ---------------------------------------------------------------------------- #
