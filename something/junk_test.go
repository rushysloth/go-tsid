package something

import (
	"fmt"
	"testing"

	"github.com/rushysloth/go-tsid"
)

func TestSomething(t *testing.T) {
	var id *tsid.Tsid
	id = tsid.Fast()
	fmt.Println(id.ToString())

	var factory *tsid.TsidFactory
	factory, _ = tsid.TsidFactoryBuilder().
		WithNode(1).
		Build()

	tsid, _ := factory.Generate()
	fmt.Println(tsid.ToString())
}
