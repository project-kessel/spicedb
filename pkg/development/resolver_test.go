package development

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authzed/spicedb/pkg/schemadsl/compiler"
	"github.com/authzed/spicedb/pkg/schemadsl/input"
)

func TestResolver(t *testing.T) {
	testSource := input.Source("test")

	tcs := []struct {
		name              string
		schema            string
		line              int
		column            int
		expectedReference *SchemaReference
	}{
		{
			name: "basic relation",
			schema: `definition user {}

			definition resource {
				relation viewer: user
				permission view = viewer
			}
			`,
			line:   4,
			column: 24,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 4, ColumnPosition: 24},
				Text:                     "viewer",
				ReferenceType:            ReferenceTypeRelation,
				ReferenceMarkdown:        "relation viewer",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 3, ColumnPosition: 4},
				TargetSourceCode:         "relation viewer: user\n",
				TargetNamePositionOffset: 9,
			},
		},
		{
			name: "basic permission",
			schema: `definition user {}

			definition resource {
				relation viewer: user
				relation editor: user
				permission edit = editor
				permission view = viewer + edit
			}
			`,
			line:   6,
			column: 33,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 6, ColumnPosition: 33},
				Text:                     "edit",
				ReferenceType:            ReferenceTypePermission,
				ReferenceMarkdown:        "permission edit",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 5, ColumnPosition: 4},
				TargetSourceCode:         "permission edit = editor\n",
				TargetNamePositionOffset: 11,
			},
		},
		{
			name: "basic type",
			schema: `definition user {}

			definition resource {
				relation viewer: user
				permission view = viewer
			}
			`,
			line:   3,
			column: 24,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 3, ColumnPosition: 24},
				Text:                     "user",
				ReferenceType:            ReferenceTypeDefinition,
				ReferenceMarkdown:        "definition user",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 0, ColumnPosition: 0},
				TargetSourceCode:         "definition user {}",
				TargetNamePositionOffset: 11,
			},
		},
		{
			name: "subject relation type",
			schema: `definition user {}

			definition group {
				relation member: user
			}

			definition resource {
				relation viewer: group#member
				permission view = viewer
			}
			`,
			line:   7,
			column: 24,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 7, ColumnPosition: 24},
				Text:                     "group",
				ReferenceType:            ReferenceTypeDefinition,
				ReferenceMarkdown:        "definition group",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 2, ColumnPosition: 3},
				TargetSourceCode:         "definition group {\n\t// ...\n}",
				TargetNamePositionOffset: 11,
			},
		},
		{
			name: "subject relation relation",
			schema: `definition user {}

			definition group {
				relation member: user
			}

			definition resource {
				relation viewer: group#member
				permission view = viewer
			}
			`,
			line:   7,
			column: 32,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 7, ColumnPosition: 32},
				Text:                     "member",
				ReferenceType:            ReferenceTypeRelation,
				ReferenceMarkdown:        "relation member",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 3, ColumnPosition: 4},
				TargetSourceCode:         "relation member: user\n",
				TargetNamePositionOffset: 9,
			},
		},
		{
			name: "filled in type",
			schema: `definition user {}

			definition resource {
				relation viewer: user | resource
				permission view = viewer
			}
			`,
			line:   3,
			column: 29,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 3, ColumnPosition: 29},
				Text:                     "resource",
				ReferenceType:            ReferenceTypeDefinition,
				ReferenceMarkdown:        "definition resource",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 2, ColumnPosition: 3},
				TargetSourceCode:         "definition resource {\n\t// ...\n}",
				TargetNamePositionOffset: 11,
			},
		},
		{
			name: "caveat reference",
			schema: `definition user {}

			caveat somecaveat(someparam int) {
				someparam < 42
			}

			definition resource {
				relation viewer: user with somecaveat
				permission view = viewer
			}
			`,
			line:   7,
			column: 35,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 7, ColumnPosition: 35},
				Text:                     "somecaveat",
				ReferenceType:            ReferenceTypeCaveat,
				ReferenceMarkdown:        "caveat somecaveat",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 2, ColumnPosition: 3},
				TargetSourceCode:         "caveat somecaveat(someparam int) {\n\t// ...\n}",
				TargetNamePositionOffset: 7,
			},
		},
		{
			name: "arrow relation reference",
			schema: `definition user {}

			definition resource {
				relation viewer: user
				permission view = viewer->foo
			}
			`,
			line:   4,
			column: 23,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 4, ColumnPosition: 23},
				Text:                     "viewer",
				ReferenceType:            ReferenceTypeRelation,
				ReferenceMarkdown:        "relation viewer",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 3, ColumnPosition: 4},
				TargetSourceCode:         "relation viewer: user\n",
				TargetNamePositionOffset: 9,
			},
		},
		{
			name: "arrow permission reference",
			schema: `definition user {}

			definition resource {
				relation viewer: user
				permission view = viewer->viewer
			}
			`,
			line:              4,
			column:            31,
			expectedReference: nil,
		},
		{
			name: "caveat parameter reference",
			schema: `definition user {}

			caveat somecaveat(someparam int) {
				someparam < 42
			}

			definition resource {
				relation viewer: user with somecaveat
				permission view = viewer
			}
			`,
			line:   3,
			column: 6,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 3, ColumnPosition: 6},
				Text:                     "someparam",
				ReferenceType:            ReferenceTypeCaveatParameter,
				ReferenceMarkdown:        "someparam int",
				TargetSource:             &testSource,
				TargetSourceCode:         "someparam int",
				TargetNamePositionOffset: 0,
			},
		},
		{
			name: "longer test",
			schema: `definition user {}

definition document {
	relation viewer: user
	relation editor: user
	relation third: user
	permission another = viewer
	permission view = third + editor + another
}`,
			line:   7,
			column: 19,
			expectedReference: &SchemaReference{
				Source:                   input.Source("test"),
				Position:                 input.Position{LineNumber: 7, ColumnPosition: 19},
				Text:                     "third",
				ReferenceType:            ReferenceTypeRelation,
				ReferenceMarkdown:        "relation third",
				TargetSource:             &testSource,
				TargetPosition:           &input.Position{LineNumber: 5, ColumnPosition: 1},
				TargetSourceCode:         "relation third: user\n",
				TargetNamePositionOffset: 9,
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			compiled, err := compiler.Compile(compiler.InputSchema{
				Source:       input.Source("test"),
				SchemaString: tc.schema,
			}, compiler.AllowUnprefixedObjectType())
			require.NoError(t, err)

			resolver, err := NewResolver(compiled)
			require.NoError(t, err)

			ref, err := resolver.ReferenceAtPosition(input.Source("test"), input.Position{
				LineNumber:     tc.line,
				ColumnPosition: tc.column,
			})

			require.NoError(t, err)
			require.Equal(t, tc.expectedReference, ref)
		})
	}
}
