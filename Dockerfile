FROM gcr.io/distroless/base
ARG BIN
COPY /bin/ep /ep
CMD ["/ep"]
