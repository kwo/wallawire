ARG REPO="ww/"
ARG VERSION="unknown"

FROM ${REPO}nodejs:latest AS builder
ARG VERSION

WORKDIR /ww/ui
COPY ./ ./
RUN yarn install && yarn build --VERSION=${VERSION}

FROM scratch AS final
COPY --from=builder /ww/ui/public/ /public/
