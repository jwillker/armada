FROM alpine:3.17.0

RUN addgroup -S -g 2000 armada && adduser -S -u 1000 armada -G armada

USER armada

# Server
COPY server /app/
COPY config/armada/config.yaml /app/config/armada/config.yaml

# Executor
COPY executor /app/
COPY config/executor/config.yaml /app/config/executor/config.yaml

# Binoculars
COPY binoculars /app/
COPY config/binoculars/config.yaml /app/config/binoculars/config.yaml

# Lookout Ingester
COPY lookoutingester /app/
COPY config/eventingester/config.yaml /app/config/eventingester/config.yaml

# Event Ingester
COPY eventingester /app/
COPY config/eventingester/config.yaml /app/config/eventingester/config.yaml

# Job Service
COPY jobservice /app/
COPY config/jobservice/config.yaml /app/config/jobservice/config.yaml

WORKDIR /app
