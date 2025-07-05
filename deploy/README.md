# Logrotate

1. Update `logrotate` path to be exactly same as `GLUTTONY_LOG_FILE_PATH`
2. Ensure that file path matching `GLUTTONY_LOG_FILE_PATH` exists (only directory is required)
3. Copy `logrotate` to `/etc/logrotate.d/gluttony`

# System.d service

1. Copy `gluttony.service` to `/lib/systemd/system/gluttony.service`
2. Copy `gluttony` binary to `/usr/bin/gluttony`
3. Create configuration env in `/etc/gluttony/config.env`

# Caddy

1. Update `/etc/caddy/Caddyfile`