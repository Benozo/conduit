package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/benozo/conduit/lib/rag"
	"github.com/benozo/conduit/lib/rag/database"
	"github.com/benozo/conduit/lib/rag/embeddings"
	"github.com/benozo/conduit/lib/rag/processors"
)

func main() {
	fmt.Println("🏢 ConduitMCP RAG - Real World Business Example")
	fmt.Println("===============================================")
	fmt.Println("Scenario: TechCorp Knowledge Management System")
	fmt.Println("Indexing: Company policies, procedures, and documentation")
	fmt.Println("")

	// Use Ollama by default for this example (easier local setup)
	provider := "ollama"
	if os.Getenv("RAG_PROVIDER") == "openai" {
		provider = "openai"
	}

	// Configuration
	var config *rag.RAGConfig
	if provider == "openai" {
		config = rag.DefaultRAGConfig()
		if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
			config.Embeddings.APIKey = apiKey
		} else {
			log.Fatal("❌ OPENAI_API_KEY required for OpenAI provider")
		}
	} else {
		config = rag.DefaultOllamaRAGConfig()
		if host := os.Getenv("OLLAMA_HOST"); host != "" {
			config.Embeddings.Host = host
		}
	}

	// Optimize for business documents
	config.Chunking.Size = 800             // Good for policy documents
	config.Chunking.Overlap = 150          // Ensure context continuity
	config.Chunking.Strategy = "paragraph" // Respect document structure

	fmt.Printf("📊 Using %s embeddings\n", strings.ToUpper(provider))
	fmt.Printf("🔧 Chunk size: %d chars, Overlap: %d chars\n",
		config.Chunking.Size, config.Chunking.Overlap)

	ctx := context.Background()

	// Initialize RAG system
	fmt.Println("\n🚀 Initializing TechCorp Knowledge Base...")

	vectorDB, err := database.NewPgVectorDB(config.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer vectorDB.Close()

	var embeddingProvider rag.EmbeddingProvider
	if provider == "openai" {
		embeddingProvider = embeddings.NewOpenAIEmbeddings(
			config.Embeddings.APIKey,
			config.Embeddings.Model,
			config.Embeddings.Dimensions,
			config.Embeddings.Timeout,
		)
	} else {
		embeddingProvider = embeddings.NewOllamaEmbeddings(
			config.Embeddings.Host,
			config.Embeddings.Model,
			config.Embeddings.Dimensions,
			config.Embeddings.Timeout,
		)
	}

	// Test embedding connection
	pingCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if err := embeddingProvider.Ping(pingCtx); err != nil {
		log.Fatalf("Embedding provider connection failed: %v", err)
	}

	chunker := processors.NewTextChunker(
		processors.Paragraph,
		config.Chunking.Size,
		config.Chunking.Overlap,
	)

	ragEngine := rag.NewRAGEngine(config, vectorDB, embeddingProvider, chunker)

	// Check existing documents
	stats, err := ragEngine.GetStats(ctx)
	if err != nil {
		log.Printf("Warning: Could not get stats: %v", err)
		stats = make(map[string]interface{})
	}

	existingDocs := getIntValue(stats, "document_count")
	fmt.Printf("📚 Current knowledge base: %d documents\n", existingDocs)

	// Index company documents if knowledge base is empty
	if existingDocs == 0 {
		fmt.Println("\n📝 Indexing TechCorp Company Documents...")
		indexCompanyDocuments(ctx, ragEngine)
	}

	// Demonstrate real-world business queries
	fmt.Println("\n🔍 Real-World Business Query Examples")
	fmt.Println("=====================================")

	businessQueries := []BusinessQuery{
		{
			Scenario:     "HR Manager - New Employee Onboarding",
			Question:     "What are the steps for onboarding a new software engineer?",
			ExpectedInfo: "HR policies, IT setup, training requirements",
		},
		{
			Scenario:     "Project Manager - Remote Work Policy",
			Question:     "What is our company policy on remote work and flexible hours?",
			ExpectedInfo: "Remote work guidelines, approval process",
		},
		{
			Scenario:     "Developer - Security Guidelines",
			Question:     "What security practices should I follow when handling customer data?",
			ExpectedInfo: "Data protection, GDPR compliance, security protocols",
		},
		{
			Scenario:     "Sales Team - Customer Onboarding",
			Question:     "How do we onboard new enterprise customers?",
			ExpectedInfo: "Customer success process, implementation steps",
		},
		{
			Scenario:     "Finance - Expense Policy",
			Question:     "What expenses can I claim and what's the approval process?",
			ExpectedInfo: "Expense categories, limits, approval workflow",
		},
	}

	for i, query := range businessQueries {
		fmt.Printf("\n%d. 👤 %s\n", i+1, query.Scenario)
		fmt.Printf("   ❓ \"%s\"\n", query.Question)

		// Perform semantic search first
		searchCtx, searchCancel := context.WithTimeout(ctx, 15*time.Second)
		searchResults, err := ragEngine.Search(searchCtx, query.Question, 3, nil)
		searchCancel()

		if err != nil {
			fmt.Printf("   ❌ Search failed: %v\n", err)
			continue
		}

		fmt.Printf("   🔍 Found %d relevant documents\n", len(searchResults))
		for j, result := range searchResults {
			fmt.Printf("      %d. Score: %.3f | %s\n",
				j+1, result.Score, truncateString(result.Chunk.Content, 60))
		}

		// Perform RAG query for comprehensive answer
		ragCtx, ragCancel := context.WithTimeout(ctx, 30*time.Second)
		ragResponse, err := ragEngine.Query(ragCtx, query.Question, 5, nil)
		ragCancel()

		if err != nil {
			fmt.Printf("   ❌ RAG query failed: %v\n", err)
			continue
		}

		fmt.Printf("   🤖 AI Answer (Confidence: %.2f):\n", ragResponse.Confidence)
		fmt.Printf("      %s\n", wrapText(ragResponse.Answer, 70, "      "))
		if len(ragResponse.Sources) > 0 {
			fmt.Printf("   📄 Sources: %d documents referenced\n", len(ragResponse.Sources))
		}
	}

	// Demonstrate search with filters
	fmt.Println("\n🎯 Advanced Search with Filters")
	fmt.Println("===============================")

	filterExamples := []FilterExample{
		{
			Description: "HR Policies Only",
			Query:       "vacation time",
			Filters:     map[string]interface{}{"department": "HR"},
		},
		{
			Description: "Security Documents Only",
			Query:       "data protection",
			Filters:     map[string]interface{}{"category": "security"},
		},
		{
			Description: "Recent Policies (2024)",
			Query:       "company policies",
			Filters:     map[string]interface{}{"year": 2024},
		},
	}

	for i, example := range filterExamples {
		fmt.Printf("\n%d. %s\n", i+1, example.Description)
		fmt.Printf("   Query: \"%s\"\n", example.Query)

		searchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		results, err := ragEngine.Search(searchCtx, example.Query, 3, example.Filters)
		cancel()

		if err != nil {
			fmt.Printf("   ❌ Filtered search failed: %v\n", err)
			continue
		}

		fmt.Printf("   📊 Results: %d documents found\n", len(results))
		for j, result := range results {
			fmt.Printf("      %d. %s...\n", j+1, truncateString(result.Chunk.Content, 50))
		}
	}

	// Show final statistics
	fmt.Println("\n📊 Knowledge Base Statistics")
	fmt.Println("============================")

	finalStats, err := ragEngine.GetStats(ctx)
	if err != nil {
		fmt.Printf("❌ Could not retrieve statistics: %v\n", err)
	} else {
		fmt.Printf("📚 Total Documents: %d\n", getIntValue(finalStats, "document_count"))
		fmt.Printf("📄 Total Chunks: %d\n", getIntValue(finalStats, "chunk_count"))
		fmt.Printf("🧠 Embedding Model: %s\n", getStringValue(finalStats, "embedding_model"))
		fmt.Printf("📏 Vector Dimensions: %d\n", getIntValue(finalStats, "embedding_dimensions"))
	}

	fmt.Println("\n🎉 TechCorp Knowledge Management Demo Complete!")
	fmt.Println("\n💼 Business Impact:")
	fmt.Println("   • Instant access to company knowledge")
	fmt.Println("   • Consistent policy interpretation")
	fmt.Println("   • Reduced time searching for information")
	fmt.Println("   • Improved employee onboarding")
	fmt.Println("   • Better compliance and governance")

	fmt.Println("\n🚀 Next Steps:")
	fmt.Println("   • Index your actual company documents")
	fmt.Println("   • Set up automated document updates")
	fmt.Println("   • Train employees on knowledge search")
	fmt.Println("   • Monitor usage and improve content")
}

