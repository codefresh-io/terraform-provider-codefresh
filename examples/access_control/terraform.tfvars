api_url = "https://onprem-tst-1.cf-cd.com/api"

default_idps = {
    local = {
       display_name = "local"
       sso = false
    }
    azure_sso = { 
        display_name = "codefresh-azure-sso-1"
        sso = true
    }
}

accounts = {
    acc1 = {}
    acc2 = {
        limits = {
            collaborators = 50
            parallel_builds = 5
        }
    }
}

users = {
    user1 = {
        email = "kosta@codefresh.io"
        personal = {
            first_name = "Kosta"
            last_name = "A"
        }
        accounts = ["acc1", "acc2"]
        admin_of_accounts = ["acc1"]
        global_admin = true
    }
    user2 = {
        email = "live.com#kosta777@gmail.com"
        personal = {
            first_name = "Q"
            last_name = "D"
        }
        accounts = ["acc2"]
        admin_of_accounts = []
        global_admin = false
    }
    user3 = {
        email = "kosta@sysadmiral.io"
        personal = {
            first_name = "Kosta"
            last_name = "sysadmiral-io"
        }
        accounts = ["acc1", "acc2"]
        admin_of_accounts = ["acc1", "acc2"]
        global_admin = true
    }
}