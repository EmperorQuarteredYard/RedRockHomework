package cli

import (
	"bufio"
	"fmt"
	"homeworkSystem/backend/internal/models"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type CLIController struct {
	db         *gorm.DB
	exitSignal *bool
}

func NewCLIController(db *gorm.DB, exitSignal *bool) *CLIController {
	return &CLIController{
		db:         db,
		exitSignal: exitSignal,
	}
}

func (c *CLIController) Serve() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("CLI Controller started. Available commands: list users, promote <user_id>, exit")
	for {
		if *c.exitSignal {
			return
		}
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		cmd := parts[0]
		switch cmd {
		case "list":
			if len(parts) < 2 || parts[1] != "users" {
				fmt.Println("Usage: list users")
				continue
			}
			c.listUsers()
		case "promote":
			if len(parts) != 2 {
				fmt.Println("Usage: promote <user_id>")
				continue
			}
			userID, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				fmt.Println("Invalid user ID")
				continue
			}
			c.promoteUser(userID)
		case "exit":
			fmt.Println("Exiting CLI controller...")
			return
		default:
			fmt.Println("Unknown command. Available: list users, promote <user_id>, exit")
		}
	}
}

func (c *CLIController) listUsers() {
	var users []models.User
	if err := c.db.Find(&users).Error; err != nil {
		fmt.Printf("Error listing users: %v\n", err)
		return
	}
	fmt.Printf("%-5s %-15s %-10s %-10s\n", "ID", "Username", "Role", "Department")
	for _, u := range users {
		fmt.Printf("%-5d %-15s %-10s %-10s\n", u.ID, u.Username, u.Role, u.Department)
	}
}

func (c *CLIController) promoteUser(userID uint64) {
	result := c.db.Model(&models.User{}).Where("id = ?", userID).Update("role", models.RoleAdmin)
	if result.Error != nil {
		fmt.Printf("Error promoting user: %v\n", result.Error)
		return
	}
	if result.RowsAffected == 0 {
		fmt.Printf("User with ID %d not found\n", userID)
		return
	}
	fmt.Printf("User %d promoted to admin successfully\n", userID)
}
