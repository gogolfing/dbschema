package refactor

import (
	"reflect"
	"runtime"
	"testing"
)

func testStmtsFunc(t *testing.T, f StmtsFunc, ctx Context, err error, stmts ...Stmt) {
	result, resultErr := f(ctx)
	if !reflect.DeepEqual(result, stmts) || !reflect.DeepEqual(resultErr, err) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("testStmtsFunc: f(ctx) = %v, %v WANT %v, %v @ %v:%v", result, resultErr, stmts, err, file, line)
	}
}
