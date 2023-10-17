# variable api_url {
#   type = string
# }

# variable token {
#   type = string
#   default = ""
# }

# teams map[team_name]usersList
# {
#   developers = ["user1", "user3"]
#   managers = ["user3", "user2"]
# }
variable teams {
  type = map(any)
}
