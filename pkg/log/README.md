## pkg/log

Simple logging package based on log/slog.

## Environment Variables

| Name | Description | Default | Required |
|------|-------------|---------| -------- |
| LOG_FORMAT | Log format (`gcp`, `json`, `text`) | `text` | |
| SERVICE | Service name. use for gcp log label. |  | |
| ENV | Environment. use for gcp log label. | | |
| GOOGLE_PROJECT_ID | Google Project Id. use for gcp logging | | |
