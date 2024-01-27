package langchain

const (
	_stuffSummarizationTemplate = `根据以下内容:


"{{.context}}"


{{.prompt}}:`

	_refineSummarizationTemplate = `您的工作是提供一个最终简洁的摘要 
我们已经提供了一个现有的摘要，直到某个特定点: "{{.existing_answer}}"
我们有机会完善现有的摘要（仅在需要时）使用下面的一些更多背景信息。
------------
"{{.context}}"
------------

根据新的背景信息，完善原始摘要
如果背景信息没有用处，则返回原始摘要。

完善后的摘要:`
)
