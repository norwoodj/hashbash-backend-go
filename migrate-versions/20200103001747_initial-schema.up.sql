CREATE TABLE rainbow_table
(
    id                 SMALLSERIAL NOT NULL,
    name               TEXT        NOT NULL,
    num_chains         BIGINT      NOT NULL,
    chain_length       BIGINT      NOT NULL,
    password_length    SMALLINT    NOT NULL,
    character_set      TEXT        NOT NULL,
    hash_function      TEXT        NOT NULL,
    final_chain_count  BIGINT      NOT NULL DEFAULT 0,
    chains_generated   BIGINT      NOT NULL DEFAULT 0,
    status             TEXT        NOT NULL,
    generate_started   TIMESTAMP            DEFAULT NULL,
    generate_completed TIMESTAMP            DEFAULT NULL,
    created            TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (name)
);


CREATE TABLE rainbow_chain
(
    start_plaintext  TEXT     NOT NULL,
    end_hash         TEXT     NOT NULL,
    rainbow_table_id SMALLINT NOT NULL,
    PRIMARY KEY (rainbow_table_id, end_hash),
    CONSTRAINT rainbow_chain_rainbowTableId_fk FOREIGN KEY (rainbow_table_id) REFERENCES rainbow_table (id) ON DELETE CASCADE
);


CREATE TABLE rainbow_table_search
(
    id               BIGSERIAL NOT NULL,
    rainbow_table_id SMALLINT  NOT NULL,
    hash             VARCHAR   NOT NULL,
    status           VARCHAR   NOT NULL,
    password         VARCHAR            DEFAULT NULL,
    search_started   TIMESTAMP          DEFAULT NULL,
    search_completed TIMESTAMP          DEFAULT NULL,
    created          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (rainbow_table_id, hash),
    CONSTRAINT rainbow_table_search_rainbowTableId_fk FOREIGN KEY (rainbow_table_id) REFERENCES rainbow_table (id) ON DELETE CASCADE
);

CREATE INDEX ON rainbow_table_search (status);
