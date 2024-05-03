-- +goose Up
DROP TEXT SEARCH CONFIGURATION IF EXISTS polish cascade;
CREATE TEXT SEARCH DICTIONARY polish_ispell (
    Template = ispell,
    DictFile = polish,
    AffFile = polish,
    StopWords = polish
    );

CREATE TEXT SEARCH CONFIGURATION polish( COPY = pg_catalog.english);

ALTER TEXT SEARCH CONFIGURATION polish
    ALTER MAPPING FOR asciiword, asciihword, hword_asciipart, word, hword, hword_part
        WITH polish_ispell;

CREATE EXTENSION pg_trgm;

-- +goose Down
DROP TEXT SEARCH CONFIGURATION polish cascade;
