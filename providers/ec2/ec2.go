package ec2

//----------------------------------------------------------------------------
// Package factored import statement:
//----------------------------------------------------------------------------

import (

	// Stdlib:
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	// Community:
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//----------------------------------------------------------------------------
// Typedefs:
//----------------------------------------------------------------------------

// Data contains variables used by EC2 API.
type Data struct {
	Region    string
	SubnetIds string
	ImageID   string
	KeyPair   string
	InsType   string
	HostName  string
	ElasticIP string
}

//--------------------------------------------------------------------------
// func: Setup
//--------------------------------------------------------------------------

// Setup an EC2 VPC and all the related components.
func (d *Data) Setup() error {

	// Connect and authenticate to the API endpoint:
	log.Printf("[setup-ec2] INFO Connecting to %s\n", d.Region)
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(d.Region)}))

	// Create an VPC:
	if err := d.createVpc(*svc); err != nil {
		return err
	}

	// Return on success:
	return nil
}

//-------------------------------------------------------------------------
// func: createVpc
//-------------------------------------------------------------------------

func (d *Data) createVpc(svc ec2.EC2) error {

	// Forge the VPC request:
	params := &ec2.CreateVpcInput {
		CidrBlock: aws.String("10.0.0.0/16"),
		DryRun:    aws.Bool(false),
		InstanceTenancy: aws.String("default"),
	}

	// Send the VPC request:
	log.Printf("[setup-ec2] INFO Creating a VPC\n")
	resp, err := svc.CreateVpc(params)
	if err != nil {
		return err
	}

	// Return on success:
	log.Printf("[setup-ec2] INFO VpcId: %s\n", *resp.Vpc.VpcId)
	return nil
}

//--------------------------------------------------------------------------
// func: Run
//--------------------------------------------------------------------------

// Run uses EC2 API to launch a new instance.
func (d *Data) Run(udata []byte) error {

	// Connect and authenticate to the API endpoint:
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(d.Region)}))

	// Forge the network interfaces:
	var networkInterfaces []*ec2.InstanceNetworkInterfaceSpecification
	subnetIds := strings.Split(d.SubnetIds, ",")

	for i := 0; i < len(subnetIds); i++ {

		// Forge the security group ids:
		var securityGroupIds []*string
		for _, gid := range strings.Split(subnetIds[i], ":")[1:] {
			securityGroupIds = append(securityGroupIds, aws.String(gid))
		}

		iface := ec2.InstanceNetworkInterfaceSpecification{
			DeleteOnTermination: aws.Bool(true),
			DeviceIndex:         aws.Int64(int64(i)),
			Groups:              securityGroupIds,
			SubnetId:            aws.String(strings.Split(subnetIds[i], ":")[0]),
		}

		networkInterfaces = append(networkInterfaces, &iface)
	}

	// Send the request:
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:           aws.String(d.ImageID),
		MinCount:          aws.Int64(1),
		MaxCount:          aws.Int64(1),
		KeyName:           aws.String(d.KeyPair),
		InstanceType:      aws.String(d.InsType),
		NetworkInterfaces: networkInterfaces,
		UserData:          aws.String(base64.StdEncoding.EncodeToString([]byte(udata))),
	})

	if err != nil {
		return err
	}

	// Pretty-print the response data:
	fmt.Println("Created instance", *runResult.Instances[0].InstanceId)

	// Add tags to the created instance:
	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(d.HostName),
			},
		},
	})

	if err != nil {
		return err
	}

	// Allocate an elastic IP address:
	if d.ElasticIP == "true" {

		params := &ec2.AllocateAddressInput{
			Domain: aws.String("vpc"),
			DryRun: aws.Bool(false),
		}

		// Send the request:
		resp, err := svc.AllocateAddress(params)

		if err != nil {
			return err
		}

		// Pretty-print the response data:
		fmt.Println(resp)
	}

	// Return on success:
	return nil
}
