# Gator - Blog Aggregator CLI

A command-line RSS feed aggregator built in Go. Gator lets you follow blogs,
aggregate their posts on a schedule, and browse posts from feeds you follow.

## Prerequisites

Before running gator, make sure you have installed:
- [Go](https://go.dev/doc/install) (version 1.21 or later)
- [PostgreSQL](https://www.postgresql.org/download/) 

## Installation

 You can then install `gator` with:

```bash
GOPROXY=direct go install github.com/thenitverse/gator@latest
```

## Config

Create a `.gatorconfig.json` file in your home directory with the following
structure:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username`, `password`, and `gator` with your own database
connection details. The `current_user_name` field is managed automatically
by the CLI once you register or log in - you don't need to set it manually.

## Usage

Register a new user (this also logs you in):

```bash
gator register <name>
```

Add a feed (this also follows it automatically):

```bash
gator addfeed <name> <url>
```

Start the aggregator, fetching feeds on a set interval:

```bash
gator agg 30s
```

Leave this running in its own terminal - it will continuously poll for new
posts based on the interval you give it (e.g. `30s`, `1m`, `5m`).

Browse posts from feeds you follow:

```bash
gator browse [limit]
```

`limit` is optional and defaults to a small number of recent posts if
omitted.

### Other commands

- `gator login <name>` - Log in as a user that already exists
- `gator users` - List all registered users
- `gator feeds` - List all feeds currently in the database
- `gator follow <url>` - Follow a feed that already exists in the database
- `gator unfollow <url>` - Unfollow a feed you're currently following
- `gator following` - List the feeds the current user follows
- `gator reset` - Delete all users (useful for testing)

## How it works

1. Users register/log in via the CLI, which tracks the "current user" in
   the config file.
2. Feeds are stored in Postgres, and any user can follow any feed via a
   many-to-many relationship table.
3. `gator agg` runs a loop that fetches the least-recently-fetched feed,
   parses its RSS XML, and saves new posts to the database.
4. `gator browse` queries posts belonging to feeds the current user
   follows, ordered by publish date.

## What I Learned

Building this project helped me practice:

- **PostgreSQL fundamentals** - designing tables, primary/foreign keys,
  and many-to-many relationships (users <-> feeds)
- **Schema migrations** - using Goose to version and apply database schema
  changes
- **Type-safe SQL** - using sqlc to generate Go code from raw SQL queries
- **RSS/XML parsing** - fetching and unmarshaling XML feeds into Go structs
- **Middleware patterns** - wrapping command handlers to require an
  authenticated/logged-in user before running certain commands
## Notes

If `gator browse` shows 0 posts even though `agg` is successfully
fetching feeds, check whether you're actually following the feed with
your current user. `browse` only shows posts from feeds the logged-in
user follows - it won't show posts from feeds that exist in the database
but are not followed. I ran into this myself: `agg` was fetching and
storing posts just fine, but `browse` kept returning 0 because I had
never run `gator follow <url>` for that feed under my user.If you're running my gator and browse shows 0 posts, check this first.

