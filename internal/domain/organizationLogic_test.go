package domain

import (
	"errors"
	"orgStruct/internal/mocks"
	"orgStruct/internal/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
				name:     " ",
				ParentID: nil,
			},
			wantErr: true,
		},
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
				//имя больше 200 символов
				name:     "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn",
				ParentID: nil,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateDivision_ErrorsSelect",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					DbMock.EXPECT().
						SelectDepartmentWhereParentID(gomock.Any()).
						Return(nil, errors.New("test error"))
					return DbMock
				},
			},
			args: args{

				name:     " test_name ",
				ParentID: nil,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateDivision_ErrorsNameRepeat",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					DbMock.EXPECT().
						SelectDepartmentWhereParentID(gomock.Any()).
						Return([]repository.Department{{Name: "test_name", ParentID: nil}}, nil)
					return DbMock
				},
			},
			args: args{
				name:     " test_name ",
				ParentID: nil,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateDivision_Success",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					DbMock.EXPECT().
						SelectDepartmentWhereParentID(gomock.Any()).
						Return([]repository.Department{{Name: "test_name2", ParentID: nil}}, nil)
					DbMock.EXPECT().InsertDepartment(gomock.Any()).DoAndReturn(func(dep repository.Department) {
						assert.Equal(t, "test_name", dep.Name)
						assert.Nil(t, dep.ParentID)
					})
					return DbMock
				},
			},
			args: args{

				name:     " test_name ",
				ParentID: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &organizationLogic{
				db: tt.fields.db(),
			}
			_, err := logic.CreateDepartmen(tt.args.name, tt.args.ParentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDivision() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_organizationLogic_CreateEmployee(t *testing.T) {
	type fields struct {
		db func() *mocks.MockDb
	}
	type args struct {
		fullName    string
		position    string
		idDepartmen int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test_organizationLogic_CreateEmployeev_ErrorsNameNotСorrect",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					return DbMock
				},
			},
			args: args{
				fullName:    "",
				position:    "testPosition",
				idDepartmen: 1,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateEmployeev_ErrorsNameNotСorrect",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					return DbMock
				},
			},
			args: args{
				fullName:    "test_name ",
				position:    "",
				idDepartmen: 1,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateEmployeev_ErrorsNameNotСorrect",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					return DbMock
				},
			},
			args: args{
				fullName:    "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn",
				position:    "testPosition",
				idDepartmen: 1,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateEmployeev_ErrorsNameNotСorrect",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					return DbMock
				},
			},
			args: args{
				fullName:    "test_name ",
				position:    "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn",
				idDepartmen: 1,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateEmployeev_ErrorsNameNotСorrect",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					return DbMock
				},
			},
			args: args{
				fullName:    "  ",
				position:    "testPosition",
				idDepartmen: 1,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateEmployeev_ErrorsNameNotСorrect",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					return DbMock
				},
			},
			args: args{
				fullName:    "test_name ",
				position:    "  ",
				idDepartmen: 1,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateEmployeev_ErrorsDB",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					DbMock.EXPECT().
						SelectDepartmentById(gomock.Any()).
						Return(repository.Department{}, errors.New("select error"))
					return DbMock
				},
			},
			args: args{
				fullName:    "test_name ",
				position:    "testPosition  ",
				idDepartmen: 1,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateEmployeev_ErrorsDepartmentNotExist",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					DbMock.EXPECT().
						SelectDepartmentById(gomock.Any()).
						Return(repository.Department{}, nil)
					return DbMock
				},
			},
			args: args{
				fullName:    "test_name ",
				position:    "testPosition  ",
				idDepartmen: 1,
			},
			wantErr: true,
		},
		{
			name: "Test_organizationLogic_CreateEmployeev_Success",
			fields: fields{
				db: func() *mocks.MockDb {
					g := gomock.NewController(t)
					DbMock := mocks.NewMockDb(g)
					DbMock.EXPECT().
						SelectDepartmentById(gomock.Any()).
						Return(repository.Department{Name: "testDepartment"}, nil)
					DbMock.EXPECT().
						InsertEmployee(gomock.Any())
					return DbMock
				},
			},
			args: args{
				fullName:    "test_name ",
				position:    "testPosition  ",
				idDepartmen: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &organizationLogic{
				db: tt.fields.db(),
			}
			if err := logic.CreateEmployee(tt.args.fullName, tt.args.position, tt.args.idDepartmen); (err != nil) != tt.wantErr {
				t.Errorf("CreateEmployee() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
