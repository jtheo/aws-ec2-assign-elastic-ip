package assigner

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"math/rand"
)

type Assigner struct {
	ec2Svc *ec2.EC2
}

func New(s *session.Session) (*Assigner, error) {
	a := &Assigner{}
	a.ec2Svc = ec2.New(s)

	return a, nil
}

func (a *Assigner) hasAssociatedAddress(instanceId string) (bool, error) {
	result, err := a.ec2Svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-id"),
				Values: []*string{ aws.String(instanceId) },
			},
		},
	})

	if err != nil {
		return false, err
	}

	if len(result.Addresses) > 0 {
		return true, nil
	}

	return false, nil
}

func (a *Assigner) getUnassociatedAddresses(key string, value string) ([]*ec2.Address, error) {
	returnAddresses := []*ec2.Address{}

	result, err := a.ec2Svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + key),
				Values: []*string{ aws.String(value) },
			},
		},
	})
	if err != nil {
		return nil, err
	}

	for _, address := range result.Addresses {
		if address.InstanceId != nil { // Skip if already associated to instance
			continue
		}

		if address.NetworkInterfaceId != nil { // Skip if network interface attached
			continue
		}

		// append to the return list of IPs available
		returnAddresses = append(returnAddresses, address)
	}

	return returnAddresses, nil
}

func (a *Assigner) associateAddress(instanceId string, address *ec2.Address) error {
	_, err := a.ec2Svc.AssociateAddress(&ec2.AssociateAddressInput{
		InstanceId: aws.String(instanceId),
		AllowReassociation: aws.Bool(false),
		AllocationId: address.AllocationId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *Assigner) AssignEIPFromPoolUsingTags(instanceId string, key string, value string) error {
	associated, err := a.hasAssociatedAddress(instanceId)
	if err != nil {
		return err
	}

	if associated {
		return fmt.Errorf("instance %v is already associated with a EIP", instanceId)
	}

	addresses, err := a.getUnassociatedAddresses(key, value)
	if err != nil {
		return err
	}

	if len(addresses) <= 0 {
		return fmt.Errorf("no EIPs available with tag of key/value: %v/%v", key, value)
	}

	// len(addresses) > 1 then pick a random one to use from the list
	err = a.associateAddress(instanceId, addresses[rand.Intn(len(addresses))])
	if err != nil {
		return err
	}

	return nil
}
