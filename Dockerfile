FROM ghcr.io/jenkins-x/jx-boot:3.1.252

ENTRYPOINT ["jx-registry"]

CMD ["create"]

COPY ./build/linux/jx-registry /usr/bin/jx-registry