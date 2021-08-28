package main

import "testing"

func TestTrimQuantityToLotSize(t *testing.T) {
	type args struct {
		quantity string
		lotSize  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Quantity longer than lotsize",
			args: args{
				quantity: "123.456789",
				lotSize:  "0.0001",
			},
			want: "123.4567",
		},
		{
			name: "Quantity longer than lotsize II",
			args: args{
				quantity: "123.456789",
				lotSize:  "0.1",
			},
			want: "123.4",
		},
		{
			name: "Quantity shorter than lotsize",
			args: args{
				quantity: "123.4",
				lotSize:  "0.0001",
			},
			want: "123.4",
		}, {
			name: "Quantity without decimal places",
			args: args{
				quantity: "123",
				lotSize:  "0.0001",
			},
			want: "123",
		}, {
			name: "lot size 1.00",
			args: args{
				quantity: "123.123",
				lotSize:  "1.00",
			},
			want: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimQuantityToLotSize(tt.args.quantity, tt.args.lotSize); got != tt.want {
				t.Errorf("TrimQuantityToLotSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