func indexCompanyDocuments(ctx context.Context, ragEngine rag.RAGEngine) {
	// Sample company documents that would typically exist
	documents := []CompanyDocument{
		{
			Title:      "Employee Handbook 2024",
			Content:    getEmployeeHandbookContent(),
			Category:   "HR",
			Department: "HR",
			Year:       2024,
			Type:       "policy",
		},
		{
			Title:      "Remote Work Policy",
			Content:    getRemoteWorkPolicyContent(),
			Category:   "HR",
			Department: "HR",
			Year:       2024,
			Type:       "policy",
		},
		{
			Title:      "Data Security Guidelines",
			Content:    getDataSecurityContent(),
			Category:   "security",
			Department: "IT",
			Year:       2024,
			Type:       "guideline",
		},
		{
			Title:      "Customer Onboarding Process",
			Content:    getCustomerOnboardingContent(),
			Category:   "sales",
			Department: "Sales",
			Year:       2024,
			Type:       "process",
		},
		{
			Title:      "Expense Reimbursement Policy",
			Content:    getExpensePolicyContent(),
			Category:   "finance",
			Department: "Finance",
			Year:       2024,
			Type:       "policy",
		},
	}

	for i, doc := range documents {
		fmt.Printf("   📄 Indexing: %s (%d/%d)\n", doc.Title, i+1, len(documents))

		metadata := map[string]interface{}{
			"category":   doc.Category,
			"department": doc.Department,
			"year":       doc.Year,
			"type":       doc.Type,
			"indexed_at": time.Now().Format(time.RFC3339),
		}

		indexCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		_, err := ragEngine.IndexContent(indexCtx, doc.Content, doc.Title, "text/plain", metadata)
		cancel()

		if err != nil {
			fmt.Printf("   ❌ Failed to index %s: %v\n", doc.Title, err)
		} else {
			fmt.Printf("   ✅ Successfully indexed %s\n", doc.Title)
		}
	}
}

