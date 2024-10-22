# mysqlDiffChecker
Check the difference of User and Variables between two MySQL instances.
# install
```
go get github.com/go-sql-driver/mysql
```

```
# ./mysqlDiffChecker -h
Usage of ./mysqlDiffChecker:
  -m string
    	Mode [Variables,User]
  -o int
    	outputWidth (default 50)
  -p string
    	Source and Target Password
  -s string
    	SourceServer[ip:port]
  -sp string
    	Source Password
  -su string
    	Source User
  -t string
    	TargetServer[ip:port]
  -tp string
    	Target Password
  -tu string
    	Target User
  -u string
    	Source and Target User
```

* Variables Check
```
# ./mysqlDiffChecker -m Variables -s 192.168.1.1:3306 -t 192.168.1.2:3306 -u test -p test
Variables                                      192.168.1.1:3306                                   192.168.1.2:3306
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
version                                        5.7.22-log                                         8.0.16
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
histogram_generation_max_mem_size              NOT_FOUND_THE_VARIABLE                             20000000
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
time_format                                    %H:%i:%s                                           NOT_FOUND_THE_VARIABLE
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
cte_max_recursion_depth                        NOT_FOUND_THE_VARIABLE                             1000
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
log_builtin_as_identified_by_password          OFF                                                NOT_FOUND_THE_VARIABLE
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
server_id                                      176134023                                          101
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
default_collation_for_utf8mb4                  NOT_FOUND_THE_VARIABLE                             utf8mb4_0900_ai_ci
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
have_symlink                                   YES                                                DISABLED
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
regexp_time_limit                              NOT_FOUND_THE_VARIABLE                             32
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
tablespace_definition_cache                    NOT_FOUND_THE_VARIABLE                             256
---------------------------------------------- -------------------------------------------------- --------------------------------------------------
relay_log_recovery                             OFF                                                ON
---------------------------------------------- -------------------------------------------------- --------------------------------------------------.
.
.
.
.

```

* User Chaek

```
# ./mysqlDiffChecker -m User -s 192.168.1.1:3306 -t 192.168.1.2:3306 -u test -p test -o 100
User                           192.168.1.1:3306                                                                                     192.168.1.2:3306
------------------------------ ---------------------------------------------------------------------------------------------------- ----------------------------------------------------------------------------------------------------
'test1'@'%'                    GRANT INSERT ON *.* TO 'test1'@'%'                                                                   GRANT SELECT ON *.* TO `test1`@`%`
------------------------------ ---------------------------------------------------------------------------------------------------- ----------------------------------------------------------------------------------------------------
'test2'@'192.168.1.2'                                                                                                               GRANT SELECT ON *.* TO `test2`@`192.168.1.2`
------------------------------ ---------------------------------------------------------------------------------------------------- ----------------------------------------------------------------------------------------------------
'test2'@'192.168.1.1'          GRANT SELECT ON *.* TO 'test2'@'192.168.1.1'
------------------------------ ---------------------------------------------------------------------------------------------------- ----------------------------------------------------------------------------------------------------
.
.
.
.
```
