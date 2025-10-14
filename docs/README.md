# Recreation MCP Server - GitHub Pages Site

This directory contains a simple, single-page website for the Recreation MCP Server project.

## View Locally

Open `index.html` in your browser:

```bash
# From the repository root
open docs/index.html

# Or use a simple HTTP server
cd docs
python3 -m http.server 8000
# Then visit http://localhost:8000
```

## Publishing to GitHub Pages

### Option 1: Via Repository Settings (Recommended)

1. Go to your GitHub repository settings
2. Navigate to **Pages** (under "Code and automation")
3. Under **Source**, select "Deploy from a branch"
4. Under **Branch**, select `main` and `/docs` folder
5. Click **Save**
6. Your site will be available at: `https://adhocteam.github.io/recreation-mcp-server/`

### Option 2: Via GitHub Actions

Create `.github/workflows/pages.yml`:

```yaml
name: Deploy GitHub Pages

on:
  push:
    branches: ["main"]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Setup Pages
        uses: actions/configure-pages@v4
      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: 'docs'
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Customization

The site is contained in a single `index.html` file with embedded CSS. 

To customize:
- Edit text content in the HTML
- Modify colors in the CSS `<style>` section
- Update repository links (currently set to `adhocteam/recreation-mcp-server`)

## Design

The site features:
- Clean, modern design with gradient backgrounds
- Fully responsive (mobile-friendly)
- Fast loading (no external dependencies)
- Accessible HTML structure
- Clear call-to-action buttons
