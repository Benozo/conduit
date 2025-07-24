# Custom Workflow Example

This example shows how to create custom agent interactions and workflows using the base components.

## What it demonstrates

- **Custom Workflow Creation**: Building workflows with base agents
- **Agent Coordination**: Using coordinators to manage multiple agents
- **Workflow Execution**: Running custom agent sequences
- **Metrics Collection**: Tracking workflow performance

## How to run

```bash
cd /home/engineone/Projects/AI/ConduitMCP/swarmV2/examples/custom_workflow
go run main.go
```

## Expected output

The example will:
1. Create a coordinator and specialist agents
2. Register agents with the coordinator
3. Create a custom workflow with the agents
4. Execute the workflow
5. Display execution metrics and results

## Code structure

- Uses base agents (`src/agents/base/`)
- Creates custom workflow patterns
- Demonstrates agent coordination
- Shows metrics and status tracking
