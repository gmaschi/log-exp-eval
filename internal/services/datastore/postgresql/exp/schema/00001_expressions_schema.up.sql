CREATE TABLE IF NOT EXISTS expressions
(
    row_id        BIGSERIAL   NOT NULL,
    expression_id uuid        NOT NULL,
    expression    TEXT        NOT NULL,
    username     TEXT        NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT expression_id_pk PRIMARY KEY (expression_id)
);
