package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/discobean/aws-ec2-assign-elastic-ip/assigner"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	var tagKey string
	var tagValue string
	var instanceId string
	var region string

	flag.StringVar(&tagKey, "eiptagkey", "", "EIP Pool Tag Key (required)")
	flag.StringVar(&tagValue, "eiptagvalue", "", "EIP Pool Tag Value (required)")
	flag.StringVar(&instanceId, "instanceid", "", "Instance ID to set (optional, if empty use metadata service)")
	flag.StringVar(&region, "region", "", "AWS Region (optional, if empty use metadata service)")

	flag.Parse()

	if tagKey == "" {
		logrus.Errorf("eiptagkey required")
		os.Exit(5)
	}

	if tagValue == "" {
		logrus.Errorf("eiptagvalue required")
		os.Exit(6)
	}

	awsSession, err := session.NewSession(&aws.Config{})
	if err != nil {
		logrus.Errorf("Failed to create AWS session: %v", err)
		os.Exit(1)
	}

	// get the instance ID information if not specified, using the metadata service
	if instanceId == "" {
		metadataSvc := ec2metadata.New(awsSession)
		if !metadataSvc.Available() {
			logrus.Error("No instance metadata available")
			os.Exit(2)
		}

		instanceIdentity, err := metadataSvc.GetInstanceIdentityDocument()
		if err != nil {
			logrus.Errorf("Failed to get instance identity document: %v", err)
			os.Exit(3)
		}

		logrus.Debugf("Got instance ID: %v", instanceIdentity.InstanceID)
		instanceId = instanceIdentity.InstanceID
	}

	// get the region information if not specified, using the metadata service
	if region == "" {
		metadataSvc := ec2metadata.New(awsSession)
		if !metadataSvc.Available() {
			logrus.Error("No instance metadata available")
			os.Exit(2)
		}

		instanceIdentity, err := metadataSvc.GetInstanceIdentityDocument()
		if err != nil {
			logrus.Errorf("Failed to get instance identity document: %v", err)
			os.Exit(3)
		}

		logrus.Debugf("Got region: %v", instanceIdentity.Region)
		region = instanceIdentity.Region
	}

	// create a new session but specify the correct region information now
	awsSession, err = session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		logrus.Errorf("Failed to create AWS session with region: %v", err)
		os.Exit(7)
	}

	assignerSvc, err := assigner.New(awsSession)
	if err != nil {
		logrus.Errorf("Failed to create new EIP assigner: %v", err)
		os.Exit(8)
	}

	err = assignerSvc.AssignEIPFromPoolUsingTags(instanceId, tagKey, tagValue)
	if err != nil {
		logrus.Errorf("Association failed: %v", err)
		os.Exit(9)
	}
}

