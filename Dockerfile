FROM alpine:3.11
ENTRYPOINT ["/post-to-socials"]
RUN apk add --update --no-cache ca-certificates
ADD post-to-socials /
