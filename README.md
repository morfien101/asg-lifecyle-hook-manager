# AutoScaling Lifecycle Hook Manager

A simple tool to handle AutoScaling lifecycle hooks.

## Use

```text
This application is used to interact with AWS Autoscaling Lifecycle hooks.
It can set the hooks to Abandon, Continue or send a heartbeat.
Only a single action can be invoked in a single run.
It will consume credentials from instance roles or ENV vars.
There is no provision for manually feeding in credentials and never will be.

  -a string
        Set the lifecycle hook action. Valid values: ABANDON, CONTINUE, HEARTBEAT
  -h    Show the help menu
  -i string
        instance_id for the EC2 instance
  -l string
        Name of the Lifecycle hook
  -n string
        Name of the autoscaling group
  -v    Show the version
  -verbose
        Will log success statements as well as errors
```

## Available on Docker Hub

[Docker Hub: morfien101/asg-lifecycle-hook-manager:latest](https://hub.docker.com/repository/docker/morfien101/asg-lifecycle-hook-manager)
