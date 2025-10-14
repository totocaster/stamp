# Release Setup Instructions

This document contains step-by-step instructions for setting up the complete release pipeline for the Stamp CLI tool.

## ✅ What I've Already Created

1. **`.goreleaser.yml`** - Complete GoReleaser configuration for automated releases
2. **`.github/workflows/release.yml`** - GitHub Actions workflow for automated releases
3. **`.github/workflows/ci.yml`** - CI workflow for testing on every push/PR
4. **`Formula/stamp.rb`** - Homebrew formula template
5. **`.golangci.yml`** - Linting configuration

## 📋 GitHub Setup Tasks

### 1. Create Homebrew Tap Repository

**Go to GitHub and create a new repository:**
- Repository name: `homebrew-tap`
- Description: "Homebrew tap for Stamp CLI and other tools"
- Make it **PUBLIC** (required for Homebrew taps)
- Initialize with README: Yes

**After creating, add the formula:**
1. Create a `Formula` directory in the new repo
2. Copy the `Formula/stamp.rb` file from this repo to `homebrew-tap/Formula/stamp.rb`
3. Commit and push

### 2. Generate GitHub Personal Access Token for Homebrew

**Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)**

Create a new token with:
- Name: `HOMEBREW_TAP_TOKEN`
- Expiration: 90 days (or no expiration)
- Scopes:
  - ✅ `repo` (all)
  - ✅ `workflow`
- Click "Generate token" and **COPY THE TOKEN** (you won't see it again)

### 3. Add Secrets to Main Repository

**Go to your `stamp` repository → Settings → Secrets and variables → Actions**

Add the following secret:
- Name: `HOMEBREW_TAP_TOKEN`
- Value: [Paste the token you copied in step 2]

### 4. (Optional) Add Codecov Integration

If you want code coverage reporting:

1. Go to https://codecov.io
2. Sign in with GitHub
3. Add the `totocaster/stamp` repository
4. Copy the CODECOV_TOKEN
5. Add it as a repository secret:
   - Name: `CODECOV_TOKEN`
   - Value: [Token from Codecov]

## 🚀 How to Create Your First Release

### Step 1: Commit and Push All Changes

```bash
# Add all the new files
git add .

# Commit the changes
git commit -m "feat: Add automated release pipeline with GoReleaser and Homebrew"

# Push to main branch
git push origin main
```

### Step 2: Create and Push Your First Tag

```bash
# Create an annotated tag
git tag -a v0.1.0 -m "Initial release of Stamp CLI

Features:
- Multiple note types (daily, fleeting, voice, analog, monthly, yearly, project)
- Smart counters with persistence
- Clipboard support (macOS)
- Cross-platform support
- Dual command names (stamp and nid)"

# Push the tag to GitHub
git push origin v0.1.0
```

### Step 3: Monitor the Release

1. Go to your repository on GitHub
2. Click on "Actions" tab
3. You should see a "Release" workflow running
4. Wait for it to complete (usually 3-5 minutes)

### Step 4: Verify the Release

Once the workflow completes successfully:

1. **Check GitHub Releases:**
   - Go to https://github.com/totocaster/stamp/releases
   - You should see v0.1.0 with all platform binaries

2. **Check Homebrew Tap:**
   - Go to your `homebrew-tap` repository
   - The `Formula/stamp.rb` should be automatically updated

3. **Test Installation:**
   ```bash
   # Add your tap
   brew tap totocaster/tap

   # Install stamp
   brew install stamp

   # Test it works
   stamp version
   nid daily
   ```

## 📝 Future Releases

For subsequent releases, you only need to:

1. Make your code changes
2. Commit and push to main
3. Create and push a new tag:
   ```bash
   git tag -a v0.2.0 -m "Release notes here"
   git push origin v0.2.0
   ```

The entire process will run automatically!

## 🏷️ Versioning Guidelines

Follow Semantic Versioning (SemVer):
- **MAJOR** (1.0.0): Breaking changes
- **MINOR** (0.1.0): New features, backwards compatible
- **PATCH** (0.1.1): Bug fixes, backwards compatible

Examples:
- `v0.1.0` - Initial release
- `v0.1.1` - Bug fixes
- `v0.2.0` - New features added
- `v1.0.0` - First stable release

## 🔧 Troubleshooting

### If the release workflow fails:

1. Check the Actions tab for error messages
2. Common issues:
   - Missing `HOMEBREW_TAP_TOKEN` secret
   - Homebrew tap repository doesn't exist
   - Formula directory missing in tap repo

### If Homebrew installation fails:

1. Make sure the tap repository is PUBLIC
2. Check that the formula was updated after release
3. Try: `brew update && brew tap totocaster/tap`

## 📦 Manual Release (Backup Method)

If automation fails, you can release manually:

```bash
# Install GoReleaser locally
brew install goreleaser

# Create a release (needs GITHUB_TOKEN)
export GITHUB_TOKEN="your_github_token"
goreleaser release --clean

# Update Homebrew formula manually in your tap repo
```

## 🎉 Success Checklist

- [ ] Created `homebrew-tap` repository on GitHub
- [ ] Added `HOMEBREW_TAP_TOKEN` secret to main repo
- [ ] Committed all release configuration files
- [ ] Pushed first tag (v0.1.0)
- [ ] Release workflow completed successfully
- [ ] Binaries available on GitHub Releases
- [ ] Homebrew formula updated in tap
- [ ] `brew install stamp` works

## 📚 Additional Resources

- [GoReleaser Documentation](https://goreleaser.com/documentation/about/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Semantic Versioning](https://semver.org/)

---

Once you complete the GitHub setup tasks above, your release pipeline will be fully automated!