# Gator - Blog Aggregator

Gator is a blog aggregator application built using Go and Postgres. It allows users to fetch and scrape blog posts from various RSS feeds and store them in a database for later consumption. The application provides a CLI to interact with the aggregator and perform various actions like adding feeds, scraping posts, and more.

## Prerequisites

Before you begin, you'll need to have the following installed:

1. **PostgreSQL**: The application uses a PostgreSQL database to store feeds and posts.
2. **Go**: Gator is built in Go, so you’ll need Go installed to build and run the application.

### Install PostgreSQL
Follow the instructions on the [PostgreSQL website](https://www.postgresql.org/download/) to install PostgreSQL for your platform.

### Install Go
Follow the instructions on the [Go website](https://golang.org/doc/install) to install Go for your platform.

## Installation

### 1. Install the `gator` CLI

To install the `gator` CLI, use the following Go command. This will download the CLI tool and install it globally.

```bash
go install github.com/tombraggg/blog-aggregator@latest
```

### 2. Create a .gatorconfig.json File:

Create a file called ".gatorconfig.json" in your home directory. It should contain something like this:

```json
{
    "db_url":"postgres://username:@localhost:5432/gator?sslmode=disable"
}
```

## Usage

Create a new user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start the aggregator:

```bash
gator agg 30s
```

View the posts:

```bash
gator browse [limit]
```

There are a few other commands you'll need as well:

- `gator login <name>` - Log in as a user that already exists
- `gator users` - List all users
- `gator feeds` - List all feeds
- `gator follow <url>` - Follow a feed that already exists in the database
- `gator unfollow <url>` - Unfollow a feed that already exists in the database