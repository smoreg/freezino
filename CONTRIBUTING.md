# Contributing to Freezino

Thank you for your interest in contributing to Freezino! This educational project aims to combat gambling addiction through awareness, and we welcome contributions from the community.

## 📜 Code of Conduct

### Our Standards

- Be respectful and inclusive
- Welcome newcomers and beginners
- Focus on constructive feedback
- Respect different viewpoints and experiences
- Prioritize the educational mission of the project

### Unacceptable Behavior

- Harassment or discrimination
- Trolling or inflammatory comments
- Publishing others' private information
- Any conduct that would be inappropriate in a professional setting

## 🚀 How to Contribute

### Reporting Bugs

Before creating a bug report:
1. Check the [existing issues](https://github.com/smoreg/freezino/issues)
2. Update to the latest version to see if the bug persists
3. Collect relevant information (browser, OS, error messages)

When creating a bug report, include:
- **Clear title** describing the issue
- **Steps to reproduce** the bug
- **Expected behavior** vs **actual behavior**
- **Screenshots** if applicable
- **Environment details** (OS, browser, versions)
- **Error logs** from console or server

### Suggesting Features

We welcome feature suggestions that align with our educational mission!

Before suggesting a feature:
1. Check [existing feature requests](https://github.com/smoreg/freezino/issues?q=label%3Aenhancement)
2. Review the [PLAN.md](./PLAN.md) to ensure it fits the project vision

When suggesting a feature:
- **Describe the problem** it solves
- **Explain the solution** you envision
- **Consider alternatives** you've thought about
- **Explain educational value** if applicable
- **Add mockups** if you have design ideas

### Pull Request Process

1. **Fork the repository** and create a branch from `main`
2. **Make your changes** following our coding standards
3. **Test your changes** thoroughly
4. **Update documentation** if needed
5. **Commit with clear messages** (see commit guidelines below)
6. **Create a pull request** with a clear description

#### Pull Request Guidelines

- Link to related issues
- Describe what changes you made and why
- Include screenshots for UI changes
- Ensure all tests pass
- Keep changes focused (one feature/fix per PR)
- Update relevant documentation

## 💻 Development Workflow

### Setting Up Development Environment

1. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/freezino.git
   cd freezino
   ```

2. **Set up upstream remote**:
   ```bash
   git remote add upstream https://github.com/smoreg/freezino.git
   ```

3. **Install dependencies**:
   ```bash
   # Backend
   cd backend
   make install

   # Frontend
   cd ../frontend
   npm install
   ```

4. **Create a branch**:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

5. **Start development servers**:
   ```bash
   # Backend (in backend/)
   make dev

   # Frontend (in frontend/)
   npm run dev
   ```

### Keeping Your Fork Updated

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

## 📝 Coding Standards

### General Principles

- **Clean Code**: Write readable, maintainable code
- **DRY**: Don't Repeat Yourself
- **KISS**: Keep It Simple, Stupid
- **YAGNI**: You Aren't Gonna Need It
- **Comments**: Explain *why*, not *what*
- **Security**: Always validate input and sanitize output

### Go Backend Standards

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Keep functions small and focused
- Use meaningful variable names
- Add comments for exported functions
- Handle errors explicitly (no silent failures)

Example:
```go
// CalculateWinnings calculates payout based on bet type and amount
func CalculateWinnings(betType BetType, amount float64, winningNumber int) float64 {
    multiplier := getMultiplier(betType, winningNumber)
    return amount * multiplier
}
```

**Run linters**:
```bash
cd backend
make fmt    # Format code
make lint   # Run linter
make test   # Run tests
```

### React Frontend Standards

- Use TypeScript for type safety
- Follow React best practices and hooks guidelines
- Use functional components (no class components)
- Keep components small and focused
- Use custom hooks for reusable logic
- Prefer composition over prop drilling
- Use Zustand for state management
- Follow ESLint rules

Example:
```typescript
interface GameCardProps {
  title: string;
  imageUrl: string;
  onPlay: () => void;
}

export const GameCard: React.FC<GameCardProps> = ({ title, imageUrl, onPlay }) => {
  return (
    <div className="game-card" onClick={onPlay}>
      <img src={imageUrl} alt={title} />
      <h3>{title}</h3>
    </div>
  );
};
```

**Run linters**:
```bash
cd frontend
npm run lint    # ESLint
npm run build   # Type check
```

### CSS/TailwindCSS Standards

- Use TailwindCSS utility classes
- Follow mobile-first responsive design
- Use semantic class names for custom CSS
- Keep custom CSS minimal
- Use CSS variables for theme colors
- Ensure accessibility (a11y)

### Commit Message Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/):

**Format**:
```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

**Examples**:
```
feat(games): add blackjack game with WebSocket support

fix(auth): prevent token refresh loop on 401 errors

docs(readme): update installation instructions for Docker

refactor(shop): extract item filtering logic to custom hook

test(work): add unit tests for work session service
```

## 🧪 Testing Guidelines

### Backend Tests

- Write unit tests for business logic
- Write integration tests for API endpoints
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for >70% code coverage

Example:
```go
func TestCalculateWinnings(t *testing.T) {
    tests := []struct {
        name     string
        betType  BetType
        amount   float64
        number   int
        expected float64
    }{
        {"straight bet win", StraightBet, 10, 7, 350},
        {"red bet win", ColorBet, 10, 1, 20},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculateWinnings(tt.betType, tt.amount, tt.number)
            assert.Equal(t, tt.expected, got)
        })
    }
}
```

### Frontend Tests

- Write unit tests for utilities and hooks
- Write component tests with Testing Library
- Write E2E tests for critical user flows
- Test accessibility
- Test responsive design

Example:
```typescript
describe('GameCard', () => {
  it('renders game title and image', () => {
    const { getByText, getByAltText } = render(
      <GameCard title="Roulette" imageUrl="/roulette.png" onPlay={() => {}} />
    );

    expect(getByText('Roulette')).toBeInTheDocument();
    expect(getByAltText('Roulette')).toBeInTheDocument();
  });

  it('calls onPlay when clicked', () => {
    const onPlay = jest.fn();
    const { getByText } = render(
      <GameCard title="Roulette" imageUrl="/roulette.png" onPlay={onPlay} />
    );

    fireEvent.click(getByText('Roulette'));
    expect(onPlay).toHaveBeenCalled();
  });
});
```

## 📁 Project Structure

Understand the codebase structure before contributing:

```
freezino/
├── backend/              # Go backend
│   ├── cmd/server/      # Main entry point
│   ├── internal/
│   │   ├── auth/        # Authentication logic
│   │   ├── config/      # Configuration
│   │   ├── database/    # DB setup and migrations
│   │   ├── handler/     # HTTP handlers
│   │   ├── middleware/  # Middleware
│   │   ├── model/       # Database models
│   │   ├── router/      # Routes
│   │   └── service/     # Business logic
│   └── ...
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # Reusable components
│   │   ├── pages/       # Page components
│   │   ├── store/       # Zustand stores
│   │   ├── hooks/       # Custom hooks
│   │   ├── utils/       # Utilities
│   │   └── i18n/        # Translations
│   └── ...
└── docs/                # Documentation
```

## 🌍 Internationalization (i18n)

When adding new UI text:

1. **Never hardcode text** in components
2. **Add translations** to all language files
3. **Use i18next** translation hooks
4. **Test with all languages**

Example:
```typescript
// ❌ Bad
<button>Start Game</button>

// ✅ Good
import { useTranslation } from 'react-i18next';

const { t } = useTranslation();
<button>{t('games.startGame')}</button>
```

Add to translation files:
```json
// frontend/src/i18n/locales/en.json
{
  "games": {
    "startGame": "Start Game"
  }
}

// frontend/src/i18n/locales/ru.json
{
  "games": {
    "startGame": "Начать игру"
  }
}
```

## 🔐 Security Guidelines

- **Never commit secrets** (API keys, passwords, etc.)
- **Validate all user input** on backend
- **Sanitize output** to prevent XSS
- **Use parameterized queries** to prevent SQL injection
- **Implement rate limiting** for sensitive endpoints
- **Use HTTPS** in production
- **Keep dependencies updated**
- **Follow OWASP** security best practices

## 📚 Documentation

- Update relevant documentation with your changes
- Add JSDoc/GoDoc comments for public APIs
- Update README if you add new features
- Add examples for complex functionality
- Keep API documentation (OpenAPI) up to date

## 🎯 Educational Mission

Remember Freezino's core mission:

- **Educational First**: Features should teach about gambling risks
- **No Real Money**: Never implement real money transactions
- **Responsible Gaming**: Promote awareness and prevention
- **Statistical Transparency**: Show real-world comparisons
- **User Safety**: Prioritize user wellbeing

## 📋 Checklist Before Submitting PR

- [ ] Code follows project coding standards
- [ ] All tests pass (`make test`, `npm run test`)
- [ ] No linting errors (`make lint`, `npm run lint`)
- [ ] Code is properly formatted
- [ ] Added/updated tests for changes
- [ ] Updated documentation
- [ ] Added translations for new UI text
- [ ] Tested manually in browser
- [ ] Tested responsive design
- [ ] No console errors or warnings
- [ ] Commit messages follow conventions
- [ ] PR description is clear and complete

## 🤝 Getting Help

- **Questions**: Open a [discussion](https://github.com/smoreg/freezino/discussions)
- **Bugs**: Open an [issue](https://github.com/smoreg/freezino/issues)
- **Chat**: Join our community (link TBD)
- **Email**: Contact maintainers (see README)

## 📄 License

By contributing to Freezino, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Freezino and helping promote responsible gaming education! 🎓
