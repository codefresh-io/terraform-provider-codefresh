# account data module

```
data "codefresh_account" "acc" {
  name = "acc1"
}

resource "codefresh_user" "user1" {
  email = "user1@example.com"
  user_name = "user1"

  accounts = [
    data.codefresh_account.acc.id
  ]

  activate = true

  roles = [
    "Admin",
    "User"
  ]

  login {
      idp_id = data.codefresh_idps.idp_azure.id
      sso = true
  }
  
  login  {
      idp_id = data.codefresh_idps.local.id
      //sso = false
  }


  personal {
    first_name = "John"
    last_name = "Smith"
  }
}
```