package mockstore

import (
	"reflect"
	"testing"

	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/stretchr/testify/mock"
)

func TestMockStore_GetLink(t *testing.T) {
	type args struct {
		arg *model.StoreRecord
	}
	tests := []struct {
		name    string
		ms      *MockStore
		args    args
		want    *model.StoreRecord
		wantErr bool
	}{
		{
			name: "1# success",
			ms:   NewMockStore(),
			args: args{
				arg: &model.StoreRecord{
					UUID: "1",
				},
			},
			want: &model.StoreRecord{
				UUID: "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.ms.On("GetLink", mock.Anything).Return(tt.args.arg, nil)
			got, err := tt.ms.GetLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockStore.GetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockStore.GetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
