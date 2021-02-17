FROM ghcr.io/jenkins-x/jx-boot:3.1.252

ENTRYPOINT ["/run.sh"]

COPY ./build/linux/jx-registry /usr/bin/jx-registry
COPY run.sh /run.sh