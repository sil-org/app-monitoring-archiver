terraform {
  cloud {
    organization = "gtis"

    workspaces {
      name = "app-monitoring-archiver"
    }
  }
}
