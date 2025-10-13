# CI/CD Security & Permissions

## Overview

This document explains the security model and permissions used in our GitHub Actions workflows, following the principle of least privilege.

## Principle of Least Privilege

Each workflow and job is granted only the minimum permissions necessary to perform its function. This reduces the attack surface and limits potential damage from compromised dependencies or malicious code.

## Workflow Permissions

### Global Permissions (`.github/workflows/ci.yml`)

```yaml
permissions:
  contents: read        # Allow reading repository contents
  checks: write        # Allow writing check runs (test results)
  pull-requests: write # Allow commenting on PRs (for coverage reports)
```

**Why these permissions?**
- `contents: read` - Workflows need to checkout and read the code
- `checks: write` - Required to post test results and status checks
- `pull-requests: write` - Needed for automated PR comments (e.g., coverage reports)

**NOT granted:**
- `contents: write` - Workflows cannot modify code or push commits
- `actions: write` - Cannot modify workflow files
- `deployments: write` - Cannot trigger deployments
- `packages: write` - Cannot publish packages
- `pages: write` - Cannot deploy to GitHub Pages
- `security-events: write` - Cannot modify security alerts

## Job-Level Permissions

### Test Job

```yaml
permissions:
  contents: read        # Read repository code
  checks: write        # Write test results
  pull-requests: write # Comment coverage on PRs
```

**Purpose:** Run unit tests, generate coverage, post results to PR

**Potential risks mitigated:**
- Cannot push malicious code back to repository
- Cannot modify other workflow files
- Cannot access repository secrets beyond what's explicitly passed

### Lint Job

```yaml
permissions:
  contents: read  # Read repository code
  checks: write   # Write lint check results
```

**Purpose:** Run code quality checks and report issues

**Reduced permissions:**
- Does NOT have `pull-requests: write` - doesn't need to comment on PRs
- Only reports check status via GitHub Checks API

### Build Job

```yaml
permissions:
  contents: read  # Read repository code
  checks: write   # Write build check results
```

**Purpose:** Compile the server binary and upload as artifact

**Notes:**
- Artifact uploads use `GITHUB_TOKEN` automatically
- Cannot publish artifacts to external registries
- Artifacts are scoped to the workflow run

## Security Best Practices

### ✅ What We're Doing Right

1. **Explicit permissions** - No implicit permission grants
2. **Read-only by default** - Start with minimal access
3. **Job-specific permissions** - Each job gets only what it needs
4. **No workflow modification** - Workflows can't alter themselves
5. **Limited PR access** - Can comment but not approve/merge
6. **Artifact isolation** - Build artifacts can't access secrets

### 🔒 Additional Protections

1. **Branch Protection Rules** (recommended)
   - Require status checks to pass
   - Require pull request reviews
   - Restrict who can push to main

2. **Secret Management**
   - Never log secrets in output
   - Use GitHub's secret masking
   - Rotate secrets regularly
   - Scope secrets to specific workflows when possible

3. **Dependency Security**
   - Pin action versions to commit SHAs (not tags)
   - Review third-party actions before use
   - Enable Dependabot for action updates

4. **Code Review**
   - Review workflow changes carefully
   - Watch for permission escalation attempts
   - Audit third-party action usage

## Permission Reference

### Available GitHub Permissions

| Permission | Read | Write | Admin | Our Usage |
|------------|------|-------|-------|-----------|
| actions | ❌ | ❌ | ❌ | None |
| checks | ✅ | ✅ | ❌ | Write (all jobs) |
| contents | ✅ | ❌ | ❌ | Read-only |
| deployments | ❌ | ❌ | ❌ | None |
| issues | ❌ | ❌ | ❌ | None |
| packages | ❌ | ❌ | ❌ | None |
| pages | ❌ | ❌ | ❌ | None |
| pull-requests | ❌ | ✅ | ❌ | Write (test job only) |
| repository-projects | ❌ | ❌ | ❌ | None |
| security-events | ❌ | ❌ | ❌ | None |
| statuses | ❌ | ❌ | ❌ | None |

## Monitoring & Auditing

### How to Monitor Workflow Permissions

1. **View Workflow Runs**
   - Go to Actions tab → Select workflow run
   - Check "Set up job" logs for granted permissions

2. **Audit Logs**
   - Organization Settings → Audit log
   - Filter by actions: `workflows.*`

3. **Security Alerts**
   - Enable Dependabot alerts for Actions
   - Review CodeQL security scanning results

### Red Flags to Watch For

⚠️ **Permission escalation attempts:**
- Workflows requesting `contents: write`
- Unexpected `pull-requests: write` on jobs that don't need it
- Any request for `actions: write` or admin permissions

⚠️ **Suspicious activity:**
- Unexpected workflow file changes
- New workflows without review
- Third-party actions from unknown sources
- Actions pinned to tags instead of SHAs

## Incident Response

If you suspect a security issue:

1. **Immediate Actions**
   - Disable the affected workflow
   - Rotate any potentially exposed secrets
   - Review recent workflow runs

2. **Investigation**
   - Check audit logs for unauthorized changes
   - Review workflow file history
   - Examine workflow run logs

3. **Remediation**
   - Remove malicious code
   - Update permissions to be more restrictive
   - Document the incident
   - Notify team

## Resources

- [GitHub Actions Security Best Practices](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
- [Automatic token authentication](https://docs.github.com/en/actions/security-guides/automatic-token-authentication)
- [Permissions for the GITHUB_TOKEN](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#permissions-for-the-github_token)
- [Security hardening for GitHub Actions](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)

## Review Schedule

This document and workflow permissions should be reviewed:
- Quarterly as part of regular security review
- Whenever new jobs are added to workflows
- After any security incident
- When GitHub updates permission model

---

**Last Updated:** October 2025  
**Owner:** Engineering Team  
**Next Review:** January 2026
