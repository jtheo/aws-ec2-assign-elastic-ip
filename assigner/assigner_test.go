package assigner

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	awsSession *session.Session
	assigner   *Assigner
)

func TestConnect(t *testing.T) {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		t.Error("Failed to create AWS session", err)
	}

	awsSession = s
}

func TestAssigner(t *testing.T) {
	a, err := New(awsSession)
	if err != nil {
		t.Error("Failed to create Assigner", err)
	}

	assigner = a
}

func TestHasAssociatedAddressFalse(t *testing.T) {
	result, err := assigner.hasAssociatedAddress("i-0c8b0472f75ee4842")
	if err != nil {
		t.Error("Got an error", err)
		return
	}

	if result {
		t.Error("Expected no addresses to be associated")
	}
}

func TestGetUnassociatedAddresses(t *testing.T) {
	result, err := assigner.getUnassociatedAddresses("Application", "minecraft")
	if err != nil {
		t.Error("Got an error", err)
		return
	}

	if len(result) != 1 {
		t.Errorf("Expected only 1 Address result, got %v results", len(result))
	}

	log.Println(result)
}

func TestAssociatingIPAddress(t *testing.T) {
	_, err := assigner.AssignEIPFromPoolUsingTags("i-0f0e97a20a05ce74b", "Application", "minecraft")
	if err != nil {
		t.Errorf("Failed to associate: %v", err)
	}
}
