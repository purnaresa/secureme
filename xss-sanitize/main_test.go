package main

import "testing"

func Test_sanitize(t *testing.T) {
	type args struct {
		untrustedPayload string
	}
	tests := []struct {
		name               string
		args               args
		wantTrustedPayload string
	}{
		{
			name: "test 1",
			args: args{
				untrustedPayload: `here <script>alert("a")</script>`,
			},
			wantTrustedPayload: `here `,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTrustedPayload := sanitize(tt.args.untrustedPayload); gotTrustedPayload != tt.wantTrustedPayload {
				t.Errorf("sanitize() = %v, want %v", gotTrustedPayload, tt.wantTrustedPayload)
			}
		})
	}
}
