# ECSSH 

![ecssh cli demo](res/demo.gif)


Elastic Container (Service) SSH - Allows you to easily navigate running containers in your ECS clusters and run a shell in them using just one simple command.

# Requirements

* AWS CLI - https://github.com/aws/aws-cli
* EC2 Instances must be in a public subnet have ssh running on open port 22.
    * Connecting to an instance in a private subnet with port 22 closed will be possible in future with the addition of AWS SSM support. https://docs.aws.amazon.com/systems-manager/index.html  

### Usage

```shell
ecssh --region "<region>" --cluster "<cluster name or ARN>"
```

Or via Docker

```shell
docker run --rm -it \
  -v ~/.aws:/app/.aws \ # Mounts your AWS config
  -v ~/<key.pem>:/app/<key.pem> \ # The pem file to connect to the instance via SSH 
  blwsh/ecssh:latest -i /app/<key.pem> --region <region>
```

*Note in future versions, support for forwarding the host SSH agent and AWS SSM will be included meaning you can omit mounting a .pem file*

### Flags

| Name    | Optional  | Description                           |
|---------|-----------|---------------------------------------|
| region  | ✔️         | The AWS region for the cluster(s).    |
| cluster | ✔️         | The cluster to search for containers. |

### Future releases

* Option to change default entrypoint command for containers.
* Automate adding and removing security groups for port 22 access for specific IP addresses.
* Automate adding ssh keys to ec2 instances.
* See if fargate support is at all possible.
