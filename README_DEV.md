# agentrc Jekyll Site

Local preview:

```bash
bundle install
bundle exec jekyll serve --livereload
```

Publish through GitHub Pages:

1. Copy these files to the root of `github.com/adeelahmad/agentrc`.
2. Commit and push to `master`.
3. In GitHub: Settings → Pages → Source: GitHub Actions.
4. The included workflow publishes the site.
5. `CNAME` is already set to `agentrc.ai`.
