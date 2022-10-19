// Code generated by github.com/layotto/protoc-gen-p6. DO NOT EDIT.

// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sms

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"mosn.io/pkg/log"

	sms "mosn.io/layotto/components/sms"
	sms1 "mosn.io/layotto/spec/proto/extension/v1/sms"

	rawGRPC "google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_api "mosn.io/layotto/pkg/grpc"
)

func NewAPI(ac *grpc_api.ApplicationContext) grpc_api.GrpcAPI {
	return &server{
		appId:      ac.AppId,
		components: ac.SmsService,
	}
}

type server struct {
	appId      string
	components map[string]sms.SmsService
}

func (s *server) SendSmsWithTemplate(ctx context.Context, in *sms1.SendSmsWithTemplateRequest) (*sms1.SendSmsWithTemplateResponse, error) {
	// find the component
	comp := s.components[in.ComponentName]
	if comp == nil {
		return nil, invalidArgumentError("SendSmsWithTemplate", grpc_api.ErrComponentNotFound, "sms", in.ComponentName)
	}

	// convert request
	req := &sms.SendSmsWithTemplateRequest{}
	err := copier.CopyWithOption(req, in, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error when converting the request: %s", err.Error())
	}

	// delegate to the component
	resp, err := comp.SendSmsWithTemplate(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// convert response
	out := &sms1.SendSmsWithTemplateResponse{}
	err = copier.CopyWithOption(out, resp, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error when converting the response: %s", err.Error())
	}
	return out, nil
}

func invalidArgumentError(method string, format string, a ...interface{}) error {
	err := status.Errorf(codes.InvalidArgument, format, a...)
	log.DefaultLogger.Errorf(fmt.Sprintf("%s fail: %+v", method, err))
	return err
}

func (s *server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *server) Register(rawGrpcServer *rawGRPC.Server) error {
	sms1.RegisterSmsServiceServer(rawGrpcServer, s)
	return nil
}
