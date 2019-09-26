FROM ubuntu

COPY bin/satellity satellity

COPY internal/configs/config.yaml internal/configs/config.yaml

CMD ["./satellity"]
