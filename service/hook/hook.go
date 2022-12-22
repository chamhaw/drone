// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hook

import (
	"context"
	"github.com/drone/drone/core"
	"github.com/drone/drone/service/httputils"
	"github.com/drone/go-scm/scm"
)

// New returns a new HookService.
func New(client *scm.Client, addr string, renew core.Renewer) core.HookService {
	return &service{client: client, addr: addr, renew: renew}
}

type service struct {
	renew  core.Renewer
	client *scm.Client
	addr   string
}

func (s *service) Create(ctx context.Context, user *core.User, repo *core.Repository) error {
	ctx, client, err := httputils.PrepareHttpClient(ctx, user, s.renew)
	if err != nil {
		return err
	}
	if client != nil {
		s.client.Client = client
	}
	hook := &scm.HookInput{
		Name:   "drone",
		Target: s.addr + "/hook",
		Secret: repo.Signer,
		Events: scm.HookEvents{
			Branch:      true,
			Deployment:  true,
			PullRequest: true,
			Push:        true,
			Tag:         true,
		},
	}
	return replaceHook(ctx, s.client, repo.Slug, hook)
}

func (s *service) Delete(ctx context.Context, user *core.User, repo *core.Repository) error {
	ctx, client, err := httputils.PrepareHttpClient(ctx, user, s.renew)
	if err != nil {
		return err
	}
	if client != nil {
		s.client.Client = client
	}
	return deleteHook(ctx, s.client, repo.Slug, s.addr)
}
