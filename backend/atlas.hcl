variable "mysql_root_password" {
  type    = string
  default = getenv("MYSQL_ROOT_PASSWORD")
}

data "ent" "app" {
  path = "./ent/schema"
}

env "local" {
  src = data.ent.app.url
  dev = "mysql://root:${var.mysql_root_password}@localhost:3306/dev?parseTime=true"
  migration {
    dir = "file://db/migrations"
  }
}
