 package main

 import (
     "flag"
     "github.com/AlecAivazis/survey/v2"
     "github.com/aws/aws-sdk-go/aws"
     "github.com/aws/aws-sdk-go/aws/session"
     "github.com/aws/aws-sdk-go/service/ec2"
     "github.com/aws/aws-sdk-go/service/ecs"
     "log"
     "os"
     "os/exec"
 )

 type Pivot struct {
     task ecs.Task
     container ecs.Container
 }

 func main() {
     var i string
     var region string
     var cluster string

     // Flags
     flag.StringVar(&i, "i", "", ".PEM file location.")
     flag.StringVar(&region, "region", "", "Example: eu-west-2")
     flag.StringVar(&cluster, "cluster", "", "The name of your cluster.")
     flag.Parse()

     sess, err := session.NewSession()

     if err != nil {
         log.Fatal("Error creating session ", err)
     }

     if region == "" {
        survey.AskOne(&survey.Input{Message: "Region",}, &region)
     }

     sess.Config.Region = aws.String(region)
     ecsSvc := ecs.New(sess)

     if cluster == "" {
         clusters, err := ecsSvc.ListClusters(&ecs.ListClustersInput{})
         if err == nil {
             var clusterArnStrings []string
             for _, arn := range clusters.ClusterArns {
                clusterArnStrings = append(clusterArnStrings, *arn)
             }
             survey.AskOne(&survey.Select{Message: "Select a task", Options: clusterArnStrings}, &cluster)
         } else {
            survey.AskOne(&survey.Input{Message: "Cluster",}, &cluster)
         }
     }

     list, _ := ecsSvc.ListTasks(&ecs.ListTasksInput{
         Cluster: &cluster,
         MaxResults: aws.Int64(100),
     })

     tasks, err := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
         Cluster: &cluster,
         Tasks:   list.TaskArns,
     })

     if err != nil {
         log.Fatal(err.Error())
     }

     var tasksStringArr []string
     containersMap := make(map[string]Pivot)
     for _, task := range tasks.Tasks {
         for _, container := range task.Containers {
             tasksStringArr = append(tasksStringArr, *task.Group + ": "+ *container.Name + " (" + *container.ContainerArn + ")")
             containersMap[*task.Group + ": "+ *container.Name + " (" + *container.ContainerArn + ")"] = Pivot{task: *task, container: *container};
         }
     }

     var taskArn string
     survey.AskOne(&survey.Select{Message: "Select a task", Options: tasksStringArr}, &taskArn)

     selected := containersMap[taskArn]

     ecsInstance, err := ecsSvc.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{Cluster: &cluster, ContainerInstances: []*string{selected.task.ContainerInstanceArn}})

     if err != nil {
        log.Fatal(err.Error())
     }

     ec2Svc := ec2.New(sess)
     ec2Instance, _ := ec2Svc.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{ecsInstance.ContainerInstances[0].Ec2InstanceId}})


     //ssmSvc := ssm.New(sess)
     //startSession, err := ssmSvc.StartSession(&ssm.StartSessionInput{Target: ec2Instance.Reservations[0].Instances[0].InstanceId})
     //fmt.Println(startSession)

     command := exec.Command(
         "ssh", "-t", "-i", i,
         "ec2-user@" + *ec2Instance.Reservations[0].Instances[0].PublicDnsName,
         "docker", "exec", "-it",
         *selected.container.RuntimeId,
         "/bin/sh",
     )
     command.Stdout = os.Stdout
     command.Stdin = os.Stdin
     command.Stderr = os.Stderr
     err = command.Start()

     if err == nil {
        command.Wait()
     } else {
         log.Fatal(err.Error())
     }
 }

