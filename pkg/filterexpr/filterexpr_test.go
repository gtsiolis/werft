package filterexpr_test

import (
	"reflect"
	"testing"

	v1 "github.com/32leaves/werft/pkg/api/v1"
	"github.com/32leaves/werft/pkg/filterexpr"
	"github.com/alecthomas/repr"
)

func TestValidBasics(t *testing.T) {
	tests := []struct {
		Input  string
		Result *v1.FilterExpression
		Error  string
	}{
		{"foo==bar", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "foo", Value: "bar", Operation: v1.FilterOp_OP_EQUALS, Negate: false}}}, ""},
		{"foo!==bar", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "foo", Value: "bar", Operation: v1.FilterOp_OP_EQUALS, Negate: true}}}, ""},
		{"foo~=bar", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "foo", Value: "bar", Operation: v1.FilterOp_OP_CONTAINS, Negate: false}}}, ""},
		{"foo!~=bar", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "foo", Value: "bar", Operation: v1.FilterOp_OP_CONTAINS, Negate: true}}}, ""},
		{"foo|=bar", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "foo", Value: "bar", Operation: v1.FilterOp_OP_STARTS_WITH, Negate: false}}}, ""},
		{"foo!|=bar", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "foo", Value: "bar", Operation: v1.FilterOp_OP_STARTS_WITH, Negate: true}}}, ""},
		{"foo=|bar", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "foo", Value: "bar", Operation: v1.FilterOp_OP_ENDS_WITH, Negate: false}}}, ""},
		{"foo!=|bar", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "foo", Value: "bar", Operation: v1.FilterOp_OP_ENDS_WITH, Negate: true}}}, ""},
		{"success==true", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "success", Value: "1", Operation: v1.FilterOp_OP_EQUALS, Negate: false}}}, ""},
		{"success==false", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "success", Value: "0", Operation: v1.FilterOp_OP_EQUALS, Negate: false}}}, ""},
		{"success!==true", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "success", Value: "1", Operation: v1.FilterOp_OP_EQUALS, Negate: true}}}, ""},
		{"success!==false", &v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "success", Value: "0", Operation: v1.FilterOp_OP_EQUALS, Negate: true}}}, ""},
		{"foo", nil, filterexpr.ErrMissingOp.Error()},
		{"phase==blabla", nil, "invalid phase: blabla"},
	}

	for _, test := range tests {
		res, err := filterexpr.Parse([]string{test.Input})
		if err != nil {
			if err.Error() != test.Error {
				t.Errorf("%s: %v != %v", test.Input, err, test.Error)
			}
			continue
		}

		if len(res) != 1 {
			t.Errorf("%s: resulted in NOT exactly one expression, but %v", test.Input, repr.String(res))
			continue
		}
		if !reflect.DeepEqual(res[0], test.Result) {
			t.Errorf("%s: expected %s but got %s", test.Input, repr.String(test.Result), repr.String(res[0]))
			continue
		}
	}
}

func TestMatchesFilter(t *testing.T) {
	md := &v1.JobMetadata{
		Owner:      "foo",
		Repository: &v1.Repository{},
	}
	tests := []struct {
		Job     *v1.JobStatus
		Expr    []*v1.FilterExpression
		Matches bool
	}{
		{
			&v1.JobStatus{Metadata: md, Phase: v1.JobPhase_PHASE_DONE},
			[]*v1.FilterExpression{&v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "phase", Value: "done", Operation: v1.FilterOp_OP_EQUALS}}}},
			true,
		},
	}

	for idx, test := range tests {
		if filterexpr.MatchesFilter(test.Job, test.Expr) != test.Matches {
			t.Errorf("test %d failed", idx)
		}
	}
}