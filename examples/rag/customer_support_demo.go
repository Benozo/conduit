package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// CustomerSupportBot demonstrates a specialized RAG use case
// for intelligent customer support with ticket classification and response generation
type CustomerSupportBot struct {
	ragAPIURL string
}

type SupportTicket struct {
	ID          string    `json:"id"`
	Customer    string    `json:"customer"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type SupportResponse struct {
	TicketID         string   `json:"ticket_id"`
	Response         string   `json:"response"`
	Confidence       float64  `json:"confidence"`
	Sources          []string `json:"sources"`
	SuggestedActions []string `json:"suggested_actions"`
	EscalateToHuman  bool     `json:"escalate_to_human"`
}

func NewCustomerSupportBot(ragAPIURL string) *CustomerSupportBot {
	return &CustomerSupportBot{ragAPIURL: ragAPIURL}
}

// IndexKnowledgeBase adds support documentation to the RAG system
func (bot *CustomerSupportBot) IndexKnowledgeBase() error {
	// Common support documents
	supportDocs := []struct {
		title    string
		content  string
		category string
	}{
		{
			title:    "WiFi Connection Troubleshooting",
			category: "technical",
			content: `# WiFi Connection Issues

## Common Solutions
1. **Restart your router** - Unplug for 30 seconds, then plug back in
2. **Check WiFi password** - Ensure you're using the correct network password
3. **Move closer to router** - Weak signal can cause connection issues
4. **Update network drivers** - Download latest drivers from manufacturer
5. **Reset network settings** - Go to Settings > Network > Reset

## Advanced Troubleshooting
- Check for interference from other devices
- Update router firmware
- Use WiFi analyzer to find best channel
- Contact ISP if issues persist

Resolution time: Usually 5-15 minutes`,
		},
		{
			title:    "Account Password Reset",
			category: "account",
			content: `# Password Reset Process

## Self-Service Reset
1. Go to login page and click "Forgot Password"
2. Enter your email address
3. Check email for reset link (may take 5-10 minutes)
4. Click link and create new password
5. Password must be 8+ characters with uppercase, lowercase, number

## If Email Not Received
- Check spam/junk folder
- Verify email address spelling
- Try different email if you have multiple accounts
- Contact support if still no email after 1 hour

## Security Requirements
- Cannot reuse last 5 passwords
- Password expires every 90 days
- Use password manager for security

Resolution time: Immediate for self-service`,
		},
		{
			title:    "Billing and Payment Issues",
			category: "billing",
			content: `# Billing Support

## Payment Problems
1. **Card Declined** - Check with bank, try different payment method
2. **Billing Address** - Ensure address matches card exactly
3. **Subscription Status** - Verify active subscription in account settings
4. **Refund Requests** - Available within 30 days of purchase

## Common Billing Questions
- Charges appear as "TechCorp Services" on statements
- Billing cycle is monthly on signup date
- Pro-rated charges for mid-cycle upgrades
- Automatic renewal 3 days before expiration

## Contacting Billing Team
- Available Mon-Fri 9 AM - 6 PM EST
- Include account email and transaction ID
- Typical response time: 24-48 hours

Resolution time: Simple issues 1-2 hours, complex billing 1-2 business days`,
		},
		{
			title:    "Software Installation Guide",
			category: "technical",
			content: `# Software Installation

## System Requirements
- Windows 10/11 or macOS 10.15+
- 4GB RAM minimum, 8GB recommended
- 2GB free disk space
- Internet connection for activation

## Installation Steps
1. Download installer from official website
2. Run as administrator (Windows) or enter password (Mac)
3. Accept license agreement
4. Choose installation directory
5. Enter license key when prompted
6. Complete installation and restart

## Common Installation Issues
- **Antivirus blocking**: Add exception for installer
- **Insufficient permissions**: Run as administrator
- **Disk space**: Free up at least 2GB
- **License issues**: Verify key is correct and unused

## Post-Installation
- Update to latest version
- Configure preferences
- Import existing data if needed

Resolution time: 15-30 minutes for standard installation`,
		},
	}

	// Index each document
	for _, doc := range supportDocs {
		payload := map[string]interface{}{
			"content": doc.content,
			"title":   doc.title,
			"metadata": map[string]interface{}{
				"category":   doc.category,
				"type":       "support_doc",
				"indexed_at": time.Now().Format(time.RFC3339),
			},
		}

		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(bot.ragAPIURL+"/documents", "application/json",
			strings.NewReader(string(jsonData)))
		if err != nil {
			return fmt.Errorf("failed to index %s: %v", doc.title, err)
		}
		resp.Body.Close()

		fmt.Printf("✅ Indexed: %s\n", doc.title)
		time.Sleep(200 * time.Millisecond) // Prevent overwhelming the system
	}

	return nil
}

// ProcessTicket analyzes a support ticket and generates response recommendations
func (bot *CustomerSupportBot) ProcessTicket(ticket SupportTicket) (*SupportResponse, error) {
	// Create query from ticket description and subject
	query := fmt.Sprintf("%s %s", ticket.Subject, ticket.Description)

	// Query RAG system
	payload := map[string]interface{}{
		"message": query,
		"limit":   3,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(bot.ragAPIURL+"/chat", "application/json",
		strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to query RAG: %v", err)
	}
	defer resp.Body.Close()

	var chatResponse struct {
		Response string                   `json:"response"`
		Sources  []map[string]interface{} `json:"sources"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract source titles
	var sources []string
	for _, source := range chatResponse.Sources {
		if title, ok := source["source"].(string); ok {
			sources = append(sources, title)
		}
	}

	// Determine if human escalation is needed
	escalateKeywords := []string{"urgent", "emergency", "angry", "escalate", "manager", "legal"}
	escalate := false
	queryLower := strings.ToLower(query)
	for _, keyword := range escalateKeywords {
		if strings.Contains(queryLower, keyword) {
			escalate = true
			break
		}
	}

	// Generate suggested actions based on category
	actions := generateSuggestedActions(ticket.Category, query)

	return &SupportResponse{
		TicketID:         ticket.ID,
		Response:         chatResponse.Response,
		Confidence:       calculateConfidence(len(sources)),
		Sources:          sources,
		SuggestedActions: actions,
		EscalateToHuman:  escalate,
	}, nil
}

