# 🔒 Security Guidelines for Smor-Ting

This document outlines security best practices and requirements for the Smor-Ting project.

## ⚠️ Critical Security Requirements

### 1. Environment Variables & Secrets

**✅ DO:**
- Use `.env.example` as a template for environment configuration
- Store actual credentials in `.env` files (which are gitignored)
- Use cryptographically secure random values for secrets
- Generate secrets using: `openssl rand -base64 32`
- Use environment variables for all sensitive configuration

**❌ DON'T:**
- Commit `.env` files to version control
- Use default/placeholder values in production
- Hardcode credentials in source code
- Share credentials in documentation files
- Use weak or predictable secrets

### 2. Credential Management

#### MongoDB Credentials
```bash
# Generate a secure MongoDB password
openssl rand -base64 32

# Use MongoDB Atlas connection string format
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/database?retryWrites=true&w=majority
```

#### JWT Secrets
```bash
# Generate secure JWT secrets (minimum 32 bytes)
JWT_SECRET=$(openssl rand -base64 32)
JWT_ACCESS_SECRET=$(openssl rand -base64 32)
JWT_REFRESH_SECRET=$(openssl rand -base64 32)
```

#### Encryption Keys
```bash
# Generate encryption keys for sensitive data
ENCRYPTION_KEY=$(openssl rand -base64 32)
PAYMENT_ENCRYPTION_KEY=$(openssl rand -base64 32)
```

### 3. Security Validation

The application automatically validates:
- ✅ No default/placeholder values in production
- ✅ Minimum 32-character length for secrets
- ✅ Secure database connection strings
- ✅ Proper SSL/TLS configuration in production
- ✅ Strong secret patterns (no dictionary words, etc.)

### 4. Database Security

#### MongoDB Atlas Security Checklist
- [ ] Enable database authentication
- [ ] Use strong, unique passwords
- [ ] Enable SSL/TLS connections
- [ ] Configure IP whitelisting
- [ ] Enable database audit logging
- [ ] Use role-based access control
- [ ] Regular security updates

#### Connection Security
```bash
# Production connections MUST use SSL
MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net/db?ssl=true&retryWrites=true&w=majority

# Enable Atlas cluster mode
MONGODB_ATLAS=true
DB_SSL_MODE=require
```

## 🛡️ Security Scanning

### Automated Security Checks

Our CI/CD pipeline includes:
- **Secret Detection**: GitGuardian & TruffleHog scan for exposed credentials
- **Vulnerability Scanning**: Dependency vulnerability checks
- **Code Security**: Semgrep static analysis
- **Infrastructure**: Dockerfile and infrastructure security
- **Environment Validation**: Ensures no .env files are committed

### Manual Security Audits

Run these commands locally before committing:

```bash
# Check for secrets in code
git diff --name-only | xargs grep -E "(password|secret|key|token)" --color=always

# Validate environment configuration
cd smor_ting_backend && go run cmd/main.go --validate-config

# Check for exposed .env files
git ls-files | grep -E '\.(env|environment)($|\.)'
```

## 🚨 Security Incident Response

### If Credentials Are Exposed

1. **Immediate Actions:**
   ```bash
   # Rotate ALL exposed credentials immediately
   # For MongoDB Atlas:
   # 1. Change database user password
   # 2. Update connection string
   # 3. Restart all services
   
   # For JWT secrets:
   JWT_SECRET=$(openssl rand -base64 32)
   JWT_ACCESS_SECRET=$(openssl rand -base64 32)
   JWT_REFRESH_SECRET=$(openssl rand -base64 32)
   ```

2. **Remove from Git History:**
   ```bash
   # Use BFG Repo-Cleaner to remove secrets from history
   java -jar bfg.jar --replace-text passwords.txt
   git reflog expire --expire=now --all && git gc --prune=now --aggressive
   ```

3. **Update Security:**
   - Force logout all users (JWT invalidation)
   - Monitor for unauthorized access
   - Update security scanning rules
   - Review access logs

### Credential Rotation Schedule

- **JWT Secrets**: Every 90 days or immediately if compromised
- **Database Passwords**: Every 90 days
- **API Keys**: Every 180 days or per provider recommendation
- **Encryption Keys**: Every 365 days or immediately if compromised

## 📋 Security Checklist

### Development Environment
- [ ] `.env` file created from `.env.example`
- [ ] All default values replaced with secure values
- [ ] `.env` file listed in `.gitignore`
- [ ] No credentials in source code

### Production Deployment
- [ ] All environment variables set securely
- [ ] Database authentication enabled
- [ ] SSL/TLS connections configured
- [ ] Security validation passes
- [ ] Monitoring and alerting configured

### Code Review
- [ ] No hardcoded credentials
- [ ] Environment variables used for configuration
- [ ] Security validation tests pass
- [ ] Documentation updated

## 🔗 Security Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [MongoDB Security Checklist](https://docs.mongodb.com/manual/administration/security-checklist/)
- [JWT Security Best Practices](https://auth0.com/blog/a-look-at-the-latest-draft-for-jwt-bcp/)
- [Go Security Guide](https://github.com/OWASP/Go-SCP)

## 📞 Security Contacts

For security vulnerabilities or concerns:
- **Internal**: Development team
- **External**: Create a private GitHub issue
- **Emergency**: Contact system administrators

---

**Remember**: Security is everyone's responsibility. When in doubt, ask questions and err on the side of caution.
