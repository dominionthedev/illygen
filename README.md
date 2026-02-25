<p align="center">
  <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/images/illygenbrandingimage.svg" width="600" height="300" alt="Illygen"/>
</p>

# Illygen

> A Go-based library and runtime for building intelligence systems.

Illygen enables developers to build AI-like systems that can **reason, make decisions, and learn** — without being full AI models. It mimics the concepts used in AI using deterministic, inspectable, resource-light Go machinery.

**It is not a replacement for AI.** It is a lightweight alternative for domains where AI is overkill — embedded systems, edge computing, domain-specific reasoning engines, and smart automation.

---

## Concepts

| Concept | What it is |
|---|---|
| **Node** | A single unit of reasoning. You consult it, it returns a verdict. Like a neuron. |
| **Flow** | A net of connected nodes — the reasoning pipeline. Like a neural network. |
| **Knowledge** | The feed of intelligence. The more, the smarter the system. |
| **Learning** | Training (high impact, offline) and Exploring (low impact, online). |

---

## Quick Start

```go
package main

import (
    "fmt"

    illygen "github.com/leraniode/illygen"
)

func main() {
    // Define nodes as plain Go functions
    profiler := illygen.NewNode("profiler", func(ctx illygen.Context) illygen.Result {
        if ctx.Get("is_programmer") == true {
            return illygen.Result{Next: "go_intro", Confidence: 0.9}
        }
        return illygen.Result{Next: "prog_intro", Confidence: 0.7}
    })

    goIntro := illygen.NewNode("go_intro", func(ctx illygen.Context) illygen.Result {
        return illygen.Result{
            Value:      "Welcome! Here's Go: https://go.dev",
            Confidence: 1.0,
        }
    })

    progIntro := illygen.NewNode("prog_intro", func(ctx illygen.Context) illygen.Result {
        return illygen.Result{
            Value:      "Welcome! Programming is the art of telling computers what to do.",
            Confidence: 1.0,
        }
    })

    // Wire the flow
    flow := illygen.NewFlow().
        Add(profiler).
        Add(goIntro).
        Add(progIntro).
        Link("profiler", "go_intro", 1.0).
        Link("profiler", "prog_intro", 1.0)

    // Run it
    engine := illygen.NewEngine()

    result, err := engine.Run(flow, illygen.Context{
        "is_programmer": true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Println(result.Value)
    // Output: Welcome! Here's Go: https://go.dev
}
```

---

## Contribution

> [!NOTE]
> Contributions are welcome! Please open an issue or submit a pull request.

If you have ideas, suggestions, or want to contribute code, please feel free to:

- Join the discussion in the [Leraniode Discussions](https://github.com/leraniode/illygen/discussions)
- Open an issue for bugs or feature requests
- Submit a pull request with your changes

---

## License

MIT

---

A [leraniode](https://github.com/leraniode) project.
<p align="left">
   <a href="https://github.com/leraniode">
       <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/footer/leraniodeproductbrandimage.png" width="400" alt="Leraniode"/>
   </a>
</p>
