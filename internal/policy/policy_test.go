package policy

import (
	"reflect"
	"testing"

	"github.com/keel-hq/keel/types"
)

func Test_getPolicyFromLabels(t *testing.T) {
	type args struct {
		labels map[string]string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name:  "policy all",
			args:  args{labels: map[string]string{types.KeelPolicyLabel: "all"}},
			want1: true,
			want:  "all",
		},
		{
			name:  "policy minor",
			args:  args{labels: map[string]string{types.KeelPolicyLabel: "minor"}},
			want1: true,
			want:  "minor",
		},
		{
			name:  "legacy policy minor",
			args:  args{labels: map[string]string{"keel.observer/policy": "minor"}},
			want1: true,
			want:  "minor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getPolicyFromLabels(tt.args.labels)
			if got != tt.want {
				t.Errorf("getPolicyFromLabels() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getPolicyFromLabels() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func mustParseGlob(g string) *GlobPolicy {
	glb, err := NewGlobPolicy(g)
	if err != nil {
		panic(err)
	}
	return glb
}

func TestGetPolicy(t *testing.T) {
	type args struct {
		policyName string
		options    *Options
	}
	tests := []struct {
		name string
		args args
		want Policy
	}{
		{
			name: "patch",
			args: args{policyName: "patch", options: &Options{}},
			want: NewSemverPolicy(SemverPolicyTypePatch),
		},
		{
			name: "glob:foo-*",
			args: args{policyName: "glob:foo-*", options: &Options{}},
			want: mustParseGlob("glob:foo-*"),
		},
		{
			name: "force match",
			args: args{policyName: "force", options: &Options{MatchTag: true}},
			want: NewForcePolicy(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPolicy(tt.args.policyName, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