func generateSuggestedActions(category, query string) []string {
	switch category {
	case "technical":
		return []string{
			"Send troubleshooting guide",
			"Schedule remote assistance session",
			"Escalate to technical team if issue persists",
		}
	case "billing":
		return []string{
			"Verify account information",
			"Process refund if eligible",
			"Transfer to billing specialist",
		}
	case "account":
		return []string{
			"Verify identity",
			"Send password reset instructions",
			"Update account settings",
		}
	default:
		return []string{
			"Gather more information",
			"Provide relevant documentation",
			"Follow up in 24 hours",
		}
	}
}

func calculateConfidence(sourceCount int) float64 {
	switch {
	case sourceCount >= 3:
		return 0.9
	case sourceCount >= 2:
		return 0.7
	case sourceCount >= 1:
		return 0.5
	default:
		return 0.2
	}
}

// DemoCustomerSupportBot function to show the customer support bot in action
func DemoCustomerSupportBot() {
	fmt.Println("🎧 Customer Support Bot - RAG Use Case Demo")
	fmt.Println("============================================")

	bot := NewCustomerSupportBot("http://localhost:8090")

	// Index support knowledge base
	fmt.Println("\n📚 Indexing support knowledge base...")
	if err := bot.IndexKnowledgeBase(); err != nil {
		log.Fatalf("Failed to index knowledge base: %v", err)
	}

	// Simulate incoming support tickets
	tickets := []SupportTicket{
		{
			ID:          "TICKET-001",
			Customer:    "john.doe@email.com",
			Subject:     "Can't connect to WiFi",
			Description: "My laptop won't connect to the office WiFi. It keeps asking for password but says it's incorrect.",
			Category:    "technical",
			Priority:    "medium",
			Status:      "open",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "TICKET-002",
			Customer:    "sarah.smith@email.com",
			Subject:     "Forgot my password",
			Description: "I can't remember my account password and the reset email isn't coming through.",
			Category:    "account",
			Priority:    "low",
			Status:      "open",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "TICKET-003",
			Customer:    "mike.wilson@email.com",
			Subject:     "Billing question - charged twice",
			Description: "I see two charges on my credit card for this month. This is urgent, please help!",
			Category:    "billing",
			Priority:    "high",
			Status:      "open",
			CreatedAt:   time.Now(),
		},
	}

	// Process each ticket
	fmt.Println("\n🎫 Processing support tickets...")
	for i, ticket := range tickets {
		fmt.Printf("\n--- Ticket %d: %s ---\n", i+1, ticket.Subject)
		fmt.Printf("Customer: %s\n", ticket.Customer)
		fmt.Printf("Description: %s\n", ticket.Description)

		response, err := bot.ProcessTicket(ticket)
		if err != nil {
			fmt.Printf("❌ Error processing ticket: %v\n", err)
			continue
		}

		fmt.Printf("\n🤖 AI Response (Confidence: %.1f):\n", response.Confidence)
		fmt.Printf("%s\n", response.Response[:min(200, len(response.Response))]+"...")

		fmt.Printf("\n📋 Suggested Actions:\n")
		for _, action := range response.SuggestedActions {
			fmt.Printf("  • %s\n", action)
		}

		if response.EscalateToHuman {
			fmt.Printf("\n⚠️  ESCALATE TO HUMAN AGENT\n")
		}

		fmt.Printf("\n📚 Sources: %v\n", response.Sources)
	}

	fmt.Println("\n✅ Customer Support Bot Demo Complete!")
	fmt.Println("\n💡 This demonstrates how RAG can be used for:")
	fmt.Println("   • Automated ticket triage")
	fmt.Println("   • Intelligent response suggestions")
	fmt.Println("   • Knowledge base search")
	fmt.Println("   • Human escalation decisions")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