// Document content functions (realistic business content)

func getEmployeeHandbookContent() string {
	return `# TechCorp Employee Handbook 2024

## Welcome to TechCorp

Welcome to TechCorp! This handbook contains important information about our company policies, procedures, and culture.

## Onboarding Process for New Employees

### Week 1: Getting Started
1. **Day 1**: Complete HR paperwork and receive equipment
2. **Day 2-3**: IT setup including laptop, accounts, and security training
3. **Day 4-5**: Department introduction and role-specific training

### Software Engineers Onboarding
- Complete security awareness training
- Set up development environment
- Review coding standards and practices
- Assign mentor for first 30 days
- Access to internal documentation and repositories

### Required Training
- Company culture and values (2 hours)
- Security awareness (1 hour)
- GDPR and data protection (1 hour)
- Role-specific technical training (varies)

## Work Schedule and Time Off

### Standard Work Hours
- Core hours: 9:00 AM - 5:00 PM
- Flexible start time: 7:00 AM - 10:00 AM
- Lunch break: 1 hour (flexible timing)

### Vacation Policy
- 25 days annual leave for full-time employees
- Additional 5 days after 5 years of service
- Must be approved by direct manager
- Minimum 2 weeks notice required for vacations longer than 5 days

### Sick Leave
- No limit on sick days for genuine illness
- Medical certificate required for absences longer than 3 consecutive days
- Notify manager as soon as possible

## Benefits

### Health Insurance
- Company pays 80% of premium for employee
- Family coverage available with employee contribution
- Dental and vision coverage included

### Professional Development
- $2,000 annual training budget per employee
- Conference attendance encouraged
- Internal training programs available
- Tuition reimbursement for relevant courses

## Performance Reviews

Annual performance reviews are conducted in January. Mid-year check-ins occur in July. Reviews include:
- Goal achievement assessment
- Skill development progress
- Career planning discussion
- Compensation review`
}

