package helpers

import (
	"context"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type Case[TActual any] struct {
	t     *testing.T
	ctx   context.Context
	given []func(context.Context) error
	when  func(context.Context) (TActual, error)
	then  []func(context.Context, TActual) error
}

func NewCase[TActual any](ctx context.Context, t *testing.T) Case[TActual] {
	return Case[TActual]{
		t:     t,
		ctx:   ctx,
		given: make([]func(context.Context) error, 0),
		then:  make([]func(context.Context, TActual) error, 0),
	}
}

func (c Case[TActual]) Given(fn func(context.Context) error) Case[TActual] {
	c.given = append(c.given, fn)
	return c
}

func (c Case[TActual]) When(fn func(context.Context) (TActual, error)) Case[TActual] {
	c.when = fn
	return c
}

func (c Case[TActual]) Then(fn func(context.Context, TActual) error) Case[TActual] {
	c.then = append(c.then, fn)
	return c
}

func (c Case[TActual]) Run() {
	for _, g := range c.given {
		if err := g(c.ctx); err != nil {
			c.t.Fatalf("failed %s with: %s", getFunctionName(g), err.Error())
		}
	}

	actual, err := c.when(c.ctx)
	if err != nil {
		c.t.Fatalf("failed %s with: %s", getFunctionName(c), err.Error())
	}

	for _, t := range c.then {
		if err := t(c.ctx, actual); err != nil {
			c.t.Fatalf("failed %s with: %s", getFunctionName(t), err.Error())
		}
	}
}

func getFunctionName(fn any) string {
	str := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	return strings.ReplaceAll(str[len(str)-1], "-fm", "")
}
