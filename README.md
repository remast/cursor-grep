# Simple Grep with Cursor.ai

I created this project to learn how to use Cursor.ai to write code.

## How to run

```bash
# Read from files
go run main.go -pattern "fox" sample1.txt sample2.txt

# Read from standard in
cat sample1.txt | go run main.go -pattern "fox"
```
