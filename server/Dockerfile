FROM python:3.11-bookworm as builder

RUN pip install poetry==1.8.3

ENV POETRY_NO_INTERACTION=1 \
    POETRY_VIRTUALENVS_IN_PROJECT=1 \
    POETRY_VIRTUALENVS_CREATE=1 \
    POETRY_CACHE_DIR=/tmp/poetry_cache

WORKDIR /backend

COPY pyproject.toml poetry.lock ./

RUN touch README.md

RUN --mount=type=cache,target=$POETRY_CACHE_DIR poetry install --no-root

FROM python:3.11-slim-bookworm

ENV VIRTUAL_ENV=/backend/.venv \
    PATH="/backend/.venv/bin:$PATH"

WORKDIR /backend

COPY --from=builder ${VIRTUAL_ENV} ${VIRTUAL_ENV}

COPY ./ /backend

EXPOSE 8000
CMD ["sh", "-c", "env"]

CMD ["sh", "-c", "python -m uvicorn main:app --host 0.0.0.0 --port 8000"]