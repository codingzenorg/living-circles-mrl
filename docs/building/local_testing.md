# Local Testing

## Purpose

This note records the current local testing loop for Living Circles so it does not depend on chat history.

## Recommended Local Loop

Use `air` when you want fast manual feedback from the Go server while changing simulation code or transport code.

```bash
source "$HOME/.nvm/nvm.sh"
nvm use
air
```

Then open:

```text
http://localhost:8080
```

This is the preferred loop for:

- checking the browser demo quickly after server-side edits
- validating interaction behavior by hand
- avoiding repeated manual `go run` restarts

## Deterministic Validation

`air` is only for fast manual feedback. It does not replace deterministic validation.

Use these commands for actual slice verification:

```bash
source "$HOME/.nvm/nvm.sh"
nvm use
npm run test:contracts
go test ./...
```

## Practical Rule

- use `air` for quick demo feedback
- use `npm run test:contracts` and `go test ./...` before accepting or committing a slice
