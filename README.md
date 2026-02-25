
# Download Cleaner

A simple Go program to clean your Downloads folder by moving files into category folders based on file type.

## Features

- Categorizes files by extension.
- Shows a preview before moving files.
- Supports dry-run mode (`-dry-run`) to generate a report without changing files.
- Moves files older than a configurable threshold (`-age`) into `Archive/`.

## Build release ZIPs for non-technical Windows users

This project includes a packaging script that creates ready-to-share ZIP files for both common Windows CPU types:

- `windows-amd64` (most Windows 10/11 PCs)
- `windows-arm64` (Windows on ARM devices)

From the project root in PowerShell:

```powershell
.\scripts\build-release.ps1 -Version v1.0.0
```

Output:

- `dist/download-cleaner-v1.0.0-windows-amd64.zip`
- `dist/download-cleaner-v1.0.0-windows-arm64.zip`

Each ZIP already contains:

- `download-cleaner.exe`
- `Run Download Cleaner.bat` (double-click launcher)
- `How to Use Download Cleaner.txt` (user instructions)

Users do **not** need Go installed.

## Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/marce/download-cleaner.git
   ```

2. Navigate to the project directory:

   ```bash
   cd download-cleaner
   ```

3. Build the application:

   - On Windows:

     ```bash
     go build -o download-cleaner.exe main.go
     ```

   - On macOS/Linux:

     ```bash
     go build -o download-cleaner main.go
     ```

4. Run the application:

   - Windows:

     ```bash
     .\download-cleaner.exe
     ```

   - macOS/Linux:

     ```bash
     ./download-cleaner
     ```

## Flags

- `-path`: path to Downloads folder (default: `~/Downloads`)
- `-age`: age threshold in days for old files (default: `365`)
- `-dry-run`: preview only; no files moved

Examples:

```bash
./download-cleaner -dry-run
./download-cleaner -age 180
./download-cleaner -path "/custom/downloads/path"
```

## Troubleshooting

- Ensure Go is installed and available in your shell (`go version`).
- Confirm you have read/write permissions for the target Downloads folder.
- Use `-dry-run` first to verify planned moves.
- If files are not categorized as expected, check extension mappings in `scanner/scanner.go`.

## Contributing

Contributions are welcome. Open an issue or submit a pull request with your proposed changes.

## Bug Reports

If you find a bug, open an issue with:

- steps to reproduce,
- expected behavior,
- actual behavior,
- and relevant console output.

## Acknowledgements

- [Go Programming Language](https://golang.org/) - The language used to build this application.

