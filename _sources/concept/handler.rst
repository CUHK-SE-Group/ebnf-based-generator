Handler 
=============

在 :ref:`grammar_graph` 这节中描述了语法图中不同的节点类型。针对不同的节点类型，就会有不同的生成策略。例如，对于 ``or`` 节点，生成策略可以是随机选择一个子节点进行下一步生成内容；也可以是根据历史的生成信息来决定到底选择哪一个子节点进行生成。

为了给不同的节点类型尽可能自由地配置生成策略，这里使用了handler来作为生成策略可配置化的底座。下面给出一个具体的例子来说明这个场景。

- **目标**：根据语法图，生成一颗所有叶子结点都是终结符的语法树。
- **输入**：一个非终结符（初始状态）；或一颗仍有叶子结点是非终结符的树（中间状态）。
- **输出**：扩展n（n>=1）个节点，将这些节点变成另一个非终结符或终结符。

通过上述步骤的不断循环，就可以把一个初始节点不断扩张成一颗所有叶子结点都是终结符的树，此时把这棵树的所有叶子结点从左到右地遍历，即可输出生成的结果。其中该结果是完全符合语法文件的要求的。

Handler 的定义如下：

.. code-block:: Go

    type Handler interface {
        Handle(*Chain, *Context, ResponseCallBack)
        HookRoute() []regexp.Regexp
        Name() string
        Type() GrammarType
    }

- ``Handle`` 是这个 handler 真正的处理部分，可以对输入的树进行操作。所有运行时信息（包括输入的树）都被保存在 ``Context`` 中。
- ``HookRoute() []regexp.Regexp`` 是用来指定这个handler被hook在哪个路径上。例如可能在有些时候只想让一部分节点来使用这个handler，此时可以修改HookRoute来控制生效的作用域。
- ``Name() string`` 返回这个handler的名称，不可重复
- ``Type() GrammarType``，用于匹配不同的节点类型

例如，一个随机选择子节点的Or Handler的实现是：

.. code-block:: Go

    type OrHandler struct {
    }

    func (h *OrHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
        cur := ctx.SymbolStack.Top()
        if len(cur.GetSymbols()) == 0 {
            chain.Next(ctx, cb)
            return
        }
        ctx.SymbolStack.Pop()
        idx := rand.Int() % len(cur.GetSymbols())
        ctx.SymbolStack.Push((cur.GetSymbols())[idx])
        ctx.Result.AddEdge(cur, (cur.GetSymbols())[idx])
        ctx.VisitedEdge[GetEdgeID(cur.GetID(), (cur.GetSymbols())[idx].GetID())]++
        chain.Next(ctx, cb)
    }

    func (h *OrHandler) HookRoute() []regexp.Regexp {
        return make([]regexp.Regexp, 0)
    }

    func (h *OrHandler) Name() string {
        return OrHandlerName
    }

    func (h *OrHandler) Type() GrammarType {
        return GrammarOR
    }

由于 Handler 被定义为一个 interface，因此用户可以自由地实现任意符合接口的 handler，除了基础实现外，还可以在原有的 handler 外部再套上一层处理函数，亦或者是实现一个不处理任何数据的纯监控函数，用于统计、监控运行信息等。


.. code-block:: Go 
    type TraceHandler struct {
    }

    func (h *TraceHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
        // 此处可以记录经过后续handler处理之前的状态
        chain.Next(ctx, cb) // 调用后续的handler进行处理
        // 此处可以记录经过后续handler处理之后的状态
    }

    func (h *TraceHandler) HookRoute() []regexp.Regexp {
        return make([]regexp.Regexp, 0)
    }

    func (h *TraceHandler) Name() string {
        return TraceHandlerName
    }

    func (h *TraceHandler) Type() GrammarType {
        return math.MaxInt // 匹配所有节点类型
    }



Handler Chain
***************

在上述的函数签名中，``Handle(chain *Chain, ctx *Context, cb ResponseCallBack)`` 有三个参数。其中 ctx 是本次运行时的所有信息。而 chain 是由 handler 组成的一条链条，ResponseCallBack则是一个函数，用于在生成结束之后调用。

Handler Chain中最核心的就是 ``Next()`` 函数，用于调用下一个handler进行处理。可以发现，在调用Next函数时，会使用一个循环来遍历后续的handler（起点是 ``ctx.HandlerIndex`` ）。如果此时要处理的符号栈是空的，那么意味着已经处理结束了，则标记为整个流程已经结束。然后再将当前待符号的类型与当前的handler能处理的类型进行匹配，如果匹配上了，那么就进入handler的处理函数。如果此时已经遍历完所有的handler，则表示当前小轮次已经结束，但是尚未结束所有的生成过程。

.. code-block:: Go

    // Next is for to handle next handler in the chain
    func (c *Chain) Next(ctx *Context, f ResponseCallBack) {
        for index := ctx.HandlerIndex; index < len(c.Handlers); index++ {
            ctx.HandlerIndex++
            if ctx.SymbolStack.Top() == nil || ctx.SymbolStack.Empty() {
                if !ctx.SymbolStack.Empty() {
                    panic(ctx)
                }
                if ctx.finish {
                    slog.Error("Warning: Symbol queue should not be empty")
                }
                ctx.finish = true
                r := NewResult(ctx)
                f(r)
                return
            }

            // 如果类型符合
            if ctx.SymbolStack.Top().GetType()&c.Handlers[index].Type() != 0 && satisfy(ctx, c.Handlers[index]) {
                c.Handlers[index].Handle(c, ctx, f)
            }
        }
    }

