package serv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCombination(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		expected *Serv
		wantErr  bool
	}{
		{
			name: "Service with message as parameter",
			src: []byte(`
				message TestMessage{
					string TestString;
				};
				service TestService(TestMessage) : string;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Message: &Message{
							Name: "TestMessage",
							Definitions: []*MessageDefinition{
								{
									Field: &Field{
										Name: "TestString",
										Type: &Type{
											Scalar: String,
										},
									},
								},
							},
						},
					},
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Reference: "TestMessage",
								},
							},
							Response: &Type{
								Scalar: String,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with message and primitive as parameter",
			src: []byte(`
				message TestMessage{
					string TestString;
				};
				service TestService(TestMessage,string) : string;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Message: &Message{
							Name: "TestMessage",
							Definitions: []*MessageDefinition{
								{
									Field: &Field{
										Name: "TestString",
										Type: &Type{
											Scalar: String,
										},
									},
								},
							},
						},
					},
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Reference: "TestMessage",
								},
								{
									Scalar: String,
								},
							},
							Response: &Type{
								Scalar: String,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with message as return",
			src: []byte(`
				message TestMessage{
					string TestString;
				};
				service TestService(string) : TestMessage;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Message: &Message{
							Name: "TestMessage",
							Definitions: []*MessageDefinition{
								{
									Field: &Field{
										Name: "TestString",
										Type: &Type{
											Scalar: String,
										},
									},
								},
							},
						},
					},
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Scalar: String,
								},
							},
							Response: &Type{
								Reference: "TestMessage",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with message as parameter and return",
			src: []byte(`
				message TestMessage{
					string TestString;
				};
				service TestService(TestMessage) : TestMessage;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Message: &Message{
							Name: "TestMessage",
							Definitions: []*MessageDefinition{
								{
									Field: &Field{
										Name: "TestString",
										Type: &Type{
											Scalar: String,
										},
									},
								},
							},
						},
					},
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Reference: "TestMessage",
								},
							},
							Response: &Type{
								Reference: "TestMessage",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	parser, err := NewServParser()
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parser.Parse(tt.src)
			assert.Equal(t, tt.expected, res)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestParseService(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		expected *Serv
		wantErr  bool
	}{
		{
			name: "Service with no parameter and returns",
			src: []byte(`
				service TestService();
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name: "TestService",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with no parameter but returns",
			src: []byte(`
				service TestService():int32;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name: "TestService",
							Response: &Type{
								Scalar: Int32,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with return primitive",
			src: []byte(`
				service TestService(string):int32;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Scalar: String,
								},
							},
							Response: &Type{
								Scalar: Int32,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with parameter primitive",
			src: []byte(`
				service TestService(string):TestMessage2;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Scalar: String,
								},
							},
							Response: &Type{
								Reference: "TestMessage2",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with two parameter primitives",
			src: []byte(`
				service TestService(string,int32):TestMessage2;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Scalar: String,
								},
								{
									Scalar: Int32,
								},
							},
							Response: &Type{
								Reference: "TestMessage2",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Only define service with message parameter and message return",
			src: []byte(`
				service TestService(TestMessage1):TestMessage2;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Reference: "TestMessage1",
								},
							},
							Response: &Type{
								Reference: "TestMessage2",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with two parameter message",
			src: []byte(`
				service TestService(TestMessage1,TestMessage2):TestMessage2;
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name: "TestService",
							Request: []*Type{
								{
									Reference: "TestMessage1",
								},
								{
									Reference: "TestMessage2",
								},
							},
							Response: &Type{
								Reference: "TestMessage2",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	parser, err := NewServParser()
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parser.Parse(tt.src)
			assert.Equal(t, tt.expected, res)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		expected *Serv
		wantErr  bool
	}{
		{
			name: "One Message, one property",
			src: []byte(`
				message TestMessage{
					string TestString;
				};
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Message: &Message{
							Name: "TestMessage",
							Definitions: []*MessageDefinition{
								{
									Field: &Field{
										Name: "TestString",
										Type: &Type{
											Scalar: String,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "One Message, two property",
			src: []byte(`
				message TestMessage{
					string TestString;
					int32 TestInt;
				};
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Message: &Message{
							Name: "TestMessage",
							Definitions: []*MessageDefinition{
								{
									Field: &Field{
										Name: "TestString",
										Type: &Type{
											Scalar: String,
										},
									},
								},
								{
									Field: &Field{
										Name: "TestInt",
										Type: &Type{
											Scalar: Int32,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple Messages",
			src: []byte(`
				message TestMessage{
					string TestString;
					int32 TestInt;
				};
				message TestMessage2{
					string TestString;
					int32 TestInt;
				};
			`),
			expected: &Serv{
				Definitions: []*Definition{
					{
						Message: &Message{
							Name: "TestMessage",
							Definitions: []*MessageDefinition{
								{
									Field: &Field{
										Name: "TestString",
										Type: &Type{
											Scalar: String,
										},
									},
								},
								{
									Field: &Field{
										Name: "TestInt",
										Type: &Type{
											Scalar: Int32,
										},
									},
								},
							},
						},
					},
					{
						Message: &Message{
							Name: "TestMessage2",
							Definitions: []*MessageDefinition{
								{
									Field: &Field{
										Name: "TestString",
										Type: &Type{
											Scalar: String,
										},
									},
								},
								{
									Field: &Field{
										Name: "TestInt",
										Type: &Type{
											Scalar: Int32,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	parser, err := NewServParser()
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parser.Parse(tt.src)
			assert.Equal(t, tt.expected, res)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}
