Context
===========

Context 包含生成过程中的所有信息。

.. code-block:: Go

    type Context struct {
        Grammar *Grammar
        context.Context
        HandlerIndex   int
        SymbolStack    *Stack
        Result         *Derivation
        finish         bool
        Storage        *memdb.MemDB

        VisitedEdge map[string]int
        Mode        Mode
    }


- Grammar: context首先包含了语法图，这是生成过程中的路径参考。
- context: 嵌套了context，以实现基本的 context 语义。
- HandlerIndex：用于标记当前遍历到的Handler。
- SymbolStack：存储即将要处理的符号，处理完后就 Pop 栈顶。
- Result：用于存储生成过程中的结果。
- finish：用来标记生成是否结束
- Storage：内存中的db，用于存储结构化的运行时信息，便于后续汇总分析。
- VisitedEdge：目前用来统计语法生成的覆盖率、做路径决策。后续可能集成到Storage中。
- Mode：用来控制当前是扩张模式（达到约束要求）还是收缩模式（迅速停止生成）。