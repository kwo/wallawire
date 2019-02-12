FROM ww/ui:latest AS ui
FROM ww/golang:latest AS builder

RUN mkdir /user \
 && echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd \
 && echo 'nobody:x:65534:' > /user/group \
 && mkdir -p /mnt/certshttp && chown nobody:nobody /mnt/certshttp \
 && mkdir -p /mnt/certsdb   && chown nobody:nobody /mnt/certsdb

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./
COPY --from=ui /public ./public

# RUN git clean -fdx
RUN go mod download
RUN go generate ./...

# Build the executable to `/app`. Mark the build as statically linked.
RUN ./scripts/build.sh

# Final stage: the running container.
FROM scratch AS final

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /mnt/certshttp /mnt/certshttp
COPY --from=builder /mnt/certsdb   /mnt/certsdb
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/wallawire /wallawire

EXPOSE 8888
VOLUME /mnt/certshttp
VOLUME /mnt/certsdb
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/wallawire"]
