# Operations Runbook

> Generated from source of truth: config.toml.example, Makefile

## Deployment Procedures

### Build Release

```bash
# Single platform
make build

# All platforms (darwin/linux/windows, amd64/arm64)
make build-all
```

Output binaries in `bin/`:
- `terminalizcrazy-darwin-amd64`
- `terminalizcrazy-darwin-arm64`
- `terminalizcrazy-linux-amd64`
- `terminalizcrazy-linux-arm64`
- `terminalizcrazy-windows-amd64.exe`

### Installation

```bash
# From source
make install  # Copies to $GOPATH/bin

# Or manual
cp bin/terminalizcrazy /usr/local/bin/
```

### First Run Setup

1. Create config directory (auto-created on first run):
   ```bash
   mkdir -p ~/.terminalizcrazy
   ```

2. Copy example config:
   ```bash
   cp config.toml.example ~/.terminalizcrazy/config.toml
   ```

3. For local AI (default):
   ```bash
   ollama pull gemma4
   ollama serve
   ```

4. For cloud AI, set API key:
   ```bash
   export GEMINI_API_KEY="your-key"
   # or
   export ANTHROPIC_API_KEY="your-key"
   ```

## Monitoring and Alerts

### Log Levels

Set via config or environment:

| Level | Description |
|-------|-------------|
| `debug` | All messages (development) |
| `info` | Standard operation (default) |
| `warn` | Warnings and errors |
| `error` | Errors only |

```bash
export LOG_LEVEL=debug
export DEBUG=true
```

### Health Checks

**Ollama Connection:**
```bash
curl http://localhost:11434/api/tags
```

**Collaboration Server (if running):**
```bash
curl http://localhost:8765/health
```

### Database Location

SQLite database: `~/.terminalizcrazy/terminalizcrazy.db`

Tables:
- `sessions` - Chat sessions
- `messages` - Chat messages
- `command_history` - Executed commands
- `agent_plans` - Agent execution plans
- `agent_tasks` - Individual tasks in plans
- `workflows` - Saved workflow templates
- `workspaces` - Layout persistence

## Common Issues and Fixes

### AI Provider Issues

**"API key is required"**
- Ensure environment variable is set correctly
- Check for typos in key format
- For Ollama: ensure `ollama_enabled = true` in config

**Ollama connection failed**
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Start Ollama
ollama serve

# Pull model if missing
ollama pull gemma4
```

**Rate limiting (cloud providers)**
- Switch to Ollama for local inference
- Reduce request frequency
- Check API quota/billing

### TUI Issues

**Display corruption**
```bash
# Reset terminal
reset
# or
tput reset
```

**Keybindings not working**
- Check terminal emulator compatibility
- Some terminals intercept Ctrl+\ or other keys
- Try running in different terminal

### Storage Issues

**Database locked**
- Only one instance can run at a time
- Kill any background processes:
  ```bash
  pkill terminalizcrazy
  ```

**Corrupted database**
```bash
# Backup and recreate
mv ~/.terminalizcrazy/terminalizcrazy.db ~/.terminalizcrazy/terminalizcrazy.db.bak
# Restart application (creates new db)
```

### Collaboration Issues

**Cannot share session**
- Check port 8765 is available
- Firewall may block WebSocket connections
- Both parties need network connectivity

**Cannot join session**
- Verify share code is correct (format: `xxxx-yyyy`)
- Host must be running and sharing
- Check network/firewall settings

## Rollback Procedures

### Configuration Rollback

```bash
# Restore previous config
cp ~/.terminalizcrazy/config.toml.bak ~/.terminalizcrazy/config.toml
```

### Database Rollback

```bash
# Restore from backup
cp ~/.terminalizcrazy/terminalizcrazy.db.bak ~/.terminalizcrazy/terminalizcrazy.db
```

### Binary Rollback

```bash
# If using make install
# Rebuild from previous tag
git checkout v0.x.x
make build
make install
```

## Data Retention (GDPR)

Configured in `config.toml` under `[retention]`:

| Setting | Default | Description |
|---------|---------|-------------|
| `message_retention_days` | 90 | Chat message retention |
| `command_history_retention_days` | 90 | Command history retention |
| `agent_plan_retention_days` | 30 | Agent plan retention |
| `auto_cleanup_enabled` | true | Auto-delete on startup |

### Manual Cleanup

```sql
-- Connect to database
sqlite3 ~/.terminalizcrazy/terminalizcrazy.db

-- Delete old messages
DELETE FROM messages WHERE created_at < datetime('now', '-90 days');

-- Delete old commands
DELETE FROM command_history WHERE created_at < datetime('now', '-90 days');

-- Vacuum to reclaim space
VACUUM;
```

## Security Checklist

- [ ] SecretGuard enabled (`secret_guard_enabled = true`)
- [ ] API keys stored in environment, not config file
- [ ] Agent mode set appropriately (prefer `suggest` over `auto`)
- [ ] Data retention configured for compliance
- [ ] Collaboration uses E2E encryption (automatic)

## Performance Tuning

### Reduce Memory Usage

```toml
[workspace]
max_workspaces = 5

history_limit = 500
```

### Disable Animations

```toml
[appearance]
enable_animations = false
```

### Limit History

```toml
history_limit = 500
```

## Backup Procedures

### Full Backup

```bash
tar -czf terminalizcrazy-backup-$(date +%Y%m%d).tar.gz \
  ~/.terminalizcrazy/config.toml \
  ~/.terminalizcrazy/terminalizcrazy.db \
  ~/.terminalizcrazy/themes/
```

### Restore

```bash
tar -xzf terminalizcrazy-backup-YYYYMMDD.tar.gz -C ~/
```
