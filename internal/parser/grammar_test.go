package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ozansz/semantix/pkg/ptrutils"
)

func TestFileParser(t *testing.T) {
	tests := []struct {
		desc string
		file string
		File *File
	}{
		{
			desc: "facts",
			file: "../../examples/v0/0x01-facts.sxql",
			File: factsFile(),
		},
		{
			desc: "queries",
			file: "../../examples/v0/0x02-queries.sxql",
			File: queriesFile(),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			parser := New()
			f, err := parser.ParseFile(tc.file)
			if err != nil {
				t.Errorf("failed to parse file: %v", err)
				return
			}
			if diff := cmp.Diff(tc.File, f); diff != "" {
				t.Errorf("unexpected file (-want +got):\n%s", diff)
			}
		})
	}
}

func factsFile() *File {
	return &File{
		Expressions: []*Expression{
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ozan"),
					Predicate: "is",
					Object:    SubjectObject{Value: "Person"},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ufuk"),
					Predicate: "is",
					Object:    SubjectObject{Value: "Person"},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ozan"),
					Predicate: "knows",
					Object:    SubjectObject{Value: "CS"},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ufuk"),
					Predicate: "knows",
					Object:    SubjectObject{Value: "CS"},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ufuk"),
					Predicate: "length",
					Object:    NumberObject{Value: 123.456},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ozan"),
					Predicate: "age",
					Object:    NumberObject{Value: 24},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ozan"),
					Predicate: "name",
					Object:    StringObject{Value: "Ozan Sazak!!!"},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ufuk"),
					Predicate: "knows",
					ObjectFact: &Fact{
						Subject:   ptrutils.Ptr("Ozan"),
						Predicate: "knows",
						Object:    SubjectObject{Value: "CS"},
					},
				},
			},
			{
				Fact: &Fact{
					SubjectFact: &Fact{
						Subject:   ptrutils.Ptr("Ozan"),
						Predicate: "knows",
						Object:    SubjectObject{Value: "CS"},
					},
					Predicate: "approvedBy",
					Object:    SubjectObject{Value: "METU"},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ezgi"),
					Predicate: "thinks",
					ObjectFact: &Fact{
						Subject:   ptrutils.Ptr("Ozan"),
						Predicate: "is",
						Object:    SubjectObject{Value: "Person"},
					},
				},
			},
			{
				Fact: &Fact{
					Subject:   ptrutils.Ptr("Ezgi"),
					Predicate: "notKnows",
					ObjectFact: &Fact{
						Subject:   ptrutils.Ptr("Ufuk"),
						Predicate: "knows",
						ObjectFact: &Fact{
							Subject:   ptrutils.Ptr("Ozan"),
							Predicate: "knows",
							Object:    SubjectObject{Value: "CS"},
						},
					},
				},
			},
		},
	}
}

func queriesFile() *File {
	return &File{
		Expressions: []*Expression{
			{
				Query: &Query{
					SubjectVar: ptrutils.Ptr("?x"),
					Predicate:  ptrutils.Ptr("is"),
					Object:     SubjectObject{Value: "Person"},
					IDInFile:   "Q1",
					Kind:       QueryKindSimple,
				},
			},
			{
				Query: &Query{
					SubjectVar: ptrutils.Ptr("?x"),
					Predicate:  ptrutils.Ptr("knows"),
					ObjectQuery: &Query{
						SubjectVar: ptrutils.Ptr("?y"),
						Predicate:  ptrutils.Ptr("knows"),
						Object:     SubjectObject{Value: "CS"},
					},
					IDInFile: "CQ1",
					Kind:     QueryKindCompound,
				},
			},
			{
				Query: &Query{
					Subject:   ptrutils.Ptr("Ozan"),
					Predicate: ptrutils.Ptr("is"),
					ObjectVar: ptrutils.Ptr("?x"),
					IDInFile:  "Q2",
					Kind:      QueryKindSimple,
				},
			},
			{
				Query: &Query{
					SubjectVar: ptrutils.Ptr("?x"),
					Predicate:  ptrutils.Ptr("knows"),
					ObjectVar:  ptrutils.Ptr("?y"),
					IDInFile:   "Q3",
					Kind:       QueryKindSimple,
				},
			},
			{
				Query: &Query{
					SubjectVar:   ptrutils.Ptr("?x"),
					PredicateVar: ptrutils.Ptr("?y"),
					ObjectVar:    ptrutils.Ptr("?z"),
					IDInFile:     "Q4",
					Kind:         QueryKindSimple,
				},
			},
			{
				Query: &Query{
					SubjectVar: ptrutils.Ptr("?x"),
					Predicate:  ptrutils.Ptr("knows"),
					ObjectVar:  ptrutils.Ptr("!y"),
					LinkedQuery: &Query{
						SubjectVar: ptrutils.Ptr("!y"),
						Predicate:  ptrutils.Ptr("subtopicOf"),
						Object:     SubjectObject{Value: "Science"},
					},
					IDInFile: "LQ1",
					Kind:     QueryKindLinked,
				},
			},
			{
				Query: &Query{
					SubjectVar: ptrutils.Ptr("?x"),
					Predicate:  ptrutils.Ptr("knows"),
					ObjectQuery: &Query{
						SubjectVar: ptrutils.Ptr("!y"),
						Predicate:  ptrutils.Ptr("knows"),
						Object:     SubjectObject{Value: "CS"},
					},
					IDInFile: "CQ2",
					Kind:     QueryKindCompound,
				},
			},
			{
				Query: &Query{
					SubjectVar: ptrutils.Ptr("!x"),
					Predicate:  ptrutils.Ptr("is"),
					Object:     SubjectObject{Value: "Person"},
					LinkedQuery: &Query{
						SubjectVar: ptrutils.Ptr("?y"),
						Predicate:  ptrutils.Ptr("is"),
						Object:     SubjectObject{Value: "Person"},
						LinkedQuery: &Query{
							SubjectVar: ptrutils.Ptr("?y"),
							Predicate:  ptrutils.Ptr("knows"),
							ObjectQuery: &Query{
								SubjectVar: ptrutils.Ptr("!x"),
								Predicate:  ptrutils.Ptr("knows"),
								Object:     SubjectObject{Value: "CS"},
							},
						},
					},
					IDInFile: "LCQ1",
					Kind:     QueryKindLinkedCompound,
				},
			},
		},
	}
}
