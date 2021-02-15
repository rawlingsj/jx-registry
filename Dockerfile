FROM gcr.io/jenkinsxio/jx-cli-base:0.0.21

ENTRYPOINT ["jx-registry"]

COPY ./build/linux/jx-registry /usr/bin/jx-registry