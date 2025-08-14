# Gator - RSS Feed Aggregator

Gator is a command-line RSS feed aggregator that allows you to follow your favorite blogs and news sites, automatically fetch new posts, and browse them directly in your terminal.

## Features
- **Multi-user support**: multiple users can register and log-in as their own user
- **Feed management**: Add, follow, and unfollow RSS feeds
- **Automatic aggregation**: Continuously fetch new posts from your followed feeds
- **Terminal browsing**: View posts directly in your command line with clean formatting
- **Duplicate handling**: Automatically prevents duplicate posts
- **Flexible scheduling**: Configure how often feeds are fetched

## Prerequisites
Before installing Gator, make sure you have the following installed:

## Required Software
1. **Go 1.19 or later**
   - Download from [golang.org](https://golang.org)
   - Verify installation: `go version`
2. **PostgreSQL**
   - Download from [postgresql.org](https://postgresql.org)
   - Or install via package manager:
   ```
   # macOS with Homebrew
   brew install postgresql

   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib

   # Windows
   # Download installer from postgresql.org
   ```
## Installation
### Install Gator CLI
```go install github.com/quduss/gator@latest```

After installation, the `gator` binary will be available in your `$GOPATH/bin` directory (usually `~/go/bin`). Make sure this directory is in your system's PATH.
### Database Setup
1. **Create a PostgreSQL database:**
```
createdb gator
```
## Configuration
### Create Config File
Create a configuration file at `~/.gatorconfig.json`:
```
{
  "db_url": "postgres://username:password@localhost:port/gator?sslmode=disable",
  "current_user_name": "your_username"
}
```
**Configuration Options:**
- `db_url`: PostgreSQL connection string
  - Format: `postgres://user:password@host:port/database?sslmode=disable`
- `current_user_name`: name of the logged-in user
### Database Migration
The application will automatically run database migrations on first use, creating the necessary tables:
- `users` - User accounts
- `feeds` - RSS feed sources
- `feed_follows` - User feed subscriptions
- `posts` - Individual RSS posts
## Usage
### Getting Started
1. **Register as a user**:
```
gator register <username>
```
2. **Add an RSS feed**:
```
gator addfeed TechCrunch https://techcrunch.com/feed/
```
3. **Follow the feed:**:
```
gator follow https://techcrunch.com/feed/
```
4. **Start the aggregator (in a separate terminal):**:
```
gator agg 1m
```
5. **Browse your posts:**:
```
gator browse
```
### Available Commands
**User Management**
```
# Register a new user
gator register <username>

# Login as existing user  
gator login <username>

# Reset database (removes all data)
 gator reset
```
**Feed Management**
```
# Add a new RSS feed
gator addfeed <feed name> <feed url>

# Follow a feed (start receiving its posts)
gator follow <url>

# Unfollow a feed
gator unfollow <url>

# List all available feeds
gator feeds

# List feeds you're following
gator following
```
**Post Management**
```
# Browse recent posts from your followed feeds
gator browse [limit]        # Default limit: 2
gator browse 10            # Show 10 recent posts

# Start continuous feed aggregation
gator agg <duration>       # e.g., "1m", "30s", "5m", "1h"
```
### Example Workflow
```
# Initial setup
gator register johndoe
gator addfeed "Hacker News" https://news.ycombinator.com/rss
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
gator follow https://news.ycombinator.com/rss
gator follow https://blog.boot.dev/index.xml

# Start aggregation (leave running in background)
gator agg 2m

# In another terminal, browse posts
gator browse 5
```
### Recommended RSS Feeds
Here are some popular RSS feeds to get you started:
- **Tech News**:
  - TechCrunch: `https://techcrunch.com/feed/`
  - Hacker News: `https://news.ycombinator.com/rss`
  - Ars Technica: `https://feeds.arstechnica.com/arstechnica/index`
- **Development**:
  - Boot.dev Blog: `https://blog.boot.dev/index.xml`
  - Go Blog: `https://blog.golang.org/feeds/posts/default`
  - Go Blog: `https://blog.golang.org/feeds/posts/default`
- **General News**:
  - BBC News: `http://feeds.bbci.co.uk/news/rss.xml`
  - Reuters: `https://feeds.reuters.com/reuters/topNews`
### Architecture
Gator is built with:
- **Go**: Core application logic
- **PostgreSQL**: Data persistence
- **SQLC**: Type-safe database queries
- **Goose**: Database migrations
The application uses a middleware pattern for authentication and follows clean architecture principles with separate concerns for database access, RSS parsing, and command handling.
### Development
### Running from Source
```
# Clone the repository
git clone https://github.com/quduss/gator.git
cd gator

# Install dependencies
go mod download

# Run in development mode
go run . <command>

# Build binary
go build -o gator

# Install locally
go install
```
### Database Schema
The application uses four main tables:
- `users`: User accounts and authentication
- `feeds`: RSS feed sources with metadata
- `feed_follows`: Many-to-many relationship between users and feeds
- `posts`: Individual articles/posts from RSS feeds
## Troubleshooting
### Common Issues
1. "**command not found: gator**"
   - Ensure `$GOPATH/bin` is in your PATH
   - Try `~/go/bin/gator` directly
2. **Database connection errors:**
   - Verify PostgreSQL is running: `pg_isready`
   - Check your connection string in `~/.gatorconfig.json`
   - Ensure database exists: `createdb gator`
3. **Config file not found:**
   - Create `~/.gatorconfig.json` with proper database URL
   - Ensure file has correct JSON syntax
4. **RSS parsing errors:**
   - Some feeds may have invalid XML/dates
   - Check feed URL is accessible in browser
   - Look for error messages in aggregator output


