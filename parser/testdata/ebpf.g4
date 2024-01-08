grammar ebpf;

controlFlowGraph: basicBlock+ ;

basicBlock: (instruction '\n')* jumpInstruction; // BPF_JMP, BPF_JMP32

instruction: arithmeticAndJump | loadAndStore;

arithmeticAndJump:
    | arithmeticInstruction // BPF_ALU, BPF_ALU64
    ;

arithmeticInstruction:
    addInstruction
    | subInstruction
    | mulInstruction
    | divInstruction
    | sdivInstruction
    | orInstruction
    | andInstruction
    | lshInstruction
    | rshInstruction
    | negInstruction
    | modInstruction
    | smodInstruction
    | xorInstruction
    | movInstruction
    | movsxInstruction
    | arshInstruction
    | endInstruction  // byte swap instruction
    ;

addInstruction: BPF_ADD source aluInsClass src dst '0' imm;
subInstruction: BPF_SUB source aluInsClass src dst '0' imm;
mulInstruction: BPF_MUL source aluInsClass src dst '0' imm;
divInstruction: BPF_DIV source aluInsClass src dst '0' imm;
sdivInstruction: BPF_SDIV source aluInsClass src dst '1' imm;
orInstruction: BPF_OR source aluInsClass src dst '0' imm;
andInstruction: BPF_AND source aluInsClass src dst '0' imm;
lshInstruction: BPF_LSH source aluInsClass src dst '0' imm;
rshInstruction: BPF_RSH source aluInsClass src dst '0' imm;
negInstruction: BPF_NEG source aluInsClass src dst '0' imm;
modInstruction: BPF_MOD source aluInsClass src dst '0' imm;
smodInstruction: BPF_SMOD source aluInsClass src dst '1' imm;
xorInstruction: BPF_XOR source aluInsClass src dst '0' imm;
movInstruction: BPF_MOV source aluInsClass src dst '0' imm;
movsxInstruction: BPF_MOVSX source aluInsClass src dst offset imm;
arshInstruction: BPF_ARSH source aluInsClass src dst '0' imm;
endInstruction: BPF_END swapsource aluInsClass src dst '0' imm; // https://docs.kernel.org/bpf/standardization/instruction-set.html#byte-swap-instructions

movsxOffset: '8'| '16' | '32';
jumpInstruction:
    jaInstruction
    | BPF_JEQ source jmpInsClass src dst offset imm
    | BPF_JGT source jmpInsClass src dst offset imm
    | BPF_JGE source jmpInsClass src dst offset imm
    | BPF_JSET source jmpInsClass src dst offset imm
    | BPF_JNE source jmpInsClass src dst offset imm
    | BPF_JSGT source jmpInsClass src dst offset imm
    | BPF_JSGE source jmpInsClass src dst offset imm
    | callInstruction
    | BPF_EXIT source jmpInsClass '0x0' dst offset imm
    | BPF_JLT source jmpInsClass src dst offset imm
    | BPF_JLE source jmpInsClass src dst offset imm
    | BPF_JSLT source jmpInsClass src dst offset imm
    | BPF_JSLE source jmpInsClass src dst offset imm;

jaInstruction:
    BPF_JA source BPF_JMP '0x0' dst offset imm
    | BPF_JA source BPF_JMP32 '0x0' dst offset imm;

callInstruction:
    BPF_CALL '0x0' jmpInsClass '0x0' dst offset imm
    | BPF_CALL '0x1' jmpInsClass '0x1' dst offset imm
    | BPF_CALL '0x2' jmpInsClass '0x2' dst offset imm
    ;



loadAndStore:
    regularLoadAndStore
    | atomicOperations
    | signExtensionLoadOperations
    | imm64bitOperations;

regularLoadAndStore: BPF_MEM size ldstClass src dst offset imm;
signExtensionLoadOperations: BPF_MEMSX size ldInsClass src dst offset imm;
atomicOperations: BPF_ATOMIC BPF_W BPF_STX src dst offset atomicImmChoice
    | BPF_ATOMIC BPF_DW BPF_STX src dst offset atomicImmChoice
    ;
imm64bitOperations: BPF_IMM BPF_DW BPF_LD immsrc dst offset imm;

immsrc: '0x1' | '0x2' | '0x3' | '0x4' | '0x5' | '0x6';

atomicImmChoice: BPF_ADD | BPF_OR | BPF_AND | BPF_XOR | BPF_FETCH | BPF_XCHG | BPF_CMPXCHG;


