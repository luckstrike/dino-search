WORK IN PROGRESS (this might work):
----------------------------------------
If you are having issues when setting up the permissions of the PostgreSQL DB:
- Try running the init_db_permissions.sh script first and then running
  the main.go file in cmd/cli
- You can test if everything works with the test_db_setup.sh script in the test folder
- If you can't run either of the scripts try changing the permissions of the
  files with chmod
