package domain

import (
	"orgStruct/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_organizationLogic_CreateDivision(t *testing.T) {
	type fields struct {
		db func() *mocks.MockDb
	}
	type args struct {
		name     string
		ParentID *int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test_organizationLogic_CreateDivision_ErrorsNameNotСorrect",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					return DbMock
				},
			},
			args: args{
				name:     "",
				ParentID: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &organizationLogic{
				db: tt.fields.db(),
			}
			if err := logic.CreateDivision(tt.args.name, tt.args.ParentID); (err != nil) != tt.wantErr {
				t.Errorf("CreateDivision() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
