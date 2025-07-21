# README Template for Conduit Examples

Use this template when creating READMEs for new examples. Replace the placeholder content with your specific example details.

````markdown
# [example_name]

## 🧠 What It Does

[Brief description of what this example demonstrates - 1-2 sentences explaining the core functionality and use case]

## ⚙️ Requirements

- **[Service/Dependency]** - [Installation instructions or link]
- **[API Key/Model]** - [How to obtain, e.g., "Get from platform.openai.com"]
- **Go 1.21+** - For building and running
- **[Memory/Hardware requirements]** - If applicable

## 🚀 How to Run

```bash
# 1. [Setup step 1 - e.g., start service]
[command]

# 2. [Setup step 2 - e.g., install dependencies]  
[command]

# 3. [Setup step 3 - e.g., set environment variables]
export KEY=value

# 4. Run the example
go run main.go
```

## 🧪 Sample Prompts/Usage

[Include 3-5 example commands or prompts that users can try]

- `"[Example prompt 1]"`
- `"[Example prompt 2]"`
- `"[Example prompt 3]"`

## ✅ Sample Output

```bash
[Show actual terminal output that users can expect to see]
[Include timestamps, tool calls, responses, etc.]
[Make this realistic and copy-pasteable]

> [User input]
🧠 [Agent reasoning or tool selection]
⚡ [Tool execution]
✅ [Result]
Agent: [Natural language response]
```

## 🔧 Tools/Components Used

[List the main tools, agents, or components this example uses]

- **[Category]**: `tool1`, `tool2`, `tool3`
- **[Category]**: `component1`, `component2`

## ⚙️ Configuration Options

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `VAR_NAME` | `default_value` | What this controls |
| `ANOTHER_VAR` | `another_default` | Another configuration option |

## 🎯 Key Features

- ✅ **[Feature 1]**: Brief description
- ✅ **[Feature 2]**: Brief description  
- ✅ **[Feature 3]**: Brief description

## 🔍 How It Works

[Optional: Brief explanation of the workflow/architecture]

1. **Step 1** → What happens first
2. **Step 2** → What happens next
3. **Step 3** → Final result

## ⚠️ Troubleshooting

**[Common Issue 1]:**
```bash
# How to check/fix
command_to_check
```

**[Common Issue 2]:**
- Solution or explanation
- Alternative approach

## 📚 Related Examples

- [`related_example1/`](../related_example1) - Brief description
- [`related_example2/`](../related_example2) - Brief description
- [`related_example3/`](../related_example3) - Brief description

## 🚀 Next Steps

After trying this example:

1. [Suggestion for next thing to try]
2. [Another progression path]
3. [Advanced usage or production deployment]
````

## Guidelines

### Required Sections
- 🧠 **What It Does** - Clear, concise purpose
- ⚙️ **Requirements** - All prerequisites listed
- 🚀 **How to Run** - Step-by-step commands
- ✅ **Sample Output** - Realistic terminal output

### Recommended Sections  
- 🧪 **Sample Prompts** - Things users can try
- 🔧 **Tools Used** - Components and capabilities
- ⚠️ **Troubleshooting** - Common issues and solutions
- 📚 **Related Examples** - Cross-references to other examples

### Style Guidelines
- Use emojis for section headers (🧠 🚀 ✅ etc.)
- Include copy-pasteable commands
- Show realistic terminal output
- Link to related examples
- Keep descriptions concise but helpful
- Use code blocks for commands and output

### Visual Elements (Future)
- Add GIFs or screenshots to show the example running
- Include diagrams for complex workflows
- Use badges for different example types
