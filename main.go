package main

import (
	"fmt"
	"time"

	"github.com/luizfpsoares/dscli/aws"
)

func main() {
	opcao := 1
	for opcao >= 1 {
		fmt.Println("[1] - Listar Security Groups")
		fmt.Println("[2] - Editar regra")
		fmt.Println("[0] - Sair")
		fmt.Scanln(&opcao)

		if opcao == 1 {
			var search string
			fmt.Print("Digite o nome completo ou parte do nome do Security Group: ")
			fmt.Scanln(&search)
			securityGroupList := aws.GetSecurityGroup(search)
			fmt.Println("")
			for i, securityGroup := range securityGroupList {
				fmt.Printf("%d | VPC: %s | ID: %s | Name: %s | Descr: %s\n", i, securityGroup.VpcName, securityGroup.Id, securityGroup.Name, securityGroup.Descr)
			}
			fmt.Println("")
		}
		if opcao == 2 {
			var userName string
			var ip string
			var securityGroupId string
			fmt.Print("Digite o nome do usuário para descrição: ")
			fmt.Scanln(&userName)
			fmt.Print("Digite o IP de origem: ")
			fmt.Scanln(&ip)
			fmt.Print("Digite o ID do Security Group: ")
			fmt.Scanln(&securityGroupId)
			response := aws.AddIngressRule(userName, ip, securityGroupId)
			fmt.Println(response)
			fmt.Println("")
			time.Sleep(3 * time.Second)
		}
	}
}
