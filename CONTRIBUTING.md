# Contributing to deployes

First off, thank you for considering contributing to deployes! It's people like you that make deployes such a great tool.

## 🌟 Ways to Contribute

- **Bug Reports**: Submit a bug report
- **Feature Requests**: Suggest new features
- **Code Contributions**: Submit a Pull Request
- **Documentation**: Improve documentation
- **Translations**: Add or improve translations

## 🐛 Bug Reports

If you find a bug, please create an [issue](https://github.com/aliicolak/deployes/issues) with:

1. **Clear title**: Brief description of the bug
2. **Steps to reproduce**: Detailed steps
3. **Expected behavior**: What should happen
4. **Actual behavior**: What actually happens
5. **Environment**: OS, Go version, Node version, etc.
6. **Screenshots**: If applicable

## 💡 Feature Requests

We love new ideas! Please create an [issue](https://github.com/aliicolak/deployes/issues) with:

1. **Clear title**: Brief description
2. **Use case**: Why is this feature needed?
3. **Proposed solution**: How should it work?
4. **Alternatives**: Other solutions you've considered

## 🔧 Development Setup

### Prerequisites

- **Go 1.24+**
- **Node.js 20+**
- **PostgreSQL 15+**
- **Docker & Docker Compose**
- **Make** (optional)

### Local Development

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/deployes.git
   cd deployes
   ```

2. **Install Dependencies**
   ```bash
   # Backend
   go mod download

   # Frontend
   cd web
   npm install
   ```

3. **Set Up Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

4. **Start Database**
   ```bash
   docker compose up -d postgres
   ```

5. **Run Backend**
   ```bash
   make run
   # or
   go run ./cmd/api
   ```

6. **Run Frontend**
   ```bash
   cd web
   npm start
   ```

## 📝 Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `go vet` before committing
- Write tests for new code
- Maintain test coverage above 70%

### TypeScript/Angular Code

- Follow [Angular Style Guide](https://angular.io/guide/styleguide)
- Use Prettier for formatting
- Write meaningful component and service names
- Add tests for components and services

### General Guidelines

- **Write clear commit messages**
  ```
  type(scope): subject

  body

  footer
  ```
  Types: feat, fix, docs, style, refactor, test, chore

- **Keep PRs small and focused**
- **Update documentation** for new features
- **Add tests** for new functionality

## 🧪 Testing

### Backend Tests
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/application/deployment -v
```

### Frontend Tests
```bash
cd web
npm test
```

## 🔍 Code Quality

### Linting
```bash
# Backend
make lint

# Frontend
cd web
npm run lint
```

### Pre-commit Hooks
```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Run manually
pre-commit run --all-files
```

## 📥 Pull Request Process

1. **Create a Branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

2. **Make Changes**
   - Write clean, tested code
   - Follow coding standards
   - Update documentation

3. **Commit Changes**
   ```bash
   git commit -m "feat: add amazing feature"
   ```

4. **Push to GitHub**
   ```bash
   git push origin feature/amazing-feature
   ```

5. **Open Pull Request**
   - Fill out PR template
   - Link related issues
   - Request review from maintainers

6. **Address Review Feedback**
   - Make requested changes
   - Push new commits
   - Keep discussion professional

### PR Requirements

- ✅ All tests pass
- ✅ Code coverage maintained
- ✅ No linting errors
- ✅ Documentation updated
- ✅ Commit messages follow convention
- ✅ At least 1 approval from maintainer

## 🏗️ Project Structure

```
deployes/
├── cmd/                    # Application entry points
│   └── api/               # API server
├── internal/
│   ├── application/       # Business logic
│   ├── domain/           # Domain entities
│   ├── infrastructure/   # DB, external services
│   └── interfaces/       # HTTP handlers
├── pkg/                   # Public utilities
├── web/                   # Angular frontend
└── docs/                  # Documentation
```

## 📚 Documentation

- Update README.md if needed
- Add inline comments for complex logic
- Update API documentation
- Add examples for new features

## 🌍 Translations

We support multiple languages! To add a new translation:

1. Copy `web/src/assets/i18n/en.json`
2. Rename to your language code (e.g., `fr.json`)
3. Translate all strings
4. Add to `app.config.ts`:
   ```typescript
   { code: 'fr', name: 'Français' }
   ```

## 🙋‍♀️ Questions?

- **GitHub Issues**: For bugs and features
- **Discussions**: For questions and ideas
- **Email**: alicolak1988@hotmail.com

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing! 🎉

Made with ❤️ by the deployes team
