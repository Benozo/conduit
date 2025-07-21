#!/bin/bash

# ConduitMCP RAG Use Cases Demo Script
# This script demonstrates various real-world applications of the RAG pipeline

BASE_URL="http://localhost:8090"

echo "🎯 ConduitMCP RAG Use Cases Demonstration"
echo "========================================"

echo ""
echo "📚 Use Case 1: Corporate Knowledge Base"
echo "-------------------------------------------"
echo "Indexing company policy documents..."

# Add HR Policy Document
curl -s -X POST "$BASE_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "# Remote Work Policy\n\n## Eligibility\nAll full-time employees who have completed 6 months of employment are eligible for remote work arrangements.\n\n## Equipment\nCompany will provide necessary equipment including laptop, monitor, and ergonomic chair for home office setup.\n\n## Communication\nRemote workers must be available during core hours 9 AM - 3 PM in their local timezone for meetings and collaboration.\n\n## Performance\nPerformance will be measured by deliverables and outcomes, not hours worked. Weekly check-ins with managers are required.",
    "title": "Remote Work Policy",
    "metadata": {"department": "HR", "type": "policy", "effective_date": "2024-01-01"}
  }' | jq '.message'

# Query the policy
echo ""
echo "💬 Employee Question: 'What equipment will I get for remote work?'"
curl -s -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"What equipment will I get for remote work?"}' | \
  jq -r '.response' | head -n 3

echo ""
echo "🏥 Use Case 2: Healthcare Knowledge Assistant"
echo "-----------------------------------------------"
echo "Indexing medical guidelines..."

# Add Medical Guideline
curl -s -X POST "$BASE_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "# Hypertension Treatment Guidelines\n\n## Diagnosis\nHypertension is diagnosed when systolic BP ≥140 mmHg or diastolic BP ≥90 mmHg on multiple readings.\n\n## First-Line Treatments\n1. ACE inhibitors (e.g., lisinopril 10-20mg daily)\n2. Calcium channel blockers (e.g., amlodipine 5-10mg daily)\n3. Thiazide diuretics (e.g., hydrochlorothiazide 25mg daily)\n\n## Lifestyle Modifications\n- Reduce sodium intake to <2.3g/day\n- Regular aerobic exercise 150 min/week\n- Weight management (BMI 18.5-24.9)\n- Limit alcohol consumption\n\n## Monitoring\nRecheck BP within 1 month of treatment initiation, then every 3 months once stable.",
    "title": "Hypertension Treatment Guidelines",
    "metadata": {"specialty": "cardiology", "type": "clinical_guideline", "version": "2024.1"}
  }' | jq '.message'

# Query medical information
echo ""
echo "💬 Healthcare Provider Question: 'What are first-line treatments for hypertension?'"
curl -s -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"What are first-line treatments for hypertension?"}' | \
  jq -r '.response' | head -n 5

echo ""
echo "🎓 Use Case 3: Educational Content Assistant"
echo "---------------------------------------------"
echo "Indexing learning materials..."

# Add Educational Content
curl -s -X POST "$BASE_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "# Introduction to Machine Learning\n\n## What is Machine Learning?\nMachine Learning (ML) is a subset of artificial intelligence that enables computers to learn and make decisions from data without being explicitly programmed.\n\n## Types of Machine Learning\n\n### Supervised Learning\n- Uses labeled training data\n- Examples: Classification, Regression\n- Algorithms: Linear Regression, Decision Trees, Neural Networks\n\n### Unsupervised Learning\n- Finds patterns in unlabeled data\n- Examples: Clustering, Dimensionality Reduction\n- Algorithms: K-means, PCA, Autoencoders\n\n### Reinforcement Learning\n- Learns through interaction and rewards\n- Examples: Game playing, Robotics\n- Algorithms: Q-Learning, Policy Gradients\n\n## Getting Started\n1. Learn Python programming\n2. Understand statistics and linear algebra\n3. Practice with datasets on Kaggle\n4. Start with scikit-learn library",
    "title": "Introduction to Machine Learning",
    "metadata": {"subject": "computer_science", "level": "beginner", "type": "tutorial"}
  }' | jq '.message'

# Query educational content
echo ""
echo "💬 Student Question: 'What are the different types of machine learning?'"
curl -s -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"What are the different types of machine learning?"}' | \
  jq -r '.response' | head -n 8

echo ""
echo "🛒 Use Case 4: E-commerce Product Assistant"
echo "--------------------------------------------"
echo "Indexing product catalog..."

# Add Product Information
curl -s -X POST "$BASE_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "# Gaming Laptop Pro X1\n\n## Specifications\n- CPU: Intel Core i7-12700H (2.3GHz base, 4.7GHz boost)\n- GPU: NVIDIA RTX 4060 8GB GDDR6\n- RAM: 16GB DDR5-4800\n- Storage: 512GB NVMe SSD\n- Display: 15.6\" 144Hz FHD IPS\n- Price: $1,299\n\n## Features\n- Ray tracing support for realistic lighting\n- DLSS 3.0 for improved performance\n- RGB backlit keyboard\n- Advanced cooling system with dual fans\n- WiFi 6E and Bluetooth 5.2\n\n## Best For\n- Gaming at 1080p high settings\n- Content creation and video editing\n- Software development\n- 3D modeling and CAD work\n\n## Customer Reviews\n- Average rating: 4.5/5 stars\n- Praised for performance and build quality\n- Some concerns about battery life during gaming",
    "title": "Gaming Laptop Pro X1",
    "metadata": {"category": "laptops", "price": 1299, "brand": "TechBrand", "rating": 4.5}
  }' | jq '.message'

