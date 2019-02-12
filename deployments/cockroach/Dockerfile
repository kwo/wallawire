FROM alpine:latest AS builder

# Create the user and group files that will be used in the running container to
# run the process an unprivileged user.
RUN mkdir /user \
 && echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd \
 && echo 'nobody:x:65534:' > /user/group \
 && mkdir -p /mnt/certs  && chown nobody:nobody /mnt/certs \
 && mkdir -p /mnt/data   && chown nobody:nobody /mnt/data

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates curl \
 && curl https://binaries.cockroachdb.com/cockroach-v2.1.4.linux-musl-amd64.tgz | tar --strip-components=1 -xz \
 && chmod 755 cockroach

# Final stage: the running container.
FROM scratch AS final

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /mnt/certs /mnt/certs
COPY --from=builder /mnt/data  /mnt/data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /cockroach /cockroach

EXPOSE 5432 8080
VOLUME /mnt/data
VOLUME /mnt/certs
USER nobody:nobody
CMD ["/cockroach", "start", "--certs-dir=/mnt/certs", "--host=0.0.0.0", "--port=5432", "--http-port=8080", "--store=path=/mnt/data"]
