FROM ghcr.io/astral-sh/uv:0.6.1 AS uv

FROM python:3.11-alpine AS builder
WORKDIR /app/

RUN  \
    --mount=type=bind,from=uv,source=/uv,target=/uv \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    /uv venv --relocatable && \
    /uv sync --no-dev --frozen --compile-bytecode --no-install-project --no-editable

COPY src /app/src

RUN  \
    --mount=type=bind,from=uv,source=/uv,target=/uv \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    /uv sync --no-dev --frozen --compile-bytecode --no-install-project --no-editable

FROM python:3.11-alpine
WORKDIR /app/

ENV PYTHONUNBUFFERED=1 \
    PYTHONUSERBASE="/app/venv" \
    PYTHONPATH="/app" \
    PATH="/app/venv/bin:$PATH"

COPY --from=builder /app/.venv /app/.venv
COPY src /app/src

USER nobody:nobody
CMD ["/app/.venv/bin/python", "src/main.py"]
