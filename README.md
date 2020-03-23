# go-backuper
![Go](https://github.com/linuxoid69/go-backuper/workflows/Go/badge.svg)
![Docker Image CI](https://github.com/linuxoid69/go-backuper/workflows/Docker%20Image%20CI/badge.svg)
The tool for backup files and DB (postgres)  
Below consider the example config.
```yaml
# path where storage backups
backup_storage_path: /mnt/data

# Name a zip file suffix. File will be look as <project_name>_19_03_2020_backup.zip
# project_name - from section backup_source_files_path. See below.
name_zip_file: backup.zip

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
