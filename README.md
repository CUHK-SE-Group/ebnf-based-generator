# ebnf-based-generator

整体效果类似于： https://github.com/renatahodovan/grammarinator


### cypher enbf示例

https://opencypher.org/resources/

### 目标

基于EBNF语法定义，生成符合定义的语句（代码），并要求尽可能地暴露可控制的参数。

## 示例

以生成SQL为例，下方是用于演示的表格以及数据

```c
mysql> SELECT * FROM Websites;
+----+--------------+---------------------------+-------+---------+
| id | name         | url                       | alexa | country |
+----+--------------+---------------------------+-------+---------+
| 1  | Google       | https://www.google.cm/    | 1     | USA     |
| 2  | 淘宝          | https://www.taobao.com/   | 13    | CN      |
| 3  | 菜鸟教程      | http://www.runoob.com/    | 4689  | CN      |
| 4  | 微博          | http://weibo.com/         | 20    | CN      |
| 5  | Facebook     | https://www.facebook.com/ | 3     | USA     |
+----+--------------+---------------------------+-------+---------+
```

1. 设置起始符号：SELECT

当前状态是SELECT集，下一个TOKEN可以是 id, name, url, alexa, country, * 中的任意一个。

要求选择的下一个TOKEN可以随机，也可以自行指定选择算法（包括自定义算法，或者硬编码返回数据）

1. 假设当前的生成Query已经是：SELECT url FROM Websites WHERE

则要求WHERE的下一个token一定是之前出现过的（在sql中是这么要求的，在pl里可能是要求后续使用的变量一定是前面定义的，且类型得相同）。

例如生成 `WHERE url="https://www.google.cm"` 则是合法的。其中， `"https://www.google.cm"`

可以是任意从数据库里sample出来的值。

第二点描述的是比较具体的针对SQL的情况，在实际实现的更加通用的结构里，这部分应该暴露出一个接口，让用户自定义自己的constraint逻辑，以及自定义生成的逻辑。

## 设计要点

1. 能接收EBNF文件，使用代码生成，生成对应的状态机代码；生成的状态机代码应该能做到下面的两点
2. 合理设计暴露的接口，使其可配置性尽可能地强（为了适应更多的场景，以及为后续加入神经网络适配）
3. 内部需要维护一个抽象的DAG，保存生成的路径，或约束的路径（可用于后续分析）
