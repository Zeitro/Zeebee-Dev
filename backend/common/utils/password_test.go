package utils

import "testing"

func TestCheckPassword(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				password: "12345678",
				hash:     "$2a$10$y3Md5HA3VdhPxW9djW46hejcHj9kpvBwnINemZAHqJiSZcZ4Ne3bq",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPassword(tt.args.password, tt.args.hash); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				password: "12345678",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Error("HashPassword() returned an empty hash")
				return
			}
			if !CheckPassword(tt.args.password, got) {
				t.Errorf("HashPassword() returned a hash that does not verify the password")
			}
		})
	}
}
