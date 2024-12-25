# S3ase

Simplifying S3 usage through open source.

> [!IMPORTANT]
> Currently in development.

## Tech stack

### Backend

- Golang
- Postgresql
- Redis

### Frontend

- Remix
- Shadcn UI
- Tailwind CSS

## Features

| Status | Feature                      |
| :----- | ---------------------------- |
| ✅     | Open source                  |
| ⚙️     | Github/Google Authentication |
| ⚙️     | Manage webhooks              |
| ⚙️     | Manage buckets               |
| ⚙️     | Manage API                   |
| ⚙️     | Manage folders               |
| ⚙️     | Manage files                 |

## Get started

To get started, you can run this project locally using 2 options

- [Docker](#docker) : Using docker-compose

```bash
# dev
docker-compose up -d ./docker-compose.dev.yml
```

- [Local](#local): To run the project locally, make sure you have the following installed:
  - [Golang](https://go.dev/doc/install)
  - [Postgresql](https://www.postgresql.org/download/)
  - [Redis](https://redis.io/docs/getting-started/installation/)
  - [Node](https://nodejs.org/en/download)

```bash
# backend
cd backend

go build -o s3ase ./cmd

# migrate the database
./s3ase migrate up

# seed the database
./s3ase migrate seed

# run the server
./s3ase server

# frontend
cd frontend && pnpm install && pnpm run dev
```

## Contributing

## License

This project is licensed under the [GNU AFFERO GENERAL PUBLIC LICENSE](https://github.com/s3ase/s3ase/blob/main/LICENSE).

## Activities

![Activities](https://repobeats.axiom.co/api/embed/2715c72646a8bcd65826fa8e45833d55404a847e.svg "Repobeats analytics image")
