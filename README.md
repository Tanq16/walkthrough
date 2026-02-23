<div align="center">
  <img src=".github/assets/logo.svg" alt="Walkthrough Logo" width="200">
  <h1>Walkthrough</h1>

  <a href="https://github.com/tanq16/walkthrough/actions/workflows/release.yml"><img alt="Build Workflow" src="https://github.com/tanq16/walkthrough/actions/workflows/release.yml/badge.svg"></a>&nbsp;<a href="https://hub.docker.com/r/tanq16/walkthrough"><img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/tanq16/walkthrough"></a><br>
  <a href="https://github.com/tanq16/walkthrough/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/tanq16/walkthrough"></a><br><br>
  <a href="#features">Features</a> &bull; <a href="#installation-and-usage">Install & Use</a> &bull; <a href="#tips-and-notes">Tips & Notes</a>
</div>

---

A JSON Canvas viewer and presentation tool. Load JSON Canvas files in a minimal infinite canvas interface, edit nodes, and present with a laser pointer. My motivation for this was my love for Excalidraw and Obsidian's JSON Canvas. Obsidian's JSON canvas is extremely easy to use for presenting information, while Excalidraw's laser pointer is fantastic to draw attention. This app is aimed at bridging that gap and allowing easy canvases for presentations.

## Features

- Renders JSON Canvas files with text, file, link, and group nodes
- Infinite canvas with pan, zoom, and node drag/resize
- Markdown rendering with syntax-highlighted code blocks
- PDF rendering directly inside nodes with scrollable pages
- Laser pointer for presentations (red trail that fades)
- Pan/Move mode for distraction-free navigation
- Link local `.md` files and upload image/PDF attachments
- Draw edges between nodes via anchor dots
- Multi-select nodes with Shift+click or Shift+drag
- Catppuccin Mocha color scheme with 6 canvas color presets
- Auto-save to `data.json`

Here's a screenshot:

![img](/.github/assets/ss.png)

## Installation and Usage

### Docker (Recommended)

```bash
docker run -d -p 8080:8080 -v ./data:/data tanq16/walkthrough:latest
```

### Binary

Download from [releases](https://github.com/tanq16/walkthrough/releases) and run:

```bash
./walkthrough serve --port 8080 --data ./my-canvas-dir
```

### Build from Source

```bash
git clone https://github.com/tanq16/walkthrough
cd walkthrough
make build
./walkthrough serve
```

## Tips and Notes

- Double-click empty canvas to create a new text note
- Double-click a text node to edit its markdown content
- Press `L` to toggle the laser pointer for presentations
- Press `Space` to toggle pan/move mode for free navigation
- Press `Z` / `X` to zoom out / zoom in
- Use `Ctrl/Cmd + scroll` to zoom, plain scroll to pan
- Drag from anchor dots on node edges to draw connections
- `Shift+click` to select multiple nodes, `Shift+drag` for box selection
- Node colors use Catppuccin-mapped presets from the JSON Canvas spec (1-6)
- The `--data` flag sets the working directory for `data.json`, `.md` files, and attachments
