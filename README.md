# ECSSH 

![ecssh cli demo](res/demo.gif)


Elastic Container (Service) SSH - Allows you to easily navigate running containers in your ECS clusters and run a shell in them using just one simple command.

# Requirements

* AWS CLI - https://github.com/aws/aws-cli

# Installation

Check releases page.

### Usage

```shell
ecssh --region "<region>" --cluster "<cluster name or ARN>"
```

### Flags

| Name    | Optional  | Description                           |
|---------|-----------|---------------------------------------|
| region  | ✔️        | The aws region of the cluster(s).     |
| cluster | ✔️        | The cluster to search for containers. |

### Future releases

* Option to change default command.
* Automate adding and removing security groups for port 22 access for specific IP addresses.
* Automate adding ssh keys to ec2 instances.
* See if fargate support is at all possible.
