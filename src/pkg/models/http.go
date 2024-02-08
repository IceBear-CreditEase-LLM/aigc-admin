package models

import (
	"context"
	"encoding/json"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/encode"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/helpers/page"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var validate = validator.New()

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, dmw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Model": ems,
	})

	r := mux.NewRouter()
	r.Handle("/models", kithttp.NewServer(
		eps.CreateModelEndpoint,
		decodeCreateModelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models", kithttp.NewServer(
		eps.ListModelsEndpoint,
		decodeListModelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/models/{id}", kithttp.NewServer(
		eps.UpdateModelEndpoint,
		decodeUpdateModelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPut)
	r.Handle("/models/{id}", kithttp.NewServer(
		eps.DeleteModelEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	r.Handle("/models/eval", kithttp.NewServer(
		eps.ListEvalEndpoint,
		decodeListEvalRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/models/{id}", kithttp.NewServer(
		eps.GetModelEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/models/{id}/deploy", kithttp.NewServer(
		eps.DeployModelEndpoint,
		decodeModelDeployRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models/{id}/undeploy", kithttp.NewServer(
		eps.UndeployModelEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models/eval", kithttp.NewServer(
		eps.CreateEvalEndpoint,
		decodeCreateEvalRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models/eval/{id}/cancel", kithttp.NewServer(
		eps.CancelEvalEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models/eval/{id}", kithttp.NewServer(
		eps.DeleteEvalEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	return r
}

func decodeCreateModelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	if req.IsPrivate && req.Parameters == 0 {
		return nil, encode.InvalidParams.Wrap(errors.New("私有模型，参数量不能为空"))
	}
	return req, nil
}

func decodeUpdateModelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("modelId is empty"))
	}
	mid, err := strconv.Atoi(id)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req UpdateModelRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	req.Id = uint(mid)
	if err = validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeListModelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListModelRequest{
		Page:     page.GetPage(r),
		PageSize: page.GetPageSize(r),
	}
	req.ModelName = r.URL.Query().Get("modelName")
	req.ProviderName = r.URL.Query().Get("providerName")
	enabled := r.URL.Query().Get("enabled")
	if enabled != "" {
		b, err := strconv.ParseBool(enabled)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(err)
		}
		req.Enabled = &b
	}
	isPrivate := r.URL.Query().Get("isPrivate")
	if isPrivate != "" {
		b, err := strconv.ParseBool(isPrivate)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(err)
		}
		req.IsPrivate = &b
	}
	isFineTuning := r.URL.Query().Get("isFineTuning")
	if isFineTuning != "" {
		b, err := strconv.ParseBool(isFineTuning)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(err)
		}
		req.IsFineTuning = &b
	}
	return req, nil
}

func decodeIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("id is empty"))
	}
	mid, err := strconv.Atoi(id)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req IdRequest
	req.Id = uint(mid)
	return req, nil
}

func decodeCreateEvalRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateEvalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeListEvalRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListEvalRequest{
		Page:     page.GetPage(r),
		PageSize: page.GetPageSize(r),
	}
	req.ModelName = r.URL.Query().Get("modelName")
	req.MetricName = r.URL.Query().Get("metricName")
	req.Status = r.URL.Query().Get("status")
	req.DatasetType = r.URL.Query().Get("datasetType")
	return req, nil
}

func decodeModelDeployRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("id is empty"))
	}
	mid, err := strconv.Atoi(id)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req ModelDeployRequest
	req.Id = uint(mid)
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if req.Quantization == "" {
		req.Quantization = "float16"
	}
	return req, nil
}
