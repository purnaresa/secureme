package main

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_readData(t *testing.T) {
	masterKey = "04076d64bdb6fcf31706eea85ec98431"
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(`SELECT id, name, national_id, create_time_unix FROM user`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "national_id", "create_time_unix"}).
			AddRow("101", "purnaresa", "WKIhX3FryJTlK3FRd2r+Q25QxuCu1YNI6RCvxaYoAfOA8cedcl1ZSFjXdOU=", "1568607939"))

	DB = db
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantUser User
		wantErr  bool
	}{
		{
			name: "read data",
			args: args{
				id: "101",
			},
			wantUser: User{
				ID:             "101",
				Name:           "purnaresa",
				NationalID:     "4444333377779999",
				CreateTimeUnix: "1568607939",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := readData(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("readData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("readData() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func Test_encrypt(t *testing.T) {
	type args struct {
		plaintext string
		key       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy test",
			args: args{
				plaintext: "kingsman",
				key:       "04076d64bdb6fcf31706eea85ec98431"},
		},
		{
			name: "KTP",
			args: args{
				plaintext: "1111222233334444",
				key:       "04076d64bdb6fcf31706eea85ec98431"},
		},
		{
			name: "KTP 2",
			args: args{
				plaintext: "3322135507770044",
				key:       "04076d64bdb6fcf31706eea85ec98431"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// encrypt the plaintext
			ciphertext, err := encrypt(tt.args.plaintext, tt.args.key)
			if err != nil {
				t.Errorf("encrypt() error = %v", err)
				return
			}
			t.Logf("ciphertext = %s", ciphertext)
			//

			// decrypt the ciphertext from previous encrypt function
			plaintext, err := decrypt(ciphertext, tt.args.key)
			if err != nil {
				t.Errorf("encrypt() error = %v", err)
				return
			}
			t.Logf("plaintext = %s", plaintext)
			//

			// compare the initial plaintext with output of previous decrypt function
			if plaintext != tt.args.plaintext {
				t.Errorf("plaintext = %v, want %v", plaintext, tt.args.plaintext)
			}
			//
		})
	}
}
