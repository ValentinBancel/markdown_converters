# Markdown Converters Project

**ALWAYS follow these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

Markdown Converters is a utility project for converting between different markdown formats and processing markdown files. The project may include Python scripts, Node.js tools, or both depending on the specific conversion needs.

## Working Effectively

### Initial Setup and Dependencies
Since this repository structure is minimal, determine the project type first:

1. **Check for project type indicators:**
   - `ls -la` to see what files exist
   - Look for `package.json` (Node.js), `requirements.txt` or `pyproject.toml` (Python)
   - Check for `Makefile`, `setup.py`, or other build files

**Common System Dependencies for Markdown Conversion:**
Many markdown converters require additional system tools. Install these if needed:
- **Pandoc:** `sudo apt-get update && sudo apt-get install -y pandoc` -- takes 1-3 minutes. NEVER CANCEL. Set timeout to 10+ minutes.
- **wkhtmltopdf:** `sudo apt-get install -y wkhtmltopdf` -- for HTML to PDF conversion
- **LaTeX (for PDF generation):** `sudo apt-get install -y texlive-latex-base texlive-fonts-recommended` -- takes 5-15 minutes. NEVER CANCEL. Set timeout to 30+ minutes.

2. **Python-based project setup (if Python files are present):**
   - Install Python 3.8+ if not available: `python3 --version`
   - Create virtual environment: `python3 -m venv venv`
   - Activate virtual environment: `source venv/bin/activate`
   - Install dependencies: `pip install -r requirements.txt` (if file exists)
   - Alternative: `pip install -e .` (if setup.py exists)
   - Alternative: `pip install .` (if pyproject.toml exists)

3. **Node.js-based project setup (if package.json exists):**
   - Install Node.js 16+ if not available: `node --version`
   - Install dependencies: `npm install` -- takes 2-5 minutes typically. NEVER CANCEL. Set timeout to 10+ minutes.
   - Alternative: `yarn install` (if yarn.lock exists)

4. **Mixed or custom project:**
   - Check for README.md or CONTRIBUTING.md for specific instructions
   - Look for shell scripts or Makefiles for build automation

### Building and Testing

**CRITICAL TIMING EXPECTATIONS:**
- Dependency installation: 2-10 minutes depending on project size. NEVER CANCEL.
- Build processes: May take 5-30 minutes for complex conversions. NEVER CANCEL.
- Test suites: Typically 1-5 minutes. NEVER CANCEL. Set timeout to 15+ minutes.

#### Python Projects
- **Run tests:** `python -m pytest` or `python -m unittest discover` -- takes 1-5 minutes. NEVER CANCEL. Set timeout to 15+ minutes.
- **Linting:** `flake8 .` or `pylint .` (if configured)
- **Type checking:** `mypy .` (if mypy is configured)
- **Formatting:** `black .` (if black is configured)

#### Node.js Projects  
- **Build:** `npm run build` -- takes 2-15 minutes. NEVER CANCEL. Set timeout to 30+ minutes.
- **Test:** `npm test` or `npm run test` -- takes 1-5 minutes. NEVER CANCEL. Set timeout to 15+ minutes.
- **Linting:** `npm run lint` (if configured)
- **Formatting:** `npm run format` or `prettier --write .` (if configured)

#### Universal Commands
- **Custom build:** `make` or `make build` (if Makefile exists) -- timing varies widely, 1-45 minutes. NEVER CANCEL. Set timeout to 60+ minutes.
- **Install script:** `./install.sh` or similar scripts (if present) -- timing varies. NEVER CANCEL. Set timeout to 30+ minutes.

### Running the Application

#### Typical markdown converter usage patterns:
1. **CLI tool (most common):**
   - `python convert.py input.md output.html`
   - `node converter.js --input input.md --output output.html`
   - `./markdown-converter --format html input.md`

2. **Batch processing:**
   - `python batch_convert.py --input-dir ./docs --output-dir ./output`
   - `npm run convert-all`

3. **Web service (if applicable):**
   - `python app.py` or `flask run` (Python Flask)
   - `node server.js` or `npm start` (Node.js Express)
   - Default ports typically 3000, 5000, or 8000

## Validation Scenarios

**ALWAYS run these validation scenarios after making changes:**

