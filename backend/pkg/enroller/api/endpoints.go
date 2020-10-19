package api

import (
	"context"
	"enroller/pkg/enroller/crypto"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PostCSREndpoint            endpoint.Endpoint
	GetPendingCSRsEndpoint     endpoint.Endpoint
	GetPendingCSRDBEndpoint    endpoint.Endpoint
	GetPendingCSRFileEndpoint  endpoint.Endpoint
	PutChangeCSRStatusEndpoint endpoint.Endpoint
	DeleteCSREndpoint          endpoint.Endpoint
	GetCRTEndpoint             endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostCSREndpoint:            MakePostCSREndpoint(s),
		GetPendingCSRsEndpoint:     MakeGetPendingCSRsEndpoint(s),
		GetPendingCSRDBEndpoint:    MakeGetPendingCSRDBEndpoint(s),
		GetPendingCSRFileEndpoint:  MakeGetPendingCSRFileEndpoint(s),
		PutChangeCSRStatusEndpoint: MakePutChangeCSRStatusEndpoint(s),
		DeleteCSREndpoint:          MakeDeleteCSREndpoint(s),
		GetCRTEndpoint:             MakeGetCTREndpoint(s),
	}
}

func MakePostCSREndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postCSRRequest)
		csr, e := s.PostCSR(ctx, req.data)
		return postCSRResponse{CSR: csr, Err: e}, nil
	}
}

func MakeGetPendingCSRsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(getPendingCSRsRequest)
		csrs := s.GetPendingCSRs(ctx)
		return getPendingCSRsResponse{CSRs: csrs}, nil
	}
}

func MakeGetPendingCSRDBEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPendingCSRRequest)
		csr, err := s.GetPendingCSRDB(ctx, req.ID)
		return getPendingCSRDBResponse{CSR: csr, Err: err}, nil
	}
}

func MakeGetPendingCSRFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPendingCSRRequest)
		data, err := s.GetPendingCSRFile(ctx, req.ID)
		return getPendingCSRFileResponse{Data: data, Err: err}, nil
	}
}

func MakePutChangeCSRStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putChangeCSRStatusRequest)
		err = s.PutChangeCSRStatus(ctx, req.CSR, req.ID)
		return putChangeCSRsResponse{Err: err}, nil
	}
}

func MakeDeleteCSREndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteCSRRequest)
		err = s.DeleteCSR(ctx, req.ID)
		return deleteCSRResponse{Err: err}, nil
	}
}

func MakeGetCTREndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCRTRequest)
		data, err := s.GetCRT(ctx, req.ID)
		return getCRTResponse{Data: data, Err: err}, nil
	}
}

type getCRTRequest struct {
	ID int
}

type getCRTResponse struct {
	Data []byte
	Err  error
}

type postCSRRequest struct {
	data []byte
}

type postCSRResponse struct {
	CSR crypto.CSR `json:"csr,omitempty"`
	Err error      `json:"err,omitempty"`
}

func (r postCSRResponse) error() error { return r.Err }

type getPendingCSRsRequest struct{}

type getPendingCSRsResponse struct {
	CSRs crypto.CSRs `json:"CSRs,omitempty"`
}

type getPendingCSRRequest struct {
	ID int
}

type getPendingCSRDBResponse struct {
	CSR crypto.CSR `json:"CSR,omitempty"`
	Err error      `json:"err,omitempty"`
}

func (r getPendingCSRDBResponse) error() error { return r.Err }

type getPendingCSRFileResponse struct {
	Data []byte
	Err  error
}

type putChangeCSRStatusRequest struct {
	CSR crypto.CSR
	ID  int
}

type putChangeCSRsResponse struct {
	Err error
}

func (r putChangeCSRsResponse) error() error { return r.Err }

type deleteCSRRequest struct {
	ID int
}

type deleteCSRResponse struct {
	Err error
}

func (r deleteCSRResponse) error() error { return r.Err }
