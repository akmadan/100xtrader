# Contributing to 100xtrader

Thank you for your interest in contributing to 100xtrader! 🎉

## Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/akmadan/100xtrader.git
   cd 100xtrader
   ```

3. **Set up development environment**
   ```bash
   # Using Docker (Recommended)
   docker-compose up -d
   
   # Or locally
   make install
   make dev
   ```

## Development Workflow

1. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write clean, readable code
   - Follow existing code style
   - Add comments for complex logic
   - Update documentation if needed

3. **Test your changes**
   ```bash
   # Backend tests
   cd go-core && go test ./...
   
   # Frontend tests (when available)
   cd web && npm test
   ```

4. **Commit your changes**
   ```bash
   git commit -m "feat: add your feature description"
   ```
   
   Use conventional commit messages:
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation
   - `style:` - Code style changes
   - `refactor:` - Code refactoring
   - `test:` - Adding tests
   - `chore:` - Maintenance tasks

5. **Push and create Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```

## Code Style

### Go (Backend)
- Follow Go conventions
- Use `gofmt` for formatting
- Keep functions focused and small
- Add comments for exported functions

### TypeScript/React (Frontend)
- Use TypeScript for type safety
- Follow React best practices
- Use functional components with hooks
- Keep components small and focused

## Pull Request Guidelines

- Provide a clear description of changes
- Reference related issues
- Ensure all tests pass
- Update documentation if needed
- Keep PRs focused (one feature/fix per PR)

## Questions?

Open an issue or start a discussion in GitHub Discussions!

