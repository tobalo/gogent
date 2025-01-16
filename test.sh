nats pub agent.technical.support '{"timestamp":"2025-01-15T02:14:23.123Z","hostname":"web-server-01","severity":"ERROR","service":"nginx","message":"Failed to bind to port 80: Address already in use","context":{"pid":1234,"user":"www-data"}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:15:01.456Z","hostname":"db-server-03","severity":"CRITICAL","service":"mysql","message":"InnoDB: Error: log sequence numbers 12345678 and 12345679 in ibdata files do not match","context":{"thread_id":5678}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:15:45.789Z","hostname":"app-server-02","severity":"WARNING","service":"systemd","message":"Process /usr/bin/app-service (PID 9876) had a non-zero exit code","context":{"exit_code":1,"unit":"app-service.service"}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:16:12.234Z","hostname":"web-server-01","severity":"ERROR","service":"apache2","message":"AH00058: Error: could not open error log file /var/log/apache2/error.log","context":{"permissions":"0644","user":"www-data"}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:17:33.567Z","hostname":"db-server-03","severity":"CRITICAL","service":"postgresql","message":"FATAL: could not write to WAL: No space left on device","context":{"database":"customers_db","disk_usage":"99.9%"}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:18:01.890Z","hostname":"app-server-02","severity":"WARNING","service":"kernel","message":"CPU temperature above threshold, cpu clock throttled","context":{"core":2,"temp":"95C","freq":"2.1GHz"}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:19:23.123Z","hostname":"web-server-01","severity":"ERROR","service":"fail2ban","message":"Banned IP 192.168.1.100 for multiple failed SSH attempts","context":{"attempts":5,"jail":"sshd","duration":"3600s"}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:20:45.456Z","hostname":"db-server-03","severity":"CRITICAL","service":"redis","message":"MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk","context":{"used_memory":"6.2G","maxmemory":"8G"}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:21:12.789Z","hostname":"app-server-02","severity":"WARNING","service":"docker","message":"Container e5a3e3d4f6g7 exceeded memory limit","context":{"container_name":"production-api","limit":"2G","actual":"2.3G"}}'
nats pub agent.technical.support '{"timestamp":"2025-01-15T02:22:33.234Z","hostname":"web-server-01","severity":"ERROR","service":"certbot","message":"Failed to renew SSL certificate for domain.com","context":{"error":"DNS validation failed","attempts":3}}'