resource "codefresh_account_user_association" "user" {
    email = "terraform-test-user+user-change@codefresh.io"
}

resource "codefresh_account_user_association" "admin" {
    email = "terraform-test-user+admin-change@codefresh.io"
    admin = true
}