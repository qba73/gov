package gov_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/gov"
)

func TestParse_ReturnsDependencyOnValidInput(t *testing.T) {
	t.Parallel()

	tt := []struct {
		input string
		want  gov.Dependency
	}{
		{
			input: `dep	google.golang.org/genproto/googleapis/api	v0.0.0-20250303144028-a0af3efb3deb	h1:p31xT4yrYrSM/G4Sn2+TNUkVhFCbG9y8itM2S6Th950=`,
			want: gov.Dependency{
				Name:    "google.golang.org/genproto/googleapis/api",
				Version: "v0.0.0-20250303144028-a0af3efb3deb",
				Digest:  "h1:p31xT4yrYrSM/G4Sn2+TNUkVhFCbG9y8itM2S6Th950=",
			},
		},
		{
			input: `dep	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go	v1.36.4-20250130201111-63bb56e20495.1	h1:4erM3WLgEG/HIBrpBDmRbs1puhd7p0z7kNXDuhHthwM=`,
			want: gov.Dependency{
				Name:    "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go",
				Version: "v1.36.4-20250130201111-63bb56e20495.1",
				Digest:  "h1:4erM3WLgEG/HIBrpBDmRbs1puhd7p0z7kNXDuhHthwM=",
			},
		},
	}

	for _, tc := range tt {
		got, err := gov.Parse(tc.input)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(tc.want, got) {
			t.Error(cmp.Diff(tc.want, got))
		}
	}
}

func TestParse_ReturnsErrorOnInvalidInput(t *testing.T) {
	t.Parallel()

	tt := []struct {
		input string
	}{
		{
			input: `dep	google.golang.org/genproto/googleapis/apiv0.0.0-20250303144028-a0af3efb3deb	h1:p31xT4yrYrSM/G4Sn2+TNUkVhFCbG9y8itM2S6Th950=`,
		},
		{
			input: `buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go	v1.36.4-20250130201111-63bb56e20495.1	h1:4erM3WLgEG/HIBrpBDmRbs1puhd7p0z7kNXDuhHthwM=`,
		},
		{
			input: `dep	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3debh1:p31xT4yrYrSM/G4Sn2+TNUkVhFCbG9y8itM2S6Th950=`,
		},
		{
			input: `dep	google.golang.org/genproto/googleapis/apiv0.0.0-20250303144028-a0af3efb3debh1:p31xT4yrYrSM/G4Sn2+TNUkVhFCbG9y8itM2S6Th950=`,
		},
		{
			input: `dep`,
		},
	}

	for _, tc := range tt {
		if _, err := gov.Parse(tc.input); err == nil {
			t.Errorf("want err on invalid input line: %q, got nil", tc.input)
		}
	}
}

func TestDependencies_ReturnsParsedDependenciesOnValidInput(t *testing.T) {
	t.Parallel()

	in, err := os.Open("testdata/gov_short.txt")
	if err != nil {
		t.Fatal(err)
	}

	p, err := gov.NewParser(
		gov.WithInput(in),
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := p.Dependencies()
	if err != nil {
		t.Fatal(err)
	}

	want := []gov.Dependency{
		{
			Name:    "cel.dev/expr",
			Version: "v0.19.1",
			Digest:  "h1:NciYrtDRIR0lNCnH1LFJegdjspNx9fI59O7TWcua/W4=",
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestToJSON_ReturnsJSONRepresentationOfParsedDependencies(t *testing.T) {
	t.Parallel()

	in, err := os.Open("testdata/gov_short.txt")
	if err != nil {
		t.Fatal(err)
	}

	p, err := gov.NewParser(
		gov.WithInput(in),
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := p.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	want := `[{"name":"cel.dev/expr","version":"v0.19.1","digest":"h1:NciYrtDRIR0lNCnH1LFJegdjspNx9fI59O7TWcua/W4="}]`
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestWithInputFromArgs_SetsInputFilePath(t *testing.T) {
	t.Parallel()

	args := []string{"testdata/gov_short.txt"}
	_, err := gov.NewParser(
		gov.WithInputFromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithInputFromArgs_IgnoresEmptyArgs(t *testing.T) {
	t.Parallel()

	in, err := os.Open("testdata/gov_short.txt")
	if err != nil {
		t.Fatal(err)
	}

	p, err := gov.NewParser(
		gov.WithInput(in),
		gov.WithInputFromArgs([]string{}),
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := p.Dependencies()
	if err != nil {
		t.Fatal(err)
	}

	want := []gov.Dependency{
		{
			Name:    "cel.dev/expr",
			Version: "v0.19.1",
			Digest:  "h1:NciYrtDRIR0lNCnH1LFJegdjspNx9fI59O7TWcua/W4=",
		},
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
