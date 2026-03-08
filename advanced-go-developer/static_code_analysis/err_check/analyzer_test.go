package errCheck_test

import (
	"testing"

	"github.com/oleshko-g/errCheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Analyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), errCheck.Analyzer, "./...")
}
