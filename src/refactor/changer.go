package refactor

import "log"

//Changer is the interface that wraps the Up and Down methods.
//
//The Up and Down methods should be thought of as opposites that offset eachother.
//For instance, a CreateTable Changer would create a table in Up and drop it in Down.
type Changer interface {
	//Up should return a slice of Stmt(s) that will be executed on the database
	//in order to "up" the state.
	Up(ctx Context) ([]*Stmt, error)

	//Down should return a slice of Stmt(s) that will be executed on the database
	//in order to "down" the state.
	Down(ctx Context) ([]*Stmt, error)
}

func CollectChangersUp(ctx Context, changers ...Changer) ([]*Stmt, error) {
	log.Println("CollectChangersUp", len(changers))
	result := []*Stmt{}
	for _, changer := range changers {
		log.Println("CollectChangersUp", "changer", changer)
		stmts, err := changer.Up(ctx)
		if err != nil {
			return nil, err
		}
		log.Println("CollectChangersUp", stmts)
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
