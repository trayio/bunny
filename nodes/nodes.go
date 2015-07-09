package nodes

import (
	"github.com/trayio/bunny/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/trayio/bunny/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/ec2"
)

type Node struct {
	Host string
}

// Collect gathers all running AWS instances with a specific tag.
// Returns hostname tags in array of *Node structs.
func Collect(cfg *aws.Config) ([]*Node, error) {
	svc := ec2.New(cfg)
	nodes := make([]*Node, 0)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:group"),
				Values: []*string{aws.String("rabbits-production")},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		return nodes, err
	}

	if len(resp.Reservations) > 0 {
		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				for _, tag := range instance.Tags {
					if *tag.Key == "hostname" {
						nodes = append(nodes, &Node{Host: *tag.Value})
					}
				}
			}
		}
	}

	return nodes, nil
}
