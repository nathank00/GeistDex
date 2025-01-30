package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"geistdex/config"
	"geistdex/internal/dex"
	"geistdex/storage"
)

func main() {
	fmt.Println("Starting GeistDex...")

	// Load configuration
	config.LoadConfig()

	// Initialize database
	storage.InitDB()

	// Define the token pair (WETH/USDC)
	tokenA := common.HexToAddress("0xC02aaa39b223FE8D0A0e5C4F27eAD9083C756Cc2") // WETH
	tokenB := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48") // USDC

	// Fetch Uniswap liquidity
	reserveA, reserveB, err := dex.GetUniswapReserves(tokenA, tokenB)
	if err != nil {
		log.Fatalf("Error fetching Uniswap reserves: %v", err)
	}

	fmt.Printf("Uniswap Reserves for WETH/USDC: %s / %s\n", reserveA.String(), reserveB.String())

	// Store in database
	query := `INSERT INTO liquidity (pair, reserve0, reserve1) VALUES (?, ?, ?)`
	_, err = storage.DB.Exec(query, "WETH/USDC", reserveA.String(), reserveB.String())
	if err != nil {
		log.Fatalf("Failed to store liquidity data: %v", err)
	}

	fmt.Println("Liquidity data stored successfully.")
}
