package dex

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"geistdex/config"
)

// Uniswap V2 Factory Contract Address
var uniswapV2Factory = common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")

// Uniswap V2 Factory ABI
var factoryABI, _ = abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"internalType":"address","name":"tokenA","type":"address"},{"internalType":"address","name":"tokenB","type":"address"}],"name":"getPair","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]`))

// Uniswap V2 Pair ABI
var pairABI, _ = abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"}]`))

func GetUniswapReserves(tokenA, tokenB common.Address) (*big.Int, *big.Int, error) {
	client, err := ethclient.Dial(config.Config.RPCURL)
	if err != nil {
		log.Fatal(err)
	}

	// Encode `getPair` function call
	data, err := factoryABI.Pack("getPair", tokenA, tokenB)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack getPair call: %w", err)
	}

	// Call `getPair` on Uniswap Factory
	msg := ethereum.CallMsg{
		To:   &uniswapV2Factory,
		Data: data,
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to call getPair: %w", err)
	}

	// Decode the response
	outputs, err := factoryABI.Unpack("getPair", result)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unpack getPair response: %w", err)
	}
	pairAddress := outputs[0].(common.Address)

	if pairAddress == (common.Address{}) {
		return nil, nil, fmt.Errorf("pair not found")
	}

	// Encode `getReserves` function call
	data, err = pairABI.Pack("getReserves")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack getReserves call: %w", err)
	}

	// Call `getReserves` on the Uniswap pair contract
	msg = ethereum.CallMsg{
		To:   &pairAddress,
		Data: data,
	}
	result, err = client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to call getReserves: %w", err)
	}

	// Decode the response
	outputs, err = pairABI.Unpack("getReserves", result)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unpack getReserves response: %w", err)
	}

	reserve0 := outputs[0].(*big.Int)
	reserve1 := outputs[1].(*big.Int)

	return reserve0, reserve1, nil
}
