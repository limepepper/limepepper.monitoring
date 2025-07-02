# Ansible role to configure auditd

This role is mostly derived from experience with auditd on Fedora, but it may
be useful for other distributions as well. However, it is not tested on
other distributions.

auditd is a special service in fedora, which doesn't follow normal systemd
systemctl commands.

```
$ sudo auditctl --signal reload
$ sudo auditctl --signal rotate
$ sudo auditctl --signal resume
$ sudo auditctl --signal state
```


```
# auditctl -s
enabled 1
failure 1
pid 1660857
rate_limit 0
backlog_limit 64
lost 0
backlog 0
backlog_wait_time 60000
backlog_wait_time_actual 0
loginuid_immutable 0 unlocked
```

## background

> The Audit system consists of two main parts: the user-space applications and utilities, and the kernel-side system call processing. The kernel component receives system calls from user-space applications and filters them through one of the following filters: user, task, fstype, or exit.
<https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/8/html/security_hardening/auditing-the-system_security-hardening>


## generating and monitoring privileged commands

```
## You can run the following commands to generate the rules (don't forget to
## add arch=b32 rules, too):
find /bin -type f -perm -04000 2>/dev/null | awk '{ printf "-a always,exit -F arch=b64 -F path=%s -F perm=x -F auid>=1000 -F auid!=unset -F key=privileged\n", $1 }' > priv.rules
find /sbin -type f -perm -04000 2>/dev/null | awk '{ printf "-a always,exit -F arch=b64 -F path=%s -F perm=x -F auid>=1000 -F auid!=unset -F key=privileged\n", $1 }' >> priv.rules
find /usr/bin -type f -perm -04000 2>/dev/null | awk '{ printf "-a always,exit -F arch=b64 -F path=%s -F perm=x -F auid>=1000 -F auid!=unset -F key=privileged\n", $1 }' >> priv.rules
find /usr/sbin -type f -perm -04000 2>/dev/null | awk '{ printf "-a always,exit -F arch=b64 -F path=%s -F perm=x -F auid>=1000 -F auid!=unset -F key=privileged\n", $1 }' >> priv.rules
filecap /bin 2>/dev/null | sed '1d' | awk '{ printf "-a always,exit -F path=%s -F arch=b64 -F perm=x -F auid>=1000 -F auid!=unset -F key=privileged\n", $2 }' >> priv.rules
filecap /sbin 2>/dev/null | sed '1d' | awk '{ printf "-a always,exit -F path=%s -F arch=b64 -F perm=x -F auid>=1000 -F auid!=unset -F key=privileged\n", $2 }' >> priv.rules
filecap /usr/bin 2>/dev/null | sed '1d' | awk '{ printf "-a always,exit -F arch=b64 -F path=%s -F perm=x -F auid>=1000 -F auid!=unset -F key=privileged\n", $2 }' >> priv.rules
filecap /usr/sbin 2>/dev/null | sed '1d' | awk '{ printf "-a always,exit -F arch=b64 -F path=%s -F perm=x -F auid>=1000 -F auid!=unset -F key=privileged\n", $2 }' >> priv.rules
```



## Notes


- audit documentation project:
  <https://github.com/linux-audit/audit-documentation/wiki>