func getRemoteWorkPolicyContent() string {
	return `# TechCorp Remote Work Policy

## Overview

TechCorp supports flexible work arrangements including remote work to promote work-life balance and productivity.

## Eligibility

### Who Can Work Remotely
- Employees who have completed 6 months of employment
- Roles that don't require physical presence
- Employees with satisfactory performance ratings
- Approval from direct manager required

### Remote Work Options
1. **Hybrid Remote**: 2-3 days per week from home
2. **Fully Remote**: Permanent remote work arrangement
3. **Temporary Remote**: Short-term arrangements for specific needs

## Application Process

### Steps to Request Remote Work
1. Complete Remote Work Request Form
2. Discuss with direct manager
3. HR review and approval
4. IT equipment assessment
5. Trial period (30 days for new arrangements)

### Required Equipment
- Company laptop with VPN access
- Secure internet connection (minimum 25 Mbps)
- Dedicated workspace
- Noise-canceling headphones for meetings

## Expectations and Guidelines

### Communication Requirements
- Available during core business hours (9 AM - 3 PM local time)
- Daily check-in with team
- Respond to messages within 4 hours during business hours
- Weekly one-on-one with manager

### Productivity Standards
- Same performance expectations as office-based work
- Regular goal setting and progress tracking
- Participation in all required meetings
- Completion of assigned tasks on schedule

### Security Requirements
- Use company-provided VPN at all times
- Secure home office setup
- No family members access to work equipment
- Regular security updates and training

## Flexible Hours Policy

### Core Hours
All employees must be available during core hours: 10 AM - 2 PM company time

### Flexible Start Times
- Earliest start: 7:00 AM
- Latest start: 10:00 AM
- Minimum 8 hours per day
- Consistent schedule preferred

### Time Zone Considerations
- Remote workers may work in different time zones
- Overlap with team time zones required
- Attendance at key meetings mandatory
- Schedule adjustments for global collaboration

## Review and Termination

Remote work arrangements are reviewed quarterly. The company reserves the right to terminate remote work privileges if:
- Performance standards are not met
- Communication expectations are not followed
- Business needs require physical presence
- Security protocols are violated`
}

func getDataSecurityContent() string {
	return `# TechCorp Data Security Guidelines

## Data Protection Overview

TechCorp is committed to protecting customer and company data in accordance with GDPR, SOC 2, and industry best practices.

## Data Classification

### Public Data
- Marketing materials
- Press releases
- Public website content
- No special handling required

### Internal Data
- Employee directories
- Internal policies
- Project documentation
- Requires basic access controls

### Confidential Data
- Customer information
- Financial records
- Strategic plans
- Requires encryption and access logging

### Restricted Data
- Personal identifiable information (PII)
- Payment card information
- Trade secrets
- Requires highest level of protection

## Customer Data Handling

### Collection and Processing
- Collect only necessary data for business purposes
- Obtain explicit consent for data processing
- Document legal basis for data processing
- Implement privacy by design principles

### Storage Requirements
- All customer data must be encrypted at rest
- Use approved cloud storage services only
- Regular backup and disaster recovery testing
- Data retention limits according to policy

### Access Controls
- Role-based access to customer data
- Minimum necessary access principle
- Regular access reviews and updates
- Multi-factor authentication required

### Data Subject Rights
- Right to access personal data
- Right to rectification of incorrect data
- Right to erasure ("right to be forgotten")
- Right to data portability
- Respond to requests within 30 days

## Security Practices for Developers

### Code Security
- Regular security code reviews
- Use of approved libraries and frameworks
- Static code analysis tools
- Dependency vulnerability scanning

### Infrastructure Security
- All systems must use HTTPS/TLS encryption
- Regular security patching schedule
- Network segmentation and firewalls
- Intrusion detection and monitoring

### Development Environment
- Separate development, staging, and production environments
- No production data in development/testing
- Secure credential management
- Code repository access controls

## Incident Response

### Security Incident Types
- Data breaches or unauthorized access
- Malware or ransomware attacks
- Phishing or social engineering attempts
- System vulnerabilities or misconfigurations

### Response Procedures
1. **Immediate Response** (0-1 hour)
   - Contain the incident
   - Notify security team
   - Document all actions taken

2. **Assessment** (1-4 hours)
   - Determine scope and impact
   - Identify affected systems and data
   - Assess risk level

3. **Notification** (4-24 hours)
   - Notify affected customers if required
   - Report to regulatory authorities if necessary
   - Internal stakeholder communication

4. **Recovery** (Ongoing)
   - Restore affected systems
   - Implement additional security measures
   - Post-incident review and lessons learned

## Training and Compliance

### Required Security Training
- Annual security awareness training for all employees
- Role-specific security training for developers and IT staff
- Phishing simulation tests quarterly
- GDPR training for employees handling personal data

### Compliance Monitoring
- Regular security audits and assessments
- Penetration testing annually
- Vulnerability scanning monthly
- Compliance reporting to leadership team`
}

