# Test
FROM node:18-alpine as test-target
# Check https://github.com/nodejs/docker-node/tree/b4117f9333da4138b03a546ec926ef50a31506c3#nodealpine to understand why libc6-compat might be needed.
RUN apk add --no-cache libc6-compat
WORKDIR /usr/src/app
RUN addgroup -S nextGroup && adduser -S nextjs -G nextGroup
ENV NODE_ENV=development
ENV PATH $PATH:/usr/src/app/node_modules/.bin

COPY --chown=nextjs:nextGroup package.json ./
RUN yarn install --frozen-lockfile

COPY --chown=nextjs:nextGroup . .

# Build
FROM test-target as build-target
ENV NODE_ENV=production
RUN yarn build

# Reduce installed packages to production-only.
RUN yarn install --production

# Archive
FROM node:18-alpine as production-target
USER nextjs
ENV NODE_ENV=production
ENV PATH $PATH:/usr/src/app/node_modules/.bin

WORKDIR /usr/src/app

# # Include only the release build and production packages.
COPY --chown=nextjs:nextGroup --from=build-target /usr/src/app/node_modules node_modules
COPY --chown=nextjs:nextGroup --from=build-target /usr/src/app/.next .next

EXPOSE 3000

ENV PORT 3000

CMD ["next", "start"]
