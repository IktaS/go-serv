package serv

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateServServer(t *testing.T) {
	tests := []struct {
		Name     string
		Serv     *Serv
		Setup    func(*testing.T) *os.File
		Teardown func(*testing.T, *os.File)
		expected string
		WantErr  bool
	}{
		{
			Name: "Test Service with reference only",
			Serv: &Serv{
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
			Setup: func(t *testing.T) *os.File {
				f, err := os.Create("./testfile.go")
				if err != nil {
					assert.Error(t, err)
				}
				return f
			},
			Teardown: func(t *testing.T, f *os.File) {
				f.Close()
				os.Remove(f.Name())
			},
			expected: "//CODE GENERATED BY GO-SERV\nfunc TestService (TestMessage1) TestMessage2\n",
			WantErr:  false,
		},
		{
			Name: "Test Service with two parameter",
			Serv: &Serv{
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
			Setup: func(t *testing.T) *os.File {
				f, err := os.Create("./testfile.go")
				if err != nil {
					assert.Error(t, err)
				}
				return f
			},
			Teardown: func(t *testing.T, f *os.File) {
				f.Close()
				os.Remove(f.Name())
			},
			expected: "//CODE GENERATED BY GO-SERV\nfunc TestService (TestMessage1, TestMessage2) TestMessage2\n",
			WantErr:  false,
		},
		{
			Name: "Test Message with scalar only",
			Serv: &Serv{
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
			Setup: func(t *testing.T) *os.File {
				f, err := os.Create("./testfile.go")
				if err != nil {
					assert.Error(t, err)
				}
				return f
			},
			Teardown: func(t *testing.T, f *os.File) {
				f.Close()
				os.Remove(f.Name())
			},
			expected: "//CODE GENERATED BY GO-SERV\ntype TestMessage struct {\n    TestString string;\n}\n",
			WantErr:  false,
		},
		{
			Name: "Test Message with two field",
			Serv: &Serv{
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
										Name: "TestString2",
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
			Setup: func(t *testing.T) *os.File {
				f, err := os.Create("./testfile.go")
				if err != nil {
					assert.Error(t, err)
				}
				return f
			},
			Teardown: func(t *testing.T, f *os.File) {
				f.Close()
				os.Remove(f.Name())
			},
			expected: "//CODE GENERATED BY GO-SERV\ntype TestMessage struct {\n    TestString string;\n    TestString2 string;\n}\n",
			WantErr:  false,
		},
		{
			Name: "Test Message and Service",
			Serv: &Serv{
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
			Setup: func(t *testing.T) *os.File {
				f, err := os.Create("./testfile.go")
				if err != nil {
					assert.Error(t, err)
				}
				return f
			},
			Teardown: func(t *testing.T, f *os.File) {
				f.Close()
				os.Remove(f.Name())
			},
			expected: "//CODE GENERATED BY GO-SERV\ntype TestMessage struct {\n    TestString string;\n}\nfunc TestService (TestMessage1) TestMessage2\n",
			WantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			f := tt.Setup(t)
			err := GenerateServServer(f, tt.Serv)
			if err != nil && !tt.WantErr {
				tt.Teardown(t, f)
				log.Fatal(err)
			}
			if tt.WantErr {
				assert.Error(t, err)
			}
			res, err := ioutil.ReadFile(f.Name())
			if err != nil && !tt.WantErr {
				tt.Teardown(t, f)
				log.Fatal(err)
			}
			if tt.WantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expected, string(res))
			tt.Teardown(t, f)
		})
	}
}
