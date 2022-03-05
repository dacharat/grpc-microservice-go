package handler

import (
	"net/http"

	"github.com/dacharat/grpc-microservice-go/proto/calculator"
	"github.com/dacharat/grpc-microservice-go/proto/machine"
	"github.com/gin-gonic/gin"
)

type handler struct {
	machineClient    machine.MachineClient
	calculatorClient calculator.CalulatorClient
}

func NewHandler(machineClient machine.MachineClient, calculatorClient calculator.CalulatorClient) handler {
	return handler{
		machineClient:    machineClient,
		calculatorClient: calculatorClient,
	}
}

func (h handler) TestHandler(c *gin.Context) {
	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "MUL"})

	runExecute(h.machineClient, &machine.InstructionSet{Instructions: instructions})

	c.Status(http.StatusNoContent)
}

func (h handler) Test2Handler(c *gin.Context) {
	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "MUL"})

	runExecuteStream(h.machineClient, &machine.InstructionSet{Instructions: instructions})

	c.Status(http.StatusNoContent)
}

func (h handler) TestStreamHandler(c *gin.Context) {
	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 100, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "FIB"})

	runServerStreamingExecute(h.machineClient, &machine.InstructionSet{Instructions: instructions})

	c.Status(http.StatusNoContent)
}

func (h handler) TestStream2Handler(c *gin.Context) {
	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 1, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operand: 2, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "MUL"})
	instructions = append(instructions, &machine.Instruction{Operand: 3, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "ADD"})
	instructions = append(instructions, &machine.Instruction{Operator: "FIB"})

	runServerStreamingExecuteStream(h.machineClient, instructions)

	c.Status(http.StatusNoContent)
}

func (h handler) InstructionHandler(c *gin.Context) {
	var req machine.InstructionSet
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	result, err := h.machineClient.Execute(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h handler) CalculatorHandler(c *gin.Context) {
	var req calculator.NumberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	result, err := h.calculatorClient.Add(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
