package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SecurityGroup struct {
	Id      string
	Name    string
	VpcName string
	Descr   string
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	resp, err := ec2Client.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		log.Fatalf("Erro ao listar Security Groups: %v", err)
	}

	var securityGRoupsList []SecurityGroup

	var search string
	fmt.Print("Digite o prefixo do Security Group: ")
	fmt.Scanln(&search)

	var sgName string
	for _, sg := range resp.SecurityGroups {

		sgName = aws.ToString(sg.GroupName)

		if strings.Contains(sgName, search) {
			resp, err := ec2Client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{
				VpcIds: []string{aws.ToString(sg.VpcId)},
			})
			if err != nil {
				log.Fatalf("Erro ao descrever VPC: %v", err)
			}
			if len(resp.Vpcs) == 0 {
				fmt.Println("VPC não encontrada")
				return
			}

			vpc := resp.Vpcs[0]

			var vpcName string
			for _, tag := range vpc.Tags {
				if aws.ToString(tag.Key) == "Name" {
					vpcName = aws.ToString(tag.Value)
					break
				}
			}
			securityGRoupsList = append(securityGRoupsList, SecurityGroup{
				Id:      aws.ToString(sg.GroupId),
				Name:    aws.ToString(sg.GroupName),
				VpcName: aws.ToString(&vpcName),
				Descr:   aws.ToString(sg.Description),
			})

		}
	}
	for i, securityGroup := range securityGRoupsList {
		fmt.Printf("%d | VPC: %s | ID: %s | Name: %s | Descr: %s\n", i, securityGroup.VpcName, securityGroup.Id, securityGroup.Name, securityGroup.Descr)
	}

	var userIP, userName, securityGroupId string
	fmt.Print("Digite o IP de origem: ")
	fmt.Scanln(&userIP)
	fmt.Print("Digite o nome do usuário: ")
	fmt.Scanln(&userName)
	fmt.Print("Digite o ID do Security Group: ")
	fmt.Scanln(&securityGroupId)

	_, err = ec2Client.AuthorizeSecurityGroupIngress(context.TODO(), &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: &securityGroupId,
		IpPermissions: []types.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(22),
				ToPort:     aws.Int32(22),
				IpRanges: []types.IpRange{
					{
						CidrIp:      aws.String(userIP),
						Description: aws.String(userName),
					},
				},
			},
		},
	})

	if err != nil {
		log.Fatal("Erro ao criar regra de entrada: %v", err)
	}

	fmt.Println("Regra criada com sucesso!")
}
