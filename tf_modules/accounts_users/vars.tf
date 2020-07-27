variable api_url {
  type = string
}

# 
variable token {
  type = string
  default = ""
}

variable default_account_features {
  type = map(any)
  default = {
    OfflineLogging = true,
    ssoManagement = true,
    teamsManagement = true,
    abac = true,
    customKubernetesCluster = true,
    launchDarklyManagement = false,
  }
}

variable default_acccount_limits {
    type = map(any)
    default = {
        collaborators = 100
        parallel_builds = 10
    }
}

variable default_idps {
    type = map(any)   
    default = {
        local  = { 
            display_name = "local"
            sso = false
        }
    } 
}

# map of accounts indexed by unique account name
# accounts = {
#     acc1 = {
#    }
#     acc2 = {
#         limits = {
#             collaborators = 50
#             parallel_builds = 5
#         }
#     }
# }
variable accounts {
  type = map(any)
}

# map of users:
# users = {
#     user1 = {
#          email = "ddd@gmail.com"
#          personal = {
#              first_name = "Q"
#              last_name = "D"
#          }
#          accounts = ["acc1", "acc2"]
#          global_admin = false
#      }
#     user2 = {
          
#         email = "ddd@gmail.com"
#         personal = {
#             first_name = "Q"
#             last_name = "D"
#         }
#         accounts = ["acc1", "acc2"]
#         global_admin = true
#     }
# }
variable users {
  //type = map(any)
}

