# go-backuper

![Go](https://github.com/linuxoid69/go-backuper/workflows/Go/badge.svg)
![Docker Image CI](https://github.com/linuxoid69/go-backuper/workflows/Docker%20Image%20CI/badge.svg)  

## Manual

The tool for backup files and DB (postgres)  
Below consider the example of config.  

```yaml
# path where storage backups
backup_storage_path: /mnt/data

# Name a zip file suffix. File will be look as <project_name>_19_03_2020_backup.zip
# project_name - from section backup_source_files_path. See below.
name_zip_file: backup.zip

# If you want to encrypt your backups set true
encrypt_backup: true
# Password is for encrypt your backups
encrypt_password: "1234"

# You time zone
time_zone: "Europe/Moscow"
# cron for database
cron_db: "0 1 * * *" # cron for backup of DB
# cron for files
cron_files: "0 2 * * *" # cron for backup of files

# project_name - name of project, dirs - directories of project must be separate comma whitout whitespace.
backup_source_files_path:
  projects:
    - - project_name: redmine
        dirs: "/opt/redmine"
      - project_name: testproject
        dirs: "/opt/testproject,/var/lib/testproject"

# Database for backup.
database:
  host: 127.0.0.1
  user: postgres
  password: 1234
  port: 5432
  options: "" # additional options e.g. "-v -j3" among flags must be one whitespace and between flag and value don't be whitespace
  dbnames:
    - demo
```

### Command line

```shell
Usage of ./go-backuper:
  -c string
        config (default "./config.yaml")
  -d    Run as daemon
  -f string
        Name of a file is for decrypt
  -p string
        Password is for decrypt
```

If you need decrypt file run that command:

```bash
./go-backuper -f test.txt.enc -p 1234
```
