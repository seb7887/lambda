package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"seb7887/lambda/pkg/logger"
	"testing"
)

func TestService_Do(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		tests  = []struct {
			name    string
			give    *Request
			want    *Response
			wantErr assert.ErrorAssertionFunc
		}{
			{
				name: "should do something",
				give: &Request{
					UserID:      "iao",
					OperationID: "123",
				},
				want: &Response{
					Status: "ok",
				},
				wantErr: assert.NoError,
			},
			{
				name: "should return error if user ID is 'x'",
				give: &Request{
					UserID:      "x",
					OperationID: "123",
				},
				want:    nil,
				wantErr: assert.Error,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New()
			r, err := svc.Do(ctx, tt.give)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, r)
		})
	}
}