### Functional Validation
1. **Basic conversion test:**
   - Create a simple test markdown file: `echo "# Test\nThis is a test." > test.md`
   - Run the conversion tool on this file
   - Verify output is generated and contains expected content
   - Clean up: `rm test.md output.*`

2. **Multiple format test (if supported):**
   - Test conversion to HTML: verify HTML tags are present
   - Test conversion to PDF: verify file is generated and opens
   - Test conversion to other formats as supported

3. **Error handling test:**
   - Test with malformed markdown to ensure graceful error handling
   - Test with missing input files to verify error messages

### Code Quality Validation
**ALWAYS run before committing:**
- Linting tools (see Building and Testing section above)
- Format code if auto-formatters are configured
- Run full test suite
- Check for any TODO or FIXME comments that need addressing

## Common File Locations and Structure

### Typical repository structure:
```
.
├── README.md                 # Main documentation
├── requirements.txt          # Python dependencies (if Python)
├── package.json             # Node.js dependencies (if Node.js)
├── src/                     # Source code
├── tests/                   # Test files
├── docs/                    # Documentation
├── examples/                # Example markdown files
├── output/                  # Generated output (usually gitignored)
└── .github/                 # GitHub workflows and templates
```

### Key files to check when making changes:
- **Main converter logic:** Usually in `src/` directory or root
- **Configuration files:** Look for `config.json`, `.env`, or similar
- **Test files:** In `tests/` or `test/` directory
- **Documentation:** README.md, docs/ directory
- **CI/CD:** `.github/workflows/` for GitHub Actions

### Troubleshooting Common Issues

### System Dependencies
**Check for required tools first:**
- **Pandoc:** `which pandoc || echo "pandoc not found - install with: sudo apt-get install -y pandoc"`
- **wkhtmltopdf:** `which wkhtmltopdf || echo "wkhtmltopdf not found - install with: sudo apt-get install -y wkhtmltopdf"`
- **Make:** `which make || echo "make not found - install with: sudo apt-get install -y build-essential"`

### Dependency Issues
- **Python:** If pip install fails, try `pip install --upgrade pip` first
- **Node.js:** If npm install fails, try `rm -rf node_modules package-lock.json && npm install`
- **Permission issues:** May need `sudo` for system-wide installs (avoid if possible, use virtual environments)

### Build Failures
- Check for missing system dependencies (pandoc, wkhtmltopdf, etc.)
- Verify correct Python/Node.js version
- Clear any cached build artifacts

### Runtime Errors
- Verify input file paths are correct and accessible
- Check output directory exists and is writable
- Validate markdown syntax in input files

## Performance Notes

- **Large file processing:** May take several minutes per file. NEVER CANCEL.
- **Batch operations:** Can take 15-60+ minutes for large document sets. NEVER CANCEL. Set timeout to 90+ minutes.
- **PDF generation:** Often the slowest operation, may take 2-10 minutes per document. NEVER CANCEL.

## CI/CD Integration

If `.github/workflows/` exists:
- **ALWAYS ensure CI passes before merging**
- Common workflow files: `ci.yml`, `test.yml`, `build.yml`
- Typical CI duration: 5-15 minutes. NEVER CANCEL GitHub Actions.

**Remember:** This is a conversion tool project - focus on input/output validation and format accuracy when testing changes.

## Common Command Reference

### Quick Repository Assessment
Run these commands in order when first working with the repository:
```bash
ls -la                          # Check repository structure
git --no-pager status          # Check git status
python3 --version              # Verify Python availability  
node --version                 # Verify Node.js availability
which pandoc || echo "pandoc not found"     # Check for pandoc
which make || echo "make not found"         # Check for make
```

### Typical Development Workflow
1. **Setup:** Follow the Initial Setup section above
2. **Make changes:** Edit source files
3. **Test locally:** Run the validation scenarios
4. **Lint/Format:** Run code quality tools
5. **Test conversion:** Create test markdown and verify output
6. **Commit:** Ensure all validations pass

### Emergency Debugging
If conversion fails or produces unexpected output:
1. Check input file validity: `cat input.md | head -20`
2. Verify output directory permissions: `ls -ld output/` 
3. Test with minimal input: `echo "# Test" | [your converter command]`
4. Check for error logs in common locations: `ls -la *.log` or `find . -name "*.log"`