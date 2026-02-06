package main

func withTimeoutExample()           { /* context.WithTimeout + select */ }
func withCancellationExample()      { /* context.WithCancel + ticker */ }
func workerPoolWithContextExample() { /* workers che ascoltano ctx.Done() */ }
func pipelineWithContextExample()   { /* stage che esce su ctx.Done() */ }
func withValueExample()             { /* key type-safe + helper getter */ }

func main() {
	withTimeoutExample()
	withCancellationExample()
	workerPoolWithContextExample()
	pipelineWithContextExample()
	withValueExample()
}


