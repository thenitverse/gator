package main

import (
	"database/sql"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	programState := &state{cfg: &cfg, db: dbQueries}
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	//cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	if len(os.Args) < 2 {
		log.Fatal("not enough arguments provided")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{Name: cmdName, Args: cmdArgs}
	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