func getCustomerOnboardingContent() string {
	return `# TechCorp Customer Onboarding Process

## Enterprise Customer Success Framework

Our enterprise customer onboarding process ensures successful implementation and adoption of TechCorp solutions.

## Pre-Onboarding Phase

### Sales to Success Handoff
- Complete customer profile and requirements
- Technical architecture review
- Implementation timeline agreement
- Success criteria definition
- Resource allocation planning

### Kickoff Preparation
- Customer Success Manager assignment
- Technical team selection
- Project plan creation
- Communication channel setup
- Documentation preparation

## Onboarding Phases

### Phase 1: Project Kickoff (Week 1)
- Welcome call with customer leadership
- Project team introductions
- Scope and timeline confirmation
- Communication protocols establishment
- Access and security setup

### Phase 2: Technical Setup (Weeks 2-4)
- Environment provisioning
- Integration planning and setup
- Data migration planning
- Custom configuration
- Security and compliance review

### Phase 3: Configuration & Training (Weeks 5-8)
- System configuration based on requirements
- User access and permissions setup
- Admin training sessions
- End-user training programs
- Documentation and resources provision

### Phase 4: Testing & Go-Live (Weeks 9-12)
- User acceptance testing
- Performance and load testing
- Go-live planning and execution
- Post-go-live support
- Success metrics tracking

## Success Metrics and Milestones

### Technical Milestones
- Environment setup completion
- Integration testing passed
- Performance benchmarks met
- Security review completed
- User acceptance testing signed off

### Business Milestones
- User adoption targets (80% within 30 days)
- Performance KPIs achievement
- Customer satisfaction scores (>8/10)
- Support ticket resolution times
- Training completion rates

## Support and Resources

### Customer Success Team
- Dedicated Customer Success Manager
- Technical Account Manager for enterprise clients
- Support team with SLA commitments
- Professional services for complex implementations

### Training Resources
- Online training portal access
- Live training sessions
- Recorded webinars library
- Best practices documentation
- User community forum access

### Ongoing Support Structure
- Regular health checks and reviews
- Quarterly business reviews
- Proactive monitoring and alerts
- Escalation procedures for issues
- Continuous improvement planning

## Implementation Best Practices

### Change Management
- Stakeholder identification and engagement
- Communication plan for organization
- Training schedule coordination
- Go-live support and troubleshooting
- Feedback collection and iteration

### Risk Mitigation
- Regular progress reviews
- Issue escalation procedures
- Contingency planning
- Resource backup plans
- Timeline adjustment processes

## Customer Success Metrics

### Onboarding Success Indicators
- Time to first value (target: 30 days)
- Feature adoption rate (target: 70%)
- User engagement levels
- Support ticket volume trends
- Customer satisfaction scores

### Long-term Success Factors
- Renewal probability
- Expansion opportunity identification
- Reference customer development
- Case study development
- Success story documentation`
}

