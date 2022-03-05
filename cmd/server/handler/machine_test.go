package handler_test

import (
	"context"
	"testing"

	"github.com/dacharat/grpc-microservice-go/cmd/server/handler"
	"github.com/dacharat/grpc-microservice-go/proto/machine"
)

func TestExecute(t *testing.T) {
	s := handler.MachineServer{}

	// set up test cases
	instruction_set_1 := []*machine.Instruction{}
	instruction_set_1 = append(instruction_set_1, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instruction_set_1 = append(instruction_set_1, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instruction_set_1 = append(instruction_set_1, &machine.Instruction{Operator: "ADD"})

	instruction_set_2 := []*machine.Instruction{}
	instruction_set_2 = append(instruction_set_2, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instruction_set_2 = append(instruction_set_2, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instruction_set_2 = append(instruction_set_2, &machine.Instruction{Operator: "MUL"})

	tests := []struct {
		instructions []*machine.Instruction
		want         float32
	}{
		{
			instructions: instruction_set_1,
			want:         11,
		},
		{
			instructions: instruction_set_2,
			want:         30,
		},
	}

	for _, tt := range tests {
		req := &machine.InstructionSet{Instructions: tt.instructions}
		resp, err := s.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("ExecuteTest(%v) got unexpected error", tt.instructions)
		}
		if resp.Output != tt.want {
			t.Errorf("got Execute(%v)=%v, wanted %v", tt.instructions, resp.Output, tt.want)
		}
	}
}