# Query product information
echo ""
echo "💬 Customer Question: 'I need a laptop for video editing under $1500. What do you recommend?'"
curl -s -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"I need a laptop for video editing under $1500. What do you recommend?"}' | \
  jq -r '.response' | head -n 5

echo ""
echo "🔬 Use Case 5: Technical Documentation Search"
echo "-----------------------------------------------"
echo "Indexing API documentation..."

# Add API Documentation
curl -s -X POST "$BASE_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "# Authentication API\n\n## OAuth 2.0 Implementation\n\n### Getting Started\n1. Register your application to get client credentials\n2. Redirect users to authorization endpoint\n3. Exchange authorization code for access token\n4. Use access token in API requests\n\n### Authorization Endpoint\n```\nGET /oauth/authorize?\n  client_id=YOUR_CLIENT_ID&\n  response_type=code&\n  redirect_uri=YOUR_REDIRECT_URI&\n  scope=read write\n```\n\n### Token Exchange\n```\nPOST /oauth/token\nContent-Type: application/x-www-form-urlencoded\n\ngrant_type=authorization_code&\ncode=AUTHORIZATION_CODE&\nclient_id=YOUR_CLIENT_ID&\nclient_secret=YOUR_CLIENT_SECRET&\nredirect_uri=YOUR_REDIRECT_URI\n```\n\n### Using Access Tokens\n```\nGET /api/user/profile\nAuthorization: Bearer YOUR_ACCESS_TOKEN\n```\n\n### Error Handling\n- 401 Unauthorized: Invalid or expired token\n- 403 Forbidden: Insufficient scope permissions\n- 400 Bad Request: Invalid request parameters",
    "title": "OAuth 2.0 API Documentation",
    "metadata": {"type": "api_docs", "version": "v2.1", "category": "authentication"}
  }' | jq '.message'

# Query technical documentation
echo ""
echo "💬 Developer Question: 'How do I implement OAuth authentication?'"
curl -s -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"How do I implement OAuth authentication?"}' | \
  jq -r '.response' | head -n 8

echo ""
echo "📊 Use Case 6: Business Intelligence Query"
echo "-------------------------------------------"
echo "Indexing market research report..."

# Add Market Research
curl -s -X POST "$BASE_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "# AI Market Analysis 2024\n\n## Market Size\nThe global AI market is projected to reach $1.8 trillion by 2030, growing at a CAGR of 36.2% from 2024 to 2030.\n\n## Key Trends\n\n### Generative AI Boom\n- ChatGPT and similar models driving enterprise adoption\n- 67% of companies planning GenAI implementations in 2024\n- Investment in LLM infrastructure growing 340% YoY\n\n### Industry Applications\n1. **Healthcare**: AI-powered diagnostics, drug discovery\n2. **Financial Services**: Fraud detection, algorithmic trading\n3. **Manufacturing**: Predictive maintenance, quality control\n4. **Retail**: Personalization, inventory optimization\n\n### Regional Growth\n- North America: 40% market share, led by US tech giants\n- Asia-Pacific: Fastest growing region at 42% CAGR\n- Europe: Focus on AI regulation and ethical AI development\n\n## Investment Landscape\n- Total AI startup funding: $50.1B in 2023\n- Top funded areas: Foundation models, autonomous vehicles, robotics\n- Corporate AI spending expected to reach $154B by 2025",
    "title": "AI Market Analysis 2024",
    "metadata": {"type": "market_research", "year": 2024, "industry": "technology", "analyst": "TechResearch Inc"}
  }' | jq '.message'

# Query market research
echo ""
echo "💬 Business Analyst Question: 'What are the AI market growth projections for 2024?'"
curl -s -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"What are the AI market growth projections for 2024?"}' | \
  jq -r '.response' | head -n 6

echo ""
echo "📈 Final Knowledge Base Statistics"
echo "----------------------------------"
curl -s -X GET "$BASE_URL/stats" | jq '.'

echo ""
echo "✅ Use Cases Demonstration Complete!"
echo ""
echo "🎉 Summary of Demonstrated Use Cases:"
echo "   📚 Corporate Knowledge Base - HR Policies"
echo "   🏥 Healthcare Assistant - Treatment Guidelines"
echo "   🎓 Educational Content - ML Tutorials"
echo "   🛒 E-commerce - Product Recommendations"
echo "   🔬 Technical Documentation - API Guides"
echo "   📊 Business Intelligence - Market Research"
echo ""
echo "💡 Next Steps:"
echo "   - Customize for your specific domain"
echo "   - Add authentication and access controls"
echo "   - Integrate with your existing systems"
echo "   - Scale with multiple RAG instances"
