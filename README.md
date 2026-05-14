# logpipe

A composable log filtering and forwarding tool that reads from stdin and routes to multiple outputs based on rules.

---

## Installation

```bash
go install github.com/yourname/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logpipe.git && cd logpipe && go build -o logpipe .
```

---

## Usage

Pipe any log-producing command into `logpipe` and define routing rules via a config file:

```bash
./myapp | logpipe --config logpipe.yaml
```

**Example `logpipe.yaml`:**

```yaml
outputs:
  - name: errors-file
    type: file
    path: /var/log/errors.log
    filter: 'level == "error"'

  - name: stdout-all
    type: stdout

  - name: remote
    type: http
    url: https://logs.example.com/ingest
    filter: 'level == "warn" || level == "error"'
```

Lines are read from stdin, evaluated against each output's filter rule, and forwarded to all matching destinations concurrently.

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `logpipe.yaml` | Path to the configuration file |
| `--dry-run` | `false` | Print matched output names without writing |
| `--verbose` | `false` | Log internal routing decisions to stderr |

---

## Features

- Reads from stdin — works with any log source
- Route the same log line to multiple outputs simultaneously
- Filter by structured fields (JSON logs) or regex patterns
- Output targets: stdout, file, HTTP endpoint, and more

---

## License

MIT © yourname
