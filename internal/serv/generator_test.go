package serv

import (
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
			WantErr: false,
		},
		{
			Name: "Test Service with reference only",
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
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			f := tt.Setup(t)
			err := GenerateServServer(f, tt.Serv)
			println()
			if tt.WantErr {
				assert.Error(t, err)
			}
			tt.Teardown(t, f)
		})
	}
}
