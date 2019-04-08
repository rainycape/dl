#ifdef __APPLE__
#define SYMBOL(s) _##s
#else
#define SYMBOL(s) s
#endif

.text
.globl SYMBOL(make_call)

/*
    Received arguments are
    %rdi: function pointer
    %rsi: arguments for registers, always 6
    %rdx: float arguments, either NULL or 8
    %rcx: stack arguments count, already rounded to n*2
    %r8: stack arguments
    %r9: wheter the function pointer returns in %xmm0
*/

SYMBOL(make_call):
    pushq %rbp         # save the base pointer
    movq %rsp, %rbp    # set new base pointer
    pushq %rbx
    pushq %r12
    pushq %r13
    pushq %r14
    pushq %r15

    // When this function was called, return address (8 bytes) was pushed on the stack.
    // Therefore the stack is already un-aligned.
    // We've pushed 6 values (6 * 8 bytes), so must subtract 8 more bytes to
    // get %rsp aligned to a 16 byte boundary.
    // We use this 8 byte space to store the %rcx value (stack arguments count)
    // available with '-40(%rbp)'
    sub $8, %rsp
    movq %rcx, -40(%rbp)        # Save the stack arguments count

    // Can be used to test if %rsp is properly aligned
    // This alligns properly, so if %rsp changes, it was not already aligned
    //andq $~0xF, %rsp

    movq %rdi, %r12    # Save function pointer
    // Save wheter we want to return from %xmm0
    movq %r9, %r13

    // Float arguments, test for no arguments first
    xor %rax, %rax
    test %rdx, %rdx
    je setup_stack
    mov $8, %rax
    movsd (%rdx), %xmm0
    movsd 8(%rdx), %xmm1
    movsd 16(%rdx), %xmm2
    movsd 24(%rdx), %xmm3
    movsd 32(%rdx), %xmm4
    movsd 40(%rdx), %xmm5
    movsd 48(%rdx), %xmm6
    movsd 56(%rdx), %xmm7

setup_stack:
    test %rcx, %rcx
    je stack_done
    pushq (%r8)
    add $8, %r8
    dec %rcx
    jmp setup_stack

stack_done:
    // Integer arguments, put the list in %r14
    // and move to registers in order.
    movq %rsi, %r14
    movq (%r14), %rdi
    movq 8(%r14), %rsi
    movq 16(%r14), %rdx
    movq 24(%r14), %rcx
    movq 32(%r14), %r8
    movq 40(%r14), %r9

    call *%r12             # Call the function
    movq -40(%rbp), %rcx   # Restore stack arguments count to %rcx to know how much we should pop from stack

restore_stack:
    test %rcx, %rcx
    je stack_restored
    add $8, %rsp
    dec %rcx
    jmp restore_stack

stack_restored:
    add $8, %rsp           # Restore %rsp stack adjustment for 16 byte boundary
    test %r13, %r13
    je restore_registers
    movq %xmm0, %rax

restore_registers:
    popq %r15
    popq %r14
    popq %r13
    popq %r12
    popq %rbx
    leaveq
    retq
