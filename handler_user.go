package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
	// "errors")
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("no username was provided")
	}
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: <%s>", cmd.Args[0])
	return nil

}
func handlerRegister(s *state, cmds command) error {
	if len(cmds.Args) != 1 {
		return fmt.Errorf("no username was provided")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmds.Args[0],
	})
	if err != nil {
		return err

	}

	err = s.cfg.SetUser(cmds.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: <%s>", cmds.Args[0])
	fmt.Println(user.Name)
	return nil

}
