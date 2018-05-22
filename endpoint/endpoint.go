package endpoint

import "context"

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

type Middleware func(Endpoint) Endpoint

type Failer interface {
	Failed() error
}
