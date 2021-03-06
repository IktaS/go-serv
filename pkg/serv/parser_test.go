package serv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCombination(t *testing.T) {
	tests := []struct {
		name     string
		src      []byte
		expected *Gserv
		wantErr  bool
	}{
		{
			name: "Service with message as parameter",
			src: []byte(`
				message TestMessage{
					string TestString;
				};
				def inbound TestService(TestMessage) : string;
			`),
			expected: &Gserv{
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
							Name:     "TestService",
							Inbound:  true,
							Outbound: false,
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
			name: "Service with message as parameter",
			src: []byte(`
				message TestMessage{
					string TestString;
				};
				def outbound TestService(TestMessage) : string;
			`),
			expected: &Gserv{
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
							Name:     "TestService",
							Inbound:  false,
							Outbound: true,
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
				def inbound TestService(TestMessage,string) : string;
			`),
			expected: &Gserv{
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
							Name:     "TestService",
							Inbound:  true,
							Outbound: false,
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
				def outbound TestService(string) : TestMessage;
			`),
			expected: &Gserv{
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
							Name:     "TestService",
							Inbound:  false,
							Outbound: true,
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
				def inbound TestService(TestMessage) : TestMessage;
			`),
			expected: &Gserv{
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
							Name:     "TestService",
							Inbound:  true,
							Outbound: false,
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
		{
			name: "Service with same name",
			src: []byte(`
				def inbound TestService(string);
				def inbound TestService(TestMessage) : TestMessage;
				message TestMessage{
					string TestString;
				};
			`),
			wantErr: true,
		},
		{
			name: "Service with nonexistent message",
			src: []byte(`
				def inbound TestService(string) : TestMessage;
			`),
			wantErr: true,
		},
	}

	parser, err := NewServParser()
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parser.Parse(tt.src)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

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
		expected *Gserv
		wantErr  bool
	}{
		{
			name: " Inbound service with no parameter and returns",
			src: []byte(`
				def inbound TestService();
			`),
			expected: &Gserv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name:     "TestService",
							Inbound:  true,
							Outbound: false,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Service with no parameter but returns",
			src: []byte(`
				def outbound TestService():int32;
			`),
			expected: &Gserv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name:     "TestService",
							Inbound:  false,
							Outbound: true,
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
				def inbound TestService(string):int32;
			`),
			expected: &Gserv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name:     "TestService",
							Inbound:  true,
							Outbound: false,
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
				def outbound TestService(string):TestMessage2;
			`),
			expected: &Gserv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name:     "TestService",
							Inbound:  false,
							Outbound: true,
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
				def inbound TestService(string,int32):TestMessage2;
			`),
			expected: &Gserv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name:     "TestService",
							Inbound:  true,
							Outbound: false,
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
				def outbound TestService(TestMessage1):TestMessage2;
			`),
			expected: &Gserv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name:     "TestService",
							Inbound:  false,
							Outbound: true,
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
				def inbound TestService(TestMessage1,TestMessage2):TestMessage2;
			`),
			expected: &Gserv{
				Definitions: []*Definition{
					{
						Service: &Service{
							Name:     "TestService",
							Inbound:  true,
							Outbound: false,
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
		expected *Gserv
		wantErr  bool
	}{
		{
			name: "One Message, one property",
			src: []byte(`
				message TestMessage{
					string TestString;
				};
			`),
			expected: &Gserv{
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
			expected: &Gserv{
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
			expected: &Gserv{
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
