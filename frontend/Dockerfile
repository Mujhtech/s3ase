# ---------------------------------------------------------#
#                   Build node image base                  #
# ---------------------------------------------------------#
FROM node:20-alpine3.20 AS base

# Install all node_modules, including dev dependencies
FROM base AS deps

WORKDIR /app

RUN corepack enable

ADD package.json yarn.lock* package-lock.json* pnpm-lock.yaml* ./
RUN pnpm install --ignore-scripts --no-frozen-lockfile


# ---------------------------------------------------------#
#              Setup production node_modules               #
# ---------------------------------------------------------#
FROM base AS production-deps

ENV NODE_ENV production

WORKDIR /app

RUN corepack enable

COPY --from=deps /app/node_modules /app/node_modules
ADD package.json yarn.lock* package-lock.json* pnpm-lock.yaml* ./
RUN pnpm install --prod --no-frozen-lockfile

# ---------------------------------------------------------#
#            Build the app & its dependencies              #
# ---------------------------------------------------------#
FROM base AS build


WORKDIR /app

RUN corepack enable

COPY --from=deps /app/node_modules /app/node_modules


ADD . .
RUN pnpm run build

# Finally, build the production image with minimal footprint
FROM base

ENV NODE_ENV production


WORKDIR /app

RUN corepack enable

COPY --from=build /app /app

EXPOSE 3000

CMD ["pnpm", "run", "start"]