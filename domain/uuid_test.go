package domain

import (
	"reflect"
	"testing"
)

func TestUUIDs_RemoveItem(t *testing.T) {

	one := NewUUID()
	two := NewUUID()
	three := NewUUID()
	type args struct {
		ID UUID
	}
	tests := []struct {
		name string
		us   UUIDs
		args args
		want []UUID
	}{
		{
			name: "remove an item in list",
			us:   UUIDs([]UUID{one, two, three}),
			args: args{
				ID: two,
			},
			want: []UUID{one, three},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.RemoveItem(tt.args.ID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UUIDs.RemoveItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
