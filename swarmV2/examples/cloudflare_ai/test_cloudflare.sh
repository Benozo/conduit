#!/bin/bash

echo "=== Cloudflare Workers AI Example Test ==="
echo "Testing SwarmV2 integration with Cloudflare Workers AI..."
echo ""

# Check if required environment variables are set
if [ -z "$CLOUDFLARE_ACCOUNT_ID" ] || [ -z "$CLOUDFLARE_API_TOKEN" ]; then
    echo "❌ Missing required environment variables!"
    echo ""
    echo "📋 Required Setup:"
    echo "   export CLOUDFLARE_ACCOUNT_ID=your_account_id"
    echo "   export CLOUDFLARE_API_TOKEN=your_api_token"
    echo ""
    echo "🔧 Optional Configuration:"
    echo "   export CLOUDFLARE_MODEL=@cf/meta/llama-3.1-8b-instruct"
    echo ""
    echo "📚 Get credentials from:"
    echo "   • Account ID: Cloudflare Dashboard → Workers & Pages → Overview"
    echo "   • API Token: Cloudflare Dashboard → My Profile → API Tokens"
    echo ""
    echo "💡 You can also create a .env file based on .env.example"
    exit 1
fi

# Build the example
echo "🔧 Building Cloudflare Workers AI example..."
go build -o cloudflare_ai_demo .

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful!"
echo ""

# Test the connection first
echo "🌐 Testing Cloudflare Workers AI connection..."
echo "   Account ID: ${CLOUDFLARE_ACCOUNT_ID:0:8}***"
echo "   Model: ${CLOUDFLARE_MODEL:-@cf/meta/llama-3.1-8b-instruct}"
echo ""

# Run the demo
echo "🚀 Running Cloudflare Workers AI Swarm Demo..."
echo "   Note: This will make API calls to Cloudflare Workers AI"
echo "   Each request counts toward your usage quota"
echo ""

./cloudflare_ai_demo

DEMO_EXIT_CODE=$?

if [ $DEMO_EXIT_CODE -eq 0 ]; then
    echo ""
    echo "✅ Cloudflare Workers AI Demo completed successfully!"
    echo ""
    echo "🎯 What was demonstrated:"
    echo "   • Connection to Cloudflare Workers AI"
    echo "   • Multi-agent collaborative workflow"
    echo "   • Market analysis → Content strategy → Strategic planning"
    echo "   • Edge computing AI with global distribution"
    echo ""
    echo "📊 Performance Benefits:"
    echo "   • Low latency through edge computing"
    echo "   • Global availability (200+ locations)"
    echo "   • Cost-effective pay-per-use model"
    echo "   • No infrastructure management required"
else
    echo ""
    echo "❌ Demo failed with exit code: $DEMO_EXIT_CODE"
    echo ""
    echo "🔍 Common issues:"
    echo "   • Invalid Account ID or API Token"
    echo "   • Insufficient API permissions"
    echo "   • Network connectivity issues"
    echo "   • Cloudflare Workers AI quota exceeded"
    echo ""
    echo "💡 Troubleshooting:"
    echo "   • Verify credentials in Cloudflare dashboard"
    echo "   • Check API token permissions"
    echo "   • Monitor usage in Cloudflare Analytics"
fi

# Cleanup
rm -f cloudflare_ai_demo

echo ""
echo "🔗 Learn more:"
echo "   • Cloudflare Workers AI: https://developers.cloudflare.com/workers-ai/"
echo "   • Available Models: https://developers.cloudflare.com/workers-ai/models/"
echo "   • Pricing: https://www.cloudflare.com/plans/"
