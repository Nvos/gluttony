env "local" {
  src = "file://migrations/schema.hcl"
  dev = "docker://postgres/16/dev?search_path=public"
  migration {
    dir = "file://migrations"
    format = goose
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}