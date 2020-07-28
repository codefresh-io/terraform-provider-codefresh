# current_account data module

Return current account (owner of the token) and its users 
```hcl
provider "codefresh" {
  api_url =  var.api_url 
  token = var.token 
}

data "codefresh_current_account" "acc" {
  
}


output "current_ac" {
  value = data.codefresh_current_account.acc
}
```

The output example: 
```
Outputs:

current_ac = {
  "_id" = "5f1fd9044d0fc94ddff0d745"
  "id" = "5f1fd9044d0fc94ddff0d745"
  "name" = "acc1"
  "users" = [
    {
      "email" = "kosta@codefresh.io"
      "id" = "5f1fd9094d0fc9c656f0d75a"
      "name" = "user1"
    },
    {
      "email" = "kosta@sysadmiral.io"
      "id" = "5f1fd9094d0fc93b52f0d75c"
      "name" = "user3"
    },
  ]
}
```