# HTML + Tailwind Landing Page Builder Agent

This example demonstrates an LLM-powered agent that can generate beautiful HTML landing pages with Tailwind CSS from natural language descriptions.

## Features

🎨 **Smart Design Generation**: LLM analyzes user requirements and creates appropriate designs
📱 **Responsive Design**: All generated pages work perfectly on mobile, tablet, and desktop
🎯 **Multiple Page Types**: SaaS, e-commerce, portfolio, and custom landing pages
🚀 **Modern Stack**: HTML5 + Tailwind CSS via CDN
📄 **Ready to Use**: Complete HTML files that work immediately in browsers
⚡ **Fast Generation**: Creates complete landing pages in seconds

## Prerequisites

1. **DeepInfra API Token**: Get your token from [DeepInfra](https://deepinfra.com)
2. **Go Environment**: Go 1.19+ installed
3. **Modern Browser**: For viewing generated HTML previews

## Setup

1. **Set Environment Variables**:
```bash
export DEEPINFRA_TOKEN=your_deepinfra_token_here
export DEEPINFRA_MODEL=meta-llama/Meta-Llama-3.1-8B-Instruct  # optional
```

2. **Run the Example**:
```bash
cd /home/engineone/Downloads/gomcp
go run examples/agents_vue_builder/main.go
```

## What It Does

The agent will automatically generate three different landing pages:

### 1. 🏢 SaaS Analytics Platform
- Modern business dashboard landing page
- Features real-time analytics theme
- Blue/purple gradient design
- Professional layout with CTA buttons

### 2. 🎧 Premium Headphones Product Page  
- E-commerce product showcase
- Premium audio equipment theme
- Black/gold minimalist design
- Product features and pricing

### 3. 👨‍💻 Developer Portfolio
- Personal portfolio landing page
- Full-stack developer theme
- Dark theme with green accents
- Skills showcase and contact info

## Generated Files

All files are saved to `./generated_pages/`:

```
generated_pages/
├── saas-landing.html         # Complete HTML page
├── product-page.html         # Complete HTML page
├── portfolio.html            # Complete HTML page
├── custom-styles.css         # Optional custom CSS
└── README.md                 # Project documentation
```

## How to View Generated Pages

### Simple and Fast (Recommended)
1. Open any `.html` file directly in your browser
2. No setup, no build process, no dependencies required
3. Tailwind CSS is loaded via CDN for instant styling

### For Development
1. Copy the HTML files to your web server
2. Customize the Tailwind classes as needed
3. Add custom CSS for additional styling

## Tools Used by the Agent

The LLM agent has access to these specialized tools:

- **`create_html_page`**: Creates complete HTML5 pages with Tailwind CSS
- **`create_css_file`**: Creates custom CSS files for additional styling
- **`create_project_files`**: Generates project documentation and config files
- **Standard tools**: Memory, text processing, UUID generation

## Example Output

When you run the agent, you'll see clean output focusing on the LLM reasoning and file generation:

```
🤖 Model: meta-llama/Meta-Llama-3.1-8B-Instruct
🧠 Creating SaaS landing page...
📄 CREATED: ./generated_pages/saas-landing.html (4521 bytes)
📋 PROJECT: ./generated_pages/README.md (892 bytes)
✅ SaaS landing page created in 2.3s!

🧠 Creating e-commerce product page...
📄 CREATED: ./generated_pages/product-page.html (3987 bytes)
✅ Product landing page created in 1.8s!

🧠 Creating personal portfolio page...
📄 CREATED: ./generated_pages/portfolio.html (4156 bytes)
✅ Portfolio landing page created in 2.1s!

📁 Generated files saved to: ./generated_pages
🎉 HTML + Tailwind landing page generation completed!
```

## Customization

To create custom landing pages, modify the task parameters in `main.go`:

```go
customTask, _ := llmAgentManager.CreateTask(
    "vue_builder_agent",
    "Your Custom Page Title",
    "Description of what to build",
    map[string]interface{}{
        "project_type": "Custom Landing Page",
        "company":      "Your Company",
        "description":  "What your product/service does",
        "features":     []string{"Feature 1", "Feature 2"},
        "target_audience": "Who this is for",
        "color_scheme": "Your preferred colors",
    },
)
```

## Technologies Used

- **HTML5**: Semantic markup with modern standards
- **Tailwind CSS**: Utility-first CSS framework via CDN
- **DeepInfra LLM**: AI-powered design and content generation
- **Go**: Backend agent orchestration

## Benefits

✅ **Instant Results**: No build process, works immediately in browsers
✅ **Professional Quality**: LLM creates modern, responsive designs
✅ **Zero Dependencies**: Pure HTML + Tailwind CSS via CDN
✅ **Learning Tool**: See how modern HTML + Tailwind work together
✅ **Production Ready**: Generated code follows best practices

## Next Steps

1. **Extend the Agent**: Add more specialized tools for specific industries
2. **Add More Frameworks**: Support for Vue, React, or other frameworks
3. **Enhanced Styling**: Add support for custom animations and interactions
4. **Content Management**: Integration with headless CMS
5. **Deployment**: Automatic deployment to static hosting services

This example showcases the power of LLM-driven development tools for rapid web development with plain HTML and Tailwind CSS!
