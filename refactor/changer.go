package refactor

import "github.com/gogolfing/dbschema/refactor/dto"

type Changer interface {
	Up(ctx Context) ([]*Stmt, error)
	Down(ctx Context) ([]*Stmt, error)

	DTO() (dto.Changer, error)
}

func CollectChangersUp(ctx Context, changers ...Changer) ([]*Stmt, error) {
	result := []*Stmt{}
	for _, changer := range changers {
		stmts, err := changer.Up(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, stmts...)
	}
	return result, nil
}

func CollectChangersDown(ctx Context, changers ...Changer) ([]*Stmt, error) {
	result := []*Stmt{}
	for _, changer := range changers {
		stmts, err := changer.Down(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, stmts...)
	}
	return result, nil
}
