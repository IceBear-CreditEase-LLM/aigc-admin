package modelevaluate

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/encode"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/helpers/page"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

var validate = validator.New()

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Evaluate": ems,
	})
	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/create", kithttp.NewServer(
		eps.CreateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/cancel/{uuid}", kithttp.NewServer(
		eps.CancelEndpoint,
		decodeCancelRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	r.Handle("/delete/{uuid}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)
	r.Handle("/fivegraph", kithttp.NewServer(
		eps.FiveGraphEndpoint,
		decodeFiveGraphRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/finish/{jobId}", kithttp.NewServer(
		eps.FinishEndpoint,
		decodeFinishRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req listRequest
	req.ModelId, _ = strconv.Atoi(r.URL.Query().Get("modelId"))
	req.Status = r.URL.Query().Get("status")
	req.EvalTargetType = r.URL.Query().Get("evalTargetType")
	req.Page = page.GetPage(r)
	req.PageSize = page.GetPageSize(r)

	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeCancelRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req cancelRequest
	vars := mux.Vars(r)
	req.Uuid = vars["uuid"]
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req deleteRequest
	vars := mux.Vars(r)
	req.Uuid = vars["uuid"]
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeFiveGraphRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req fiveGraphRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeFinishRequest(_ context.Context, r *http.Request) (interface{}, error) {
	bodyBytes, _ := io.ReadAll(r.Body)
	fmt.Println("decodeFinishRequest", string(bodyBytes))

	var req finishRequest
	vars := mux.Vars(r)
	req.JobId = vars["jobId"]
	req.Result = string(bodyBytes)
	return req, nil
}
