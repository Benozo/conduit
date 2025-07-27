package main

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

func testConnection() {
	ctx := context.Background()

	// Connect to Milvus
	c, err := client.NewGrpcClient(ctx, "localhost:19530")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer c.Close()

	// List collections
	collections, err := c.ListCollections(ctx)
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}

	fmt.Printf("Found %d collections:\n", len(collections))
	for _, coll := range collections {
		fmt.Printf("- %s\n", coll.Name)
	}

	// Check knowledge_base collection specifically
	collectionName := "knowledge_base"
	has, err := c.HasCollection(ctx, collectionName)
	if err != nil {
		log.Fatalf("Failed to check collection: %v", err)
	}

	if !has {
		fmt.Printf("Collection '%s' does not exist!\n", collectionName)
		return
	}

	fmt.Printf("\nCollection '%s' exists! ✅\n", collectionName)

	// Get collection statistics
	stats, err := c.GetCollectionStatistics(ctx, collectionName)
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
	} else {
		fmt.Printf("Collection statistics:\n")
		for key, value := range stats {
			fmt.Printf("- %s: %s\n", key, value)
		}
	}

	// Describe the collection
	collection, err := c.DescribeCollection(ctx, collectionName)
	if err != nil {
		log.Printf("Failed to describe collection: %v", err)
	} else {
		fmt.Printf("\nCollection schema:\n")
		fmt.Printf("- Name: %s\n", collection.Name)
		fmt.Printf("- Fields:\n")
		for _, field := range collection.Schema.Fields {
			fmt.Printf("  - %s (%s)\n", field.Name, field.DataType)
			if field.TypeParams != nil {
				for k, v := range field.TypeParams {
					fmt.Printf("    %s: %s\n", k, v)
				}
			}
		}
	}
}