func getExpensePolicyContent() string {
	return `# TechCorp Expense Reimbursement Policy

## Overview

This policy outlines approved business expenses and the reimbursement process for TechCorp employees.

## Approved Expense Categories

### Travel Expenses
- **Flights**: Economy class for domestic, business class for international trips >8 hours
- **Hotels**: Reasonable business hotels, up to $200/night in major cities
- **Ground Transportation**: Taxis, rideshares, public transit, rental cars
- **Meals**: Reasonable meal costs during business travel
- **Parking**: Airport and hotel parking fees

### Business Meals and Entertainment
- **Client Meals**: Business meals with customers or prospects
- **Team Meals**: Occasional team building or working meals
- **Entertainment**: Client entertainment with business purpose
- **Limits**: $100 per person for meals, $200 per person for entertainment

### Professional Development
- **Training and Conferences**: Registration fees for approved events
- **Books and Subscriptions**: Professional publications and learning resources
- **Certification Fees**: Industry certifications relevant to role
- **Memberships**: Professional association memberships

### Office and Equipment
- **Office Supplies**: Basic supplies for remote workers
- **Equipment**: Small equipment items under $500
- **Software**: Job-related software subscriptions
- **Internet**: Portion of home internet for remote workers ($50/month max)

## Expense Limits and Guidelines

### Daily Limits
- **Meals (Domestic)**: Breakfast $15, Lunch $25, Dinner $50
- **Meals (International)**: 150% of domestic rates
- **Incidental Expenses**: $25 per day for travel-related items
- **Ground Transportation**: Reasonable and necessary costs

### Approval Requirements
- **$0-$500**: Direct manager approval
- **$501-$2,000**: Department head approval
- **$2,001+**: VP and Finance approval
- **International Travel**: Always requires pre-approval

## Reimbursement Process

### Expense Submission
1. Submit expenses within 30 days of incurrence
2. Use company expense management system
3. Include all required receipts and documentation
4. Provide business justification for each expense
5. Route to appropriate approver based on amount

### Required Documentation
- **Receipts**: Original receipts for all expenses over $25
- **Business Purpose**: Clear explanation of business need
- **Attendees**: Names and affiliations for meal expenses
- **Mileage**: Start/end locations and business purpose
- **Foreign Currency**: Original currency amount and exchange rate

### Approval and Payment
- Manager review and approval within 5 business days
- Finance review for compliance and accuracy
- Payment processed within 2 weeks of approval
- Direct deposit to employee bank account
- Monthly expense reports for management review

## Non-Reimbursable Expenses

### Personal Expenses
- Personal portions of travel (hotel stays extended for leisure)
- Personal meals not related to business
- Personal entertainment or shopping
- Alcoholic beverages (except client entertainment)
- Traffic violations and parking tickets

### Excessive or Luxury Expenses
- First-class flights (unless pre-approved for medical reasons)
- Luxury hotel accommodations above policy limits
- Expensive meals without business justification
- Personal services (laundry, gym, spa)
- Excessive taxi/rideshare costs when alternatives available

## Special Situations

### International Travel
- Currency conversion at date of transaction
- VAT recovery where applicable
- Travel insurance coverage provided by company
- Emergency contact and assistance available
- Per diem rates may apply for extended trips

### Remote Work Expenses
- Home office setup allowance: $1,000 per year
- Internet subsidy: $50 per month
- Equipment replacement on as-needed basis
- Ergonomic accessories: up to $500 per year
- Co-working space fees: up to $200 per month

### Training and Development
- Pre-approval required for events over $1,000
- Travel expenses follow standard travel policy
- Conference materials and networking events included
- Professional certification exam fees covered
- Skills training directly related to current role

## Policy Violations

### Common Violations
- Submitting personal expenses
- Exceeding policy limits without approval
- Missing or inadequate documentation
- Late submission of expense reports
- Fraudulent or duplicate submissions

### Consequences
- First violation: Warning and training
- Second violation: Written warning
- Repeated violations: Disciplinary action up to termination
- Fraudulent expenses: Immediate termination and legal action
- Overpayments must be repaid immediately`
}

// Helper types and functions

type CompanyDocument struct {
	Title      string
	Content    string
	Category   string
	Department string
	Year       int
	Type       string
}

type BusinessQuery struct {
	Scenario     string
	Question     string
	ExpectedInfo string
}

type FilterExample struct {
	Description string
	Query       string
	Filters     map[string]interface{}
}

func getIntValue(m map[string]interface{}, key string) int {
	if val, exists := m[key]; exists {
		if intVal, ok := val.(int); ok {
			return intVal
		}
		if floatVal, ok := val.(float64); ok {
			return int(floatVal)
		}
	}
	return 0
}

func getStringValue(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func wrapText(text string, width int, prefix string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var result strings.Builder
	var line strings.Builder

	for _, word := range words {
		if line.Len()+len(word)+1 > width {
			if line.Len() > 0 {
				result.WriteString(prefix + line.String() + "\n")
				line.Reset()
			}
		}
		if line.Len() > 0 {
			line.WriteString(" ")
		}
		line.WriteString(word)
	}

	if line.Len() > 0 {
		result.WriteString(prefix + line.String())
	}

	return result.String()
}
