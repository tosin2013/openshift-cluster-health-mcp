# Remove Example Tokens from Git History

Following GitHub's official guide: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/removing-sensitive-data-from-a-repository

## ⚠️ Important: Backup First

```bash
# Create a backup
git clone --mirror https://github.com/tosin2013/openshift-cluster-health-mcp.git backup-repo
```

---

## Option 1: BFG Repo-Cleaner (Recommended - Fastest)

### Step 1: Install BFG

```bash
# Download BFG
wget https://repo1.maven.org/maven2/com/madgag/bfg/1.14.0/bfg-1.14.0.jar
# Or: brew install bfg (on macOS)
```

### Step 2: Run BFG to Remove the Token

```bash
# Replace the specific token string across all history
java -jar bfg-1.14.0.jar --replace-text remove-tokens.txt

# Or using brew version:
bfg --replace-text remove-tokens.txt
```

The `remove-tokens.txt` file contains:
```
sha256~n49Ls7eKyPdzqXLz8Pf2a5Ml1AgiPFbTExMNsVRjAVQ
```

BFG will replace this with `***REMOVED***` in all commits.

### Step 3: Clean Up

```bash
git reflog expire --expire=now --all
git gc --prune=now --aggressive
```

### Step 4: Force Push

```bash
git push --force --all
git push --force --tags
```

---

## Option 2: git filter-repo (More Control)

### Step 1: Install git filter-repo

```bash
# macOS
brew install git-filter-repo

# Or download from: https://github.com/newren/git-filter-repo
pip3 install git-filter-repo
```

### Step 2: Create Expression File

Create `expressions.txt`:
```
regex:sha256~n49Ls7eKyPdzqXLz8Pf2a5Ml1AgiPFbTExMNsVRjAVQ==>sha256~YOUR_TOKEN_HERE
```

### Step 3: Run git filter-repo

```bash
git filter-repo --replace-text expressions.txt
```

### Step 4: Force Push

```bash
# Re-add remote (filter-repo removes it)
git remote add origin https://github.com/tosin2013/openshift-cluster-health-mcp.git

git push --force --all
git push --force --tags
```

---

## Option 3: Manual git filter-branch (Advanced)

```bash
git filter-branch --tree-filter '
  find . -type f -name "*.md" -exec sed -i "s/sha256~n49Ls7eKyPdzqXLz8Pf2a5Ml1AgiPFbTExMNsVRjAVQ/sha256~YOUR_TOKEN_HERE/g" {} +
' --prune-empty --tag-name-filter cat -- --all

git reflog expire --expire=now --all
git gc --prune=now --aggressive

git push --force --all
git push --force --tags
```

---

## ⚠️ After Force Push

**Important:** All collaborators need to:

```bash
# Don't pull! Re-clone instead
cd ..
rm -rf openshift-cluster-health-mcp
git clone https://github.com/tosin2013/openshift-cluster-health-mcp.git
cd openshift-cluster-health-mcp
```

---

## Verify It Worked

```bash
# Search for the old token in all history
git log --all --full-history -p -S "sha256~n49Ls7eKyPdzqXLz8Pf2a5Ml1AgiPFbTExMNsVRjAVQ"

# Should return no results
```

---

## ✅ Recommended Approach

**Use BFG Repo-Cleaner** - it's the fastest and safest:

1. `wget https://repo1.maven.org/maven2/com/madgag/bfg/1.14.0/bfg-1.14.0.jar`
2. `java -jar bfg-1.14.0.jar --replace-text remove-tokens.txt`
3. `git reflog expire --expire=now --all && git gc --prune=now --aggressive`
4. `git push --force --all`

Done! 🎉

---

## Current Status

- ✅ `remove-tokens.txt` created with the example token
- ⏳ Ready to run BFG or git filter-repo
- ⏳ Waiting for force push

## After Cleaning

The GitHub secret scanner will automatically:
- Close the alert within 24 hours
- No longer flag the old commits (they'll be unreachable)
