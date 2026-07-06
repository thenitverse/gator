package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("...: %w", err)
	}
	fmt.Println("...")
	return nil
}
func handlerListUsers(s *state, cmd command) error {
	results, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, item := range results {
		if item.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", item.Name)
		} else {
			fmt.Printf("* %v\n", item.Name)
		}
	}
	return nil
}
