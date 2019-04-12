### 定时任务触发器  

#### conf/app.conf内容  
runmode = test  
appname = wScheduler  
httpaddr = 0.0.0.0  
httpport = 8888  
dbUser = xxx    
dbPassword = xxx  
dbName = wScheduler  
mailbox = xxx  
mailboxPassword = xxx  

[test]  
servername = 0.0.0.0  
[prod]  
servername = xxx  

#### job
job cron格式:详见quartz/quartz_test.go  
job urls:多个url之间以;间隔 目前的触发策略是依次访问每一个url  
job Params使用json格式  
