# 🚀 Codescribe – AI‑Powered Commit Messages & Code Annotation CLI

Codescribe is a lightweight, intelligent, and developer‑friendly **Go CLI tool** that enhances your workflow by generating AI‑powered **conventional commit messages** and adding **inline code comments** to any programming language — safely and automatically.

Built for speed. Powered by **Groq Llama models**.

---

# ⭐ Features

### 🔥 AI‑Generated Conventional Commit Messages
- Detects both **staged and unstaged** changes.
- Generates clean, meaningful **Conventional Commit** messages.
- Provides **5 message options** using AI.
- Interactive message selection.
- Automatically commits & pushes.

### 🧠 AI Code Annotator
- Adds comments to **ANY programming language**.
- Copy‑safe: never modifies original code — writes to an output directory.
- Guaranteed:
  - ❌ No rewriting
  - ❌ No formatting changes
  - ❌ No renaming
  - ✔ Adds comments ABOVE the relevant lines
- Handles large files using chunking to avoid token limits.

### ⚡ Powered by Groq
- Uses Llama‑3 or Llama‑4‑Scout models.
- Ultra‑fast inference.
- OpenAI‑compatible API format.

---

# 📦 Installation

Clone the repository:

```sh
git clone https://github.com/Parth-11/Codescribe.git
cd Codescribe
```

Install dependencies:

```sh
go mod tidy
```

Build the CLI binary:

```sh
go build -o codescribe
```

(Optional) Move to system PATH:

```sh
sudo mv codescribe /usr/local/bin/
```

---

# 🔐 Environment Setup

Codescribe uses Groq’s LLMs.  
Create a `.env` file in the project root:

```
GROQ_API_KEY=your_api_key_here
```

Or export manually:

```sh
export GROQ_API_KEY=your_api_key_here
```

(Recommended) Load `.env` automatically by importing `github.com/joho/godotenv`.

---

# 🧑‍💻 Usage

Codescribe provides two major commands:

---

# 1️⃣ **Generate Commit Messages**

```sh
codescribe commit
```

This will:

1. Detect unstaged & staged changes.
2. Stage them automatically.
3. Generate 5 AI‑powered commit messages.
4. Let you select one interactively.
5. Commit & push to your Git repository.

### Example:

```
Detecting changes...
Generating commit messages...

Choose a commit message:
  1. feat: add annotation logic with chunking support
  2. fix: handle empty Groq response safely
  3. refactor: improve git diff detection
  4. docs: update README for installation
  5. perf: optimize comment chunk merging
```

---

# 2️⃣ **Annotate Code With AI Comments**

Command syntax:

```sh
codescribe comment --src <source_folder> --out <output_folder>
```

Example:

```sh
codescribe comment --src ./cmd --out ./annotated
```

This will:

1. Copy your codebase to the output folder.
2. Read each file.
3. Split into safe chunks (<= 12KB).
4. Add comments ABOVE lines while preserving original code.
5. Reassemble & save annotated files.

Your original project remains untouched.

---

# 📂 Project Structure

```
Codescribe/
├── cmd/
│   ├── commit.go
│   ├── comment.go
│   └── root.go
├── internal/
│   ├── ai/
│   │   ├── groq.go
│   │   └── comment.go
│   ├── git/
│   │   └── git.go
│   ├── fs/
│   │   └── copy.go
│   └── prompt/
│       └── prompt.go
├── .env (optional)
├── go.mod
├── main.go
└── README.md
```

---

# 🧪 Troubleshooting

### ❌ “missing GROQ_API_KEY”
You must export your API key or load `.env`.

### ❌ “empty Groq response”
Enable debug:

```
export CODESCRIBE_DEBUG=1
```

Then run again to inspect raw API response.

### ❌ Comments overwrite entire file
Already fixed using strict system instructions ensuring:
- Code is NOT rewritten
- Comments appear above lines only

---


# 📄 License
MIT License — free to use and modify.

---
