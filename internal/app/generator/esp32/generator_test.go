package serv

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/IktaS/go-serv/pkg/serv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateServServer(t *testing.T) {
	tests := []struct {
		Name     string
		Serv     *serv.Gserv
		Setup    func(*testing.T) *os.File
		Teardown func(*testing.T, *os.File)
		expected string
		WantErr  bool
	}{}
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
