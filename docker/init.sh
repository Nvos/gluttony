#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" "$POSTGRES_DB" <<-EOSQL
  create text search dictionary polish_spell (template=ispell, dictfile=polish, afffile=polish, stopwords=polish);
  create text search configuration polish (copy=english);
  alter text search configuration polish alter mapping for word, asciiword with polish_spell, simple;
  create extension pg_trgm;
EOSQL