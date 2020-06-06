// Copyright 2018 The Range Core Authors
// Copyright 2018 The go-ethereum Authors
// This file is part of the Range Core library.
//
// The Range Core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Range Core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Range Core library. If not, see <http://www.gnu.org/licenses/>.

package http

import (
	"context"

	"range/core/gen3/swarm/api"
	"range/core/gen3/swarm/sctx"
)

type uriKey struct{}

func GetRUID(ctx context.Context) string {
	v, ok := ctx.Value(sctx.HTTPRequestIDKey{}).(string)
	if ok {
		return v
	}
	return "xxxxxxxx"
}

func SetRUID(ctx context.Context, ruid string) context.Context {
	return context.WithValue(ctx, sctx.HTTPRequestIDKey{}, ruid)
}

func GetURI(ctx context.Context) *api.URI {
	v, ok := ctx.Value(uriKey{}).(*api.URI)
	if ok {
		return v
	}
	return nil
}

func SetURI(ctx context.Context, uri *api.URI) context.Context {
	return context.WithValue(ctx, uriKey{}, uri)
}
