# AgentRC

**AgentRC** stands for **Agent Run Config**.

AgentRC is an open specification for declaring, packaging, securing, and sharing portable AI agents.

The website in this repository is a complete Jekyll site for **agentrc.ai**.

## Public positioning

AgentRC is not a runtime, sandbox, cloud platform, model provider, or agent framework.

AgentRC defines the contract an agent declares:

- what the agent is;
- how it starts;
- what tools, functions, skills, and MCP servers it needs;
- what files, hosts, and secrets it may access;
- what policies govern those accesses;
- how the package is pinned, signed, shared, and inspected.

Compatible runners execute AgentRC packages on their chosen substrate.

## Local preview

```bash
bundle install
bundle exec jekyll serve --livereload
```

## Publish

The repo includes:

- `_config.yml`
- `Gemfile`
- `CNAME` for `agentrc.ai`
- `.github/workflows/pages.yml`
- Jekyll layouts, CSS, docs, spec pages, schemas, grammar, and examples

Push to `main`, then set GitHub Pages to **GitHub Actions**.

## Branding and theme

This site uses the hybrid AgentRC direction: a sober standards-style documentation structure with a subtle holographic technical mark. It supports both dark and light themes, follows the user's system preference on first load, and stores manual theme selection locally in the browser.

