# S3ase

S3ase is an open-source file management service backed by S3-compatible object storage. It provides a Remix dashboard and a Go API for managing apps, files, folders, members, API keys, webhooks, custom domains, usage, and plan state.

## Stack

- Go, Chi, PostgreSQL, Redis, Asynq, and the AWS S3 SDK
- Remix, React, Tailwind CSS, and shadcn/ui
- MinIO for local S3-compatible development

## Run with Docker

Docker Compose is the quickest way to start the full development stack:

```bash
docker compose -f docker-compose.dev.yml up -d --build
```

The migration container runs before the API starts. Services are then available at:

- Dashboard: `http://127.0.0.1:3000`
- API: `http://127.0.0.1:5555/api`
- MinIO console: `http://127.0.0.1:9001`
- PostgreSQL: `localhost:5434`
- Redis: `localhost:6380`

Configure either GitHub or Google OAuth credentials before using browser sign-in. The callback URL is `http://localhost:3000/auth/callback`; the backend redirect URL is `http://localhost:5555`.

Stop the stack without deleting data:

```bash
docker compose -f docker-compose.dev.yml down
```

## Run locally

You need Go 1.26.5+, Node.js 20+, pnpm 9, PostgreSQL, Redis, and S3-compatible storage.

```bash
cd backend
go build -o s3ase ./cmd
./s3ase migrate up
./s3ase server
```

In another terminal:

```bash
cd frontend
pnpm install --frozen-lockfile
pnpm dev
```

The backend reads configuration from environment variables. Important production settings include a unique 32-byte `ENCRYPTION_KEY`, database and Redis connection values, S3 credentials, OAuth credentials, allowed CORS origins, and `DOMAIN_CNAME_TARGET`.

## Object storage

AWS S3 and S3-compatible services such as MinIO remain the default. Configure them with `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, `AWS_DEFAULT_REGION`, and optionally `AWS_ENDPOINT` and `AWS_USE_PATH_STYLE`.

To use Cloudflare R2 instead, create an R2 API token with Admin Read & Write permission because S3ase creates and deletes an app bucket as part of the existing app lifecycle. Then configure:

```dotenv
OBJECT_STORAGE_PROVIDER=r2
R2_ACCOUNT_ID=your_cloudflare_account_id
R2_ACCESS_KEY_ID=your_r2_access_key_id
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
# Optional; supported values are eu and fedramp.
R2_JURISDICTION=
```

S3ase derives the correct R2 endpoint and uses R2's required `auto` region. The dashboard, direct uploads, and resumable tus uploads use the same UI and API paths for either provider.

The provider setting applies to the whole deployment. Changing it does not migrate objects that already exist in another storage provider.

## API keys

API-key secrets are displayed only once when created in the dashboard. Send the secret as a bearer token to `/api/v1`:

```bash
curl http://localhost:5555/api/v1/files \
  -H "Authorization: Bearer s3ase_your_key"

curl -X POST "http://localhost:5555/api/v1/files?name=hello.txt" \
  -H "Authorization: Bearer s3ase_your_key" \
  -H "Content-Type: text/plain" \
  --data-binary "hello"
```

Read keys can list and download. Write keys can mutate. Full keys can do both.

## Webhooks and domains

Webhooks support upload started, completed, failed, and deleted events. Delivery runs through the background queue with retries; targets must use HTTPS and cannot resolve to private network addresses.

Custom-domain verification requires both records shown in the dashboard:

- A CNAME pointing to `DOMAIN_CNAME_TARGET`
- A TXT record at `_s3ase.<your-domain>` containing the generated ownership token

## Billing

The app persists and displays its current plan. Paid checkout is intentionally unavailable until a payment provider is configured; the UI does not simulate upgrades or charges.

## Verification

```bash
cd backend && go test ./...

cd frontend
pnpm lint
pnpm typecheck
pnpm test -- --run
pnpm build

docker compose -f docker-compose.dev.yml config --quiet
```

## License

S3ase is licensed under the [GNU Affero General Public License](LICENSE).
