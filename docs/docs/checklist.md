---
title: Troubleshooting
---

# Troubleshooting

Find below a list of actions used to troubleshoot *Káto*, this list is based on real issues and their solutions.
Some issues might have been permanently fixed but are keept here for its troubleshooting value.

<br>
<h4><span class="glyphicon glyphicon glyphicon-pencil" aria-hidden="true"></span> <em>Multiple HA-Proxy PIDs inside one container</em></h4>
<hr>

**Diagnose:**

```
loopssh worker "docker exec -i marathon-lb ps auxf | grep 'haproxy -p'"
```

**Mitigate:**

```
sudo systemctl restart marathon-lb
```

<br>
<h4><span class="glyphicon glyphicon glyphicon-pencil" aria-hidden="true"></span> <em>Disk usage</em></h4>
<hr>

**Diagnose:**

```
for i in quorum master worker border; do loopssh ${i} df -h; done
```

**Mitigate:**

```
sudo journalctl --vacuum-time=1h
docker rmi $(docker images -qf dangling=true)
```

<br>
<h4><span class="glyphicon glyphicon glyphicon-pencil" aria-hidden="true"></span> <em>Mixed CoreOS versions</em></h4>
<hr>

**Diagnose:**

```
for i in quorum master worker border; do
  loopssh ${i} "cat /etc/os-release | grep VERSION="
done
```

**Mitigate:**

```
update_engine_client -check_for_update
```

<br>
<h4><span class="glyphicon glyphicon glyphicon-pencil" aria-hidden="true"></span> <em>Summary of running containers (not realtime)</em></h4>
<hr>

```
for i in $(etcdctl ls /docker/images); do etcdctl get ${i}; done | \
sort | uniq -c | sort -n
```

<br>
<h4><span class="glyphicon glyphicon glyphicon-pencil" aria-hidden="true"></span> <em>The resource demand for a given task is higher than the available resources co-located on a single worker node. Therefore, the Marathon task stays in the waiting state forever.</em></h4>
<hr>

This is not really an error, you can:

 - Exo-scale up your cluster.
 - Redefine the task so it requires less resources.
 - Kill existing tasks in order to free resources.

<br>
<h4><span class="glyphicon glyphicon glyphicon-pencil" aria-hidden="true"></span> <em>Multiple Marathon frameworks registered but only one is expected to be up and running.</em></h4>
<hr>

Try to teardown the unexpected framework ID:

```
curl -H "Content-Type: application/json" -X POST \
  -d 'frameworkId=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee-ffff' \
  -v -L http://master:5050/teardown
```

<br>
<h4><span class="glyphicon glyphicon glyphicon-pencil" aria-hidden="true"></span> <em>Retrieve user-data</em></h4>
<hr>

To view `user-data` on EC2:

```
curl -s http://169.254.169.254/2009-04-04/user-data | gzip -d | less
```
