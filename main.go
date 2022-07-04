package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jtheo/aws-ec2-assign-elastic-ip/assigner"
)

var Version = "v0.0"

func main() {
	var tagKey string
	var tagValue string
	var instanceId string
	var region string
	var ver bool

	flag.StringVar(&tagKey, "tag-name", "", "EIP Pool Tag Key (required)")
	flag.StringVar(&tagValue, "tag-value", "", "EIP Pool Tag Value (required)")
	flag.StringVar(&instanceId, "instanceid", "", "Instance ID to set (optional, if empty use metadata service)")
	flag.StringVar(&region, "region", "", "AWS Region (optional, if empty use metadata service)")
	flag.BoolVar(&ver, "version", false, "show the version")
	flag.Parse()

	if ver {
		fmt.Printf("aws-ec2-assign-elastic-ip version: %s\n", Version)
		return
	}

	if tagKey == "" {
		log.Println("tag-name required")
		flag.Usage()
		os.Exit(5)
	}

	if tagValue == "" {
		log.Println("tag-name required")
		flag.Usage()
		os.Exit(6)
	}

	awsSession, err := session.NewSession(&aws.Config{})
	if err != nil {
		log.Printf("Failed to create AWS session: %v\n", err)
		os.Exit(1)
	}

	// get the instance ID information if not specified, using the metadata service
	if instanceId == "" {
		metadataSvc := ec2metadata.New(awsSession)
		if !metadataSvc.Available() {
			log.Println("No instance metadata available")
			os.Exit(2)
		}

		instanceIdentity, err := metadataSvc.GetInstanceIdentityDocument()
		if err != nil {
			log.Printf("Failed to get instance identity document: %v\n", err)
			os.Exit(3)
		}

		log.Printf("Got instance ID: %v\n", instanceIdentity.InstanceID)
		instanceId = instanceIdentity.InstanceID
	}

	// get the region information if not specified, using the metadata service
	if region == "" {
		metadataSvc := ec2metadata.New(awsSession)
		if !metadataSvc.Available() {
			log.Println("No instance metadata available")
			os.Exit(2)
		}

		instanceIdentity, err := metadataSvc.GetInstanceIdentityDocument()
		if err != nil {
			log.Printf("Failed to get instance identity document: %v\n", err)
			os.Exit(3)
		}

		log.Printf("Got region: %v\n", instanceIdentity.Region)
		region = instanceIdentity.Region
	}

	// create a new session but specify the correct region information now
	awsSession, err = session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Printf("Failed to create AWS session with region: %v\n", err)
		os.Exit(7)
	}

	assignerSvc, err := assigner.New(awsSession)
	if err != nil {
		log.Printf("Failed to create new EIP assigner: %v\n", err)
		os.Exit(8)
	}

	result, err := assignerSvc.AssignEIPFromPoolUsingTags(instanceId, tagKey, tagValue)
	if err != nil {
		log.Printf("Association failed: %v\n", err)
		os.Exit(9)
	}

	fmt.Println(result)
}