size: BPF_W
    | BPF_H
    | BPF_B
    | BPF_DW;

aluInsClass: BPF_ALU | BPF_ALU64;
jmpInsClass: BPF_JMP32 | BPF_JMP;
ldstClass: ldInsClass | stInsClass;
ldInsClass: BPF_LD | BPF_LDX;
stInsClass: BPF_ST | BPF_STX;

source: BPF_K | BPF_X;
swapsource: BPF_TO_LE | BPF_TO_BE | Reserved;


src: reg
    | imm
    ;

dst: reg
    ;

offset: NUMBER
    ;

imm: NUMBER
    ;

reg: 'R0' | 'R1' | 'R2' | 'R3' | 'R4' | 'R5' | 'R6' | 'R7' | 'R8' | 'R9' | 'R10';

BPF_ADD : 'BPF_ADD' | '0x00';
BPF_SUB : 'BPF_SUB' | '0x10';
BPF_MUL : 'BPF_MUL' | '0x20';
BPF_DIV : 'BPF_DIV' | '0x30';
BPF_SDIV : 'BPF_SDIV' | '0x30';
BPF_OR : 'BPF_OR' | '0x40';
BPF_AND : 'BPF_AND' | '0x50';
BPF_LSH : 'BPF_LSH' | '0x60';
BPF_RSH : 'BPF_RSH' | '0x70';
BPF_NEG : 'BPF_NEG' | '0x80';
BPF_MOD : 'BPF_MOD' | '0x90';
BPF_SMOD : 'BPF_SMOD' | '0x90';
BPF_XOR : 'BPF_XOR' | '0xa0';
BPF_MOV : 'BPF_MOV' | '0xb0';
BPF_MOVSX : 'BPF_MOVSX' | '0b00';
BPF_ARSH : 'BPF_ARSH' | '0xc0';
BPF_END : 'BPF_END' | '0xd0';

BPF_JA : 'BPF_JA' | '0x0';
BPF_JEQ : 'BPF_JEQ' | '0x1';
BPF_JGT : 'BPF_JGT' | '0x2';
BPF_JGE : 'BPF_JGE' | '0x3';
BPF_JSET : 'BPF_JSET' | '0x4';
BPF_JNE : 'BPF_JNE' | '0x5';
BPF_JSGT : 'BPF_JSGT' | '0x6';
BPF_JSGE : 'BPF_JSGE' | '0x7';
BPF_CALL : 'BPF_CALL' | '0x8';
BPF_EXIT : 'BPF_EXIT' | '0x9';
BPF_JLT : 'BPF_JLT' | '0xa';
BPF_JLE : 'BPF_JLE' | '0xb';
BPF_JSLT : 'BPF_JSLT' | '0xc';
BPF_JSLE : 'BPF_JSLE' | '0xd';

BPF_IMM : 'BPF_IMM' | '0x00';
BPF_MEM : 'BPF_MEM' | '0x60';
BPF_MEMSX : 'BPF_MEMSX' | '0x80';
BPF_ATOMIC : 'BPF_ATOMIC' | '0xc0';

BPF_W : 'BPF_W' | '0x00';
BPF_H : 'BPF_H' | '0x08';
BPF_B : 'BPF_B' | '0x10';
BPF_DW : 'BPF_DW' | '0x18';

BPF_ALU : 'BPF_ALU' | '0x04';
BPF_ALU64 : 'BPF_ALU64' | '0x07';
BPF_JMP32 : 'BPF_JMP32' | '0x06';
BPF_JMP : 'BPF_JMP' | '0x05';
BPF_LD : 'BPF_LD' | '0x00';
BPF_LDX : 'BPF_LDX' | '0x01';
BPF_ST : 'BPF_ST' | '0x02';
BPF_STX : 'BPF_STX' | '0x03';

BPF_K : 'BPF_K' | '0x00';
BPF_X : 'BPF_X' | '0x08';
BPF_TO_LE : 'BPF_TO_LE' | '0x00';
BPF_TO_BE : 'BPF_TO_BE' | '0x08';
BPF_FETCH: 'BPF_FETCH' | '0x01';
BPF_XCHG: 'BPF_XCHG' | '0xe1';
BPF_CMPXCHG: 'BPF_CMPXCHG' | '0xf1';
Reserved: 'Reserved' | '0x00';

NUMBER : [0-9]+ ;