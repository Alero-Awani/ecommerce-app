# --------- Stage 1: Dependencies + Build ---------
FROM node:17-alpine AS builder
WORKDIR /app

# Install only what's needed to build the app
COPY package.json yarn.lock ./
RUN npm install --legacy-peer-deps

# Copy all app files
COPY . .

# Build the app (skip if your dev script does this in-memory)
RUN yarn build

# --------- Stage 2: Dev Environment (for local use) ---------
FROM node:17-alpine AS dev
WORKDIR /app

# Copy source and dependencies from builder
COPY --from=builder /app /app

EXPOSE 3030

# Start the app in dev mode
CMD ["yarn", "dev"]