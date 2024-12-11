package netutils

import "testing"

func TestCheckSubnet(t *testing.T) {
	type args struct {
		trustedSubnet string
		ip            string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "1# success",
			args: args{
				ip:            "176.14.86.83",
				trustedSubnet: "176.14.64.0/18",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "2# empty trustedSubnet",
			args: args{
				ip:            "176.14.86.83",
				trustedSubnet: "",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "3# wrong ip",
			args: args{
				ip:            "192.14.86.83",
				trustedSubnet: "176.14.64.0/18",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "4# err parsing subnet",
			args: args{
				ip:            "192.14.86.83",
				trustedSubnet: "176.14.64.0/450",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckSubnet(tt.args.trustedSubnet, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}
