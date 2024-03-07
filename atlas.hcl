env "local" {
  src = "file://internal/database/schema.hcl"
  dev = "docker://postgres/16/dev?search_path=public"
  migration {
    dir = "file://internal/database/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}