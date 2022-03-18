package aws

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	Session *session.Session
)

func init() {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	region := os.Getenv("AWS_REGION")

	if accessKeyID == "" {
		log.Fatal("You need to set AWS_ACCESS_KEY_ID")
	}

	if secretAccessKey == "" {
		log.Fatal("You need to set AWS_SECRET_ACCESS_KEY")
	}

	if sessionToken == "" {
		log.Fatal("You need to set AWS_SESSION_TOKEN")
	}

	if region == "" {
		region = "eu-west-1"
	}

	CreateSessionWithConfig(accessKeyID, secretAccessKey, sessionToken, region)
}

// CreateSessionWithConfig creates a session with the configuration provided
// sessionToken can be empty
func CreateSessionWithConfig(accessKeyID, secretAccessKey, sessionToken, region string) {
	s, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken),
		Region:      aws.String(region),
	})

	if err != nil {
		log.Printf("Could not initiate AWS Session: %v", err)
	}

	Session = s
}

func GetIPs() ([]string, error) {
	var ips []string
	svc := ec2.New(Session)

	input := ec2.DescribeNetworkInterfacesInput{}

	output, err := svc.DescribeNetworkInterfaces(&input)
	if err != nil {
		return ips, err
	}

	if len(output.NetworkInterfaces) < 0 {
		log.Fatal("no IPs")
	}

	for _, v := range output.NetworkInterfaces {
		if v.Association != nil {
			ips = append(ips, *v.Association.PublicIp)
		}
	}

	return ips, nil
}
