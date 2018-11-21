package domain

import (
	"reflect"
	"testing"
)

func TestVertices_Join(t *testing.T) {
	aID := NewUUID()
	bID := NewUUID()
	type args struct {
		b Vertices
	}
	tests := []struct {
		name string
		a    Vertices
		args args
		want Vertices
	}{
		{
			name: "join 1 vs 1",
			a: Vertices{
				Vertex{
					Node: &Node{
						ID: aID,
					},
				},
			},
			args: args{
				b: Vertices{
					Vertex{
						Node: &Node{
							ID: bID,
						},
					},
				},
			},
			want: Vertices{
				Vertex{
					Node: &Node{
						ID: aID,
					},
				},
				Vertex{
					Node: &Node{
						ID: bID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Join(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vertices.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}
