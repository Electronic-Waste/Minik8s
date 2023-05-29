#!/usr/bin/expect
spawn ssh stu1642@pilogin.hpc.sjtu.edu.cn "cd test && ls"
expect "*\[fingerprint\]\)\?"
send "yes\r"
expect "*Password:"
send "h4L&\$IQW\r"
interact